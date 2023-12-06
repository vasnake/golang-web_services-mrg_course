package main

import (
	"testing"
)

func BenchmarkXmlDoc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = findAllLoginsInDoc()
	}
}

func BenchmarkXmlStream(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = findAllLoginsInStream()
	}
}
