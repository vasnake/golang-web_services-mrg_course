package main

import (
	"bytes"
	"encoding/json"
	"sync"
	"testing"
)

const bytesBufferSize = 512 // 64

type PublicPage struct {
	ID          int
	Name        string
	Url         string
	OwnerID     int
	ImageUrl    string
	Tags        []string
	Description string
	Rules       string
}

var (
	coolGolangPublic = PublicPage{
		ID:          1,
		Name:        "CoolGolangPublic",
		Url:         "http://example.com",
		OwnerID:     100500,
		ImageUrl:    "http://example.com/img.png",
		Tags:        []string{"programming", "go", "golang"},
		Description: "Best page about golang programming",
		Rules:       "",
	}

	Pages = []PublicPage{
		coolGolangPublic,
		coolGolangPublic,
		coolGolangPublic,
	}
)

func BenchmarkAllocMemNaive(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {

			bufRef := bytes.NewBuffer(make([]byte, 0, bytesBufferSize)) // no pool, just simple allocation
			_ = json.NewEncoder(bufRef).Encode(Pages)

		}
	})
}

func BenchmarkAllocMemFromPool(b *testing.B) {
	var memPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, bytesBufferSize)) // alloc mem for pool
		},
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {

			bufRef := memPool.Get().(*bytes.Buffer) // allocate using pool
			_ = json.NewEncoder(bufRef).Encode(Pages)

			bufRef.Reset() // return resources to pool
			memPool.Put(bufRef)
		}
	})
}

/*
	go test -bench . -benchmem pool_test.go
*/
