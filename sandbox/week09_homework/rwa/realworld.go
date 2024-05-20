package main

import (
	"fmt"
	"net/http"
	"time"
)

// сюда писать код

func GetApp() http.Handler {
	return &ConduitAppHttpHandlers{}
}

type ConduitAppHttpHandlers struct {
	storage Storage
}

func (srv *ConduitAppHttpHandlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// app_test.go:518: [0] bad status code, want: 201, have:200
	show("not yet")
}

type Storage interface {
	Set(key, value string) error
	Get(key string) (value string, err error)
}

func panicOnError(msg string, err error) {
	if err != nil {
		panic(msg + ": " + err.Error())
	}
}

// ts returns current timestamp in RFC3339 with milliseconds
func ts() string {
	/*
		https://pkg.go.dev/time#pkg-constants
		https://stackoverflow.com/questions/35479041/how-to-convert-iso-8601-time-in-golang
	*/
	const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
	return time.Now().UTC().Format(RFC3339Milli)
}

// show writes message to standard output. Message combined from prefix msg and slice of arbitrary arguments
func show(msg string, xs ...any) {
	var line = ts() + ": " + msg

	for _, x := range xs {
		// https://pkg.go.dev/fmt
		// line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
}
