// go test -bench . -benchmem prealloc_test.go
package main

import (
	"testing"
)

const appendIterationsCount = 1000

func BenchmarkAppendNaive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data := make([]int, 0) // capacity = 0
		for j := 0; j < appendIterationsCount; j++ {
			data = append(data, j)
		}
	}
}

func BenchmarkAppendPrealloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data := make([]int, 0, appendIterationsCount) // len=0, capacity=maxLex
		for j := 0; j < appendIterationsCount; j++ {
			data = append(data, j)
		}
	}
}

/*
	go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1 pool_test.go

	go tool pprof main.test.exe cpu.out
	go tool pprof main.test.exe mem.out

	go tool pprof -svg -inuse_space main.test.exe mem.out > mem_is.svg
	go tool pprof -svg -inuse_objects main.test.exe mem.out > mem_io.svg
	go tool pprof -svg main.test.exe cpu.out > cpu.svg

	go tool pprof -png main.test.exe cpu.out > cpu.png
*/
