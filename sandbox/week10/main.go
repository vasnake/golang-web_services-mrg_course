package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
	oauth_demo "week10/0_oauth"
	v_oauth "week10/1_photolist_oauth"
	fronte "week10/2_frontend"
	ajax_api "week10/3_ajax"
)

const (
	port    = 8080
	portStr = ":8080"
	host    = "127.0.0.1"
)

func lessonTemplate() {
	show("lessonTemplate: program started ...")
	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
}

func main() {
	// oauthDemo()
	// oauthPhotolist()
	// refactoredFrontPhotolist()
	ajaxApiPhotolist()
}

func ajaxApiPhotolist() {
	show("ajaxApiPhotolist: program started ...")
	ajax_api.MainDemo()
	show("end of program. ")
}

func refactoredFrontPhotolist() {
	show("refactoredFrontPhotolist: program started ...")
	fronte.MainDemo()
	show("end of program. ")
}

func oauthPhotolist() {
	show("oauthPhotolist: program started ...")
	v_oauth.MainDemo()
	show("end of program. ")
}

func oauthDemo() {
	show("oauthDemo: program started ...")
	oauth_demo.MainDemo()
	show("end of program. ")
}

// --- useful little functions ---

var atomicCounter = new(atomic.Uint64)

func nextID() string {
	return strconv.FormatInt(int64(atomicCounter.Add(1)), 36)
}

func panicOnError(msg string, err error) {
	if err != nil {
		panic(msg + ": " + err.Error())
	}
}

func strRef(in string) *string {
	return &in
}

// ts returns current timestamp in RFC3339 with milliseconds
func ts() string {
	/*
		https://pkg.go.dev/time#pkg-constants
		https://stackoverflow.com/questions/35479041/how-to-convert-iso-8601-time-in-golang
	*/
	const (
		RFC3339      = "2006-01-02T15:04:05Z07:00"
		RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
	)
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
