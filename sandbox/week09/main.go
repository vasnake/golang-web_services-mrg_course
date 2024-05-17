package main

import (
	"fmt"
	"net/http"
	"time"

	basic "week09/0_basic"
	storage "week09/1_storage"
	sql_storage "week09/2_sql_storage"
	sess_auth "week09/3_auth"
	csrf_token "week09/4_csrf_token"
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
	// basicPrototypeDemo()
	// storageLayerDemo()
	// sqlStorageDemo()
	// sessionAuthDemo()
	// passSaltDemo()
	csrfTokenDemo()
}

func csrfTokenDemo() {
	show("csrfTokenDemo: program started ...")
	csrf_token.MainCsrf()
	show("end of program. ")
}

func passSaltDemo() {
	show("passSaltDemo: program started ...")
	PassSaltMain1()
	PassSaltMain2()
	show("end of program. ")
}

func sessionAuthDemo() {
	show("sessionAuthDemo: program started ...")
	sess_auth.MainAuth()
	show("end of program. ")
}

func sqlStorageDemo() {
	show("sqlStorageDemo: program started ...")
	sql_storage.MainSqlStorage()
	show("end of program. ")
}

func storageLayerDemo() {
	show("storageLayerDemo: program started ...")
	storage.MainStorage()
	show("end of program. ")
}

func basicPrototypeDemo() {
	show("basicPrototypeDemo: program started ...")
	basic.MainBasic()
	show("end of program. ")
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
