package main

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

// defined in other file
/*
type TCache struct {
	*memcache.Client
}
*/

func main() {
	MemcachedAddresses := []string{"127.0.0.1:11211"}
	memcacheClient := memcache.New(MemcachedAddresses...)

	tc := &TCache{memcacheClient}

	mkey := "habrposts"
	tc.Delete(mkey) // start with empty cache

	// universal interface, retrieve items from backend (not from cache), func => (items, tags, error)
	rebuild := func() (interface{}, []string, error) {
		habrPosts, err := GetHabrPosts()
		if err != nil {
			return nil, nil, err
		}
		return habrPosts, []string{"habrTag", "geektimes"}, nil
	}

	// get items, first call, cache miss
	fmt.Println("\nTGet call #1")
	posts := RSS{}
	err := tc.TGet(mkey, 30, &posts, rebuild)
	fmt.Println("#1", len(posts.Items), "err:", err)

	// get items, second call, from cache
	fmt.Println("\nTGet call #2")
	posts = RSS{}
	err = tc.TGet(mkey, 30, &posts, rebuild)
	fmt.Println("#2", len(posts.Items), "err:", err)

	// set new value in some tags => cache values invalid
	fmt.Println("\ninc tag habrTag")
	tc.Increment("habrTag", 1)

	// async get items, cache invalid and shoud be rebuilded
	go func() {
		// time.Sleep(time.Millisecond)
		fmt.Println("\nTGet call #async")
		posts = RSS{}
		err = tc.TGet(mkey, 30, &posts, rebuild)
		fmt.Println("#async", len(posts.Items), "err:", err)
	}()

	// concurrent get items, who's first get the lock?
	fmt.Println("\nTGet call #3")
	posts = RSS{}
	err = tc.TGet(mkey, 30, &posts, rebuild)
	fmt.Println("#3", len(posts.Items), "err:", err)
}
