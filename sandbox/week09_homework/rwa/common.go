package main

import (
	"fmt"
	"log"
	"math/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// https://gist.github.com/ik5/d8ecde700972d4378d87

var (
	Info = LogTeal
	Warn = LogYellow
	Fata = LogRed
)

var (
	LogBlack   = Color("\033[1;30m%s\033[0m")
	LogRed     = Color("\033[1;31m%s\033[0m")
	LogGreen   = Color("\033[1;32m%s\033[0m")
	LogYellow  = Color("\033[1;33m%s\033[0m")
	LogPurple  = Color("\033[1;34m%s\033[0m")
	LogMagenta = Color("\033[1;35m%s\033[0m")
	LogTeal    = Color("\033[1;36m%s\033[0m")
	LogWhite   = Color("\033[1;37m%s\033[0m")
)

func Color(colorString string) func(...interface{}) {
	sprint := func(args ...interface{}) {
		log.Printf(colorString, fmt.Sprint(args...))
	}
	return sprint
}
