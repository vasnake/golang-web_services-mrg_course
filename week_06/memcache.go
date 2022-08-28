package main

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	MemcachedAddresses := []string{"127.0.0.1:11211"}
	memcacheClient := memcache.New(MemcachedAddresses...)

	mkey := "coursera"

	// SET
	memcacheClient.Set(&memcache.Item{
		Key:        mkey,
		Value:      []byte("1"),
		Expiration: 3,
	})

	// INCREMENT
	memcacheClient.Increment("habrTag", 1)

	// GET
	item, err := memcacheClient.Get(mkey)
	if err != nil && err != memcache.ErrCacheMiss {
		fmt.Println("MC error", err)
	}

	fmt.Printf("mc value %#v\n", item)

	// DELETE
	memcacheClient.Delete(mkey)

	// GET
	item, err = memcacheClient.Get(mkey)
	if err == memcache.ErrCacheMiss {
		fmt.Println("record not found in MC")
	}

}
