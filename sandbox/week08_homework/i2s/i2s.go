package main

import (
	"fmt"
	"time"
)

func i2s(data interface{}, out interface{}) error {
	show("i2s params: ", data, out)
	// i2s params:
	// map[string]interface {}{"Active":true, "ID":42, "Username":"rvasily"};
	// &main.Simple{ID:0, Username:"", Active:false};

	// iterate over out's fields, set each field from map `data`

	return nil
}

// func userInput(msg string) (res string, err error) {
// 	show(msg)
// 	if n, e := fmt.Scanln(&res); n != 1 || e != nil {
// 		return "", e
// 	}
// 	return res, nil
// }

// func panicOnError(msg string, err error) {
// 	if err != nil {
// 		panic(msg + ": " + err.Error())
// 	}
// }

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
