package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// universal cached item
type CacheItemStore struct {
	Data interface{}
	Tags map[string]int
}

// item with tags for my app
type CacheItem struct {
	Data json.RawMessage
	Tags map[string]int
}

// type for func that retrieve data from source (not cache)
type RebuildFunc func() (interface{}, []string, error)

type TCache struct {
	*memcache.Client
}

func (tc *TCache) TGet(
	mkey string,
	ttl int32,
	in interface{}, // it's out really, method should write data to it
	rebuildCb RebuildFunc,
) (err error) {
	inKind := reflect.ValueOf(in).Kind()
	if inKind != reflect.Ptr {
		return fmt.Errorf("in must be ptr, got %s", inKind)
	}

	tc.checkLock(mkey) // cache lock

	itemRaw, err := tc.Get(mkey)
	if err == memcache.ErrCacheMiss {
		fmt.Println("Record not found in memcache")
		return tc.rebuild(mkey, ttl, in, rebuildCb) // go to source for items, cache them, and return them
	} else if err != nil {
		return err
	}

	// items found in cache

	item := &CacheItem{}
	err = json.Unmarshal(itemRaw.Value, &item)
	if err != nil {
		return err
	}

	tagsValid, err := tc.isTagsValid(item.Tags)
	if err != nil {
		return fmt.Errorf("isTagsValid error %s", err)
	}

	if tagsValid {
		// write data to out buffer (named as 'in', WTF?)
		err = json.Unmarshal(item.Data, &in)
		return err
	}

	// tags are invalid: go to source for items, cache them, and return them
	return tc.rebuild(mkey, ttl, in, rebuildCb)
}

func (tc *TCache) isTagsValid(itemTags map[string]int) (bool, error) {
	tags := make([]string, 0, len(itemTags))
	for tagKey := range itemTags {
		tags = append(tags, tagKey)
	}

	// tags from cache
	curr, err := tc.GetMulti(tags)
	if err != nil {
		return false, err
	}

	// unpack tags from cache
	currentTagsMap := make(map[string]int, len(curr))
	for tagKey, tagItem := range curr {
		i, err := strconv.Atoi(string(tagItem.Value))
		if err != nil {
			return false, err
		}
		currentTagsMap[tagKey] = i
	}

	// compare tags (it should be tag:count pair, if not equal => cached item is invalid)
	return reflect.DeepEqual(itemTags, currentTagsMap), nil
}

func (tc *TCache) rebuild(
	mkey string,
	ttl int32,
	in interface{}, // output buffer
	rebuildCb RebuildFunc, // real items getter
) error {
	// cache lock
	tc.lockRebuild(mkey)
	defer tc.unlockRebuild(mkey)

	// go to source, N.B. tags is a []string, not a map[string]int
	result, tags, err := rebuildCb()

	// ожидаем и возвращаем одинаковые типы // check items type
	if reflect.TypeOf(result) != reflect.TypeOf(in) {
		return fmt.Errorf(
			"data type mismatch, expected %s, got %s", reflect.TypeOf(in),
			reflect.TypeOf(result),
		)
	}

	// get tags from cache (if tags not exists => create them with current time as counts)
	currTags, err := tc.getCurrentItemTags(tags, ttl)
	if err != nil {
		return err
	}

	cacheData := CacheItemStore{result, currTags}
	rawData, err := json.Marshal(cacheData)
	if err != nil {
		return err
	}

	// cache actual data
	err = tc.Set(&memcache.Item{
		Key:        mkey,
		Value:      rawData,
		Expiration: int32(ttl),
	})

	// write items to output buffer (in universal way)
	inVal := reflect.ValueOf(in)
	resultVal := reflect.ValueOf(result)
	rv := reflect.Indirect(inVal)
	rvpresult := reflect.Indirect(resultVal)
	rv.Set(rvpresult) // *in = *result

	return nil // no errors
}

func (tc *TCache) checkLock(mkey string) error {
	// only good for demo
	for i := 0; i < 4; i++ {
		_, err := tc.Get("lock_" + mkey)
		if err == memcache.ErrCacheMiss {
			return nil
		}

		if err != nil {
			return err
		}

		time.Sleep(10 * time.Millisecond) // still locked
	}

	// could be still locked
	return nil
}

func (tc *TCache) lockRebuild(mkey string) (bool, error) {
	// try to acquire lock; return (success, errObj)

	// пытаемся взять лок на перестроение кеша
	// чтобы все не ломанулись его перестраивать
	// параметры надо тюнить
	lockKey := "lock_" + mkey
	lockAccuired := false

	for i := 0; i < 4; i++ {
		// add добавляет запись если её ещё нету
		err := tc.Add(&memcache.Item{
			Key:        lockKey,
			Value:      []byte("1"),
			Expiration: int32(3),
		})

		if err == memcache.ErrNotStored {
			// locked already
			fmt.Println("get lock try", i)
			time.Sleep(time.Millisecond * 10)
			continue
		} else if err != nil {
			// oops, actual error
			return false, err
		}

		// set value successfully
		lockAccuired = true
		break
	}

	if !lockAccuired {
		return false, fmt.Errorf("Can't get lock")
	}

	// allright
	return true, nil
}

func (tc *TCache) unlockRebuild(mkey string) {
	// remove lock
	tc.Delete("lock_" + mkey)
}

func (tc *TCache) getCurrentItemTags(tags []string, ttl int32) (map[string]int, error) {
	// get existing tags from cache, if no such tag => set it with current time as count

	// WHY? do we need to read first? We need to write new map.
	currTags, err := tc.GetMulti(tags)
	if err != nil {
		return nil, err
	}

	resultTags := make(map[string]int, len(tags))
	now := int(time.Now().Unix())
	nowBytes := []byte(fmt.Sprint(now))

	for _, tagKey := range tags {
		tagItem, tagExist := currTags[tagKey]
		if !tagExist {
			// no such tag in cache, set new value
			err := tc.Set(&memcache.Item{
				Key:        tagKey,
				Value:      nowBytes,
				Expiration: int32(ttl),
			})
			if err != nil {
				return nil, err
			}

			resultTags[tagKey] = now
		} else {
			// tag exists in cahe, set existing value
			i, err := strconv.Atoi(string(tagItem.Value))
			if err != nil {
				return nil, err
			}

			resultTags[tagKey] = i
		}
	}

	return resultTags, nil
}
