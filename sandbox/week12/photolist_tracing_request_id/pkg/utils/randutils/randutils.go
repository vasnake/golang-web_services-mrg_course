package randutils

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandBytesHex(n int) string {
	return fmt.Sprintf("%x", RandBytes(n))
}

func RandCryptBytesHex(n int) string {
	return fmt.Sprintf("%x", RandCryptBytes(n))
}

func RandBytes(n int) []byte {
	res := make([]byte, n)
	rand.Read(res)
	return res
}

func RandCryptBytes(n int) []byte {
	res := make([]byte, n)
	crand.Read(res)
	return res
}
