package main

import (
	"encoding/json"
	"testing"
)

var (
	jsonBytes = []byte(`{"RealName":"Vasily", "Login":"v.romanov", "Status":1, "Flags": 1}`)
)

// go test -v -bench=. -benchmem json/*.go

func BenchmarkEncodeStandart(b *testing.B) {
	var (
		c = Client{}
	)
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(&c)
	}
}

func BenchmarkDecodeStandart(b *testing.B) {
	var (
		c = Client{}
	)
	for i := 0; i < b.N; i++ {
		_ = json.Unmarshal(jsonBytes, &c)
	}
}

func BenchmarkEncodeEasyjson(b *testing.B) {
	var (
		u = UserV2{}
	)
	for i := 0; i < b.N; i++ {
		// codegen
		_, _ = u.MarshalJSON()
	}
}

func BenchmarkDecodeEasyjson(b *testing.B) {
	var (
		u = UserV2{}
	)
	for i := 0; i < b.N; i++ {
		// codegen
		_ = u.UnmarshalJSON(jsonBytes)
	}
}
