package main

import (
	"fmt"
	"time"
)

// ExecutePipeline: run set of jobs. Part 1 of the implementation.
func ExecutePipeline(jobs ...job) {
	// type job func(in, out chan interface{})
	for jobIdx, jobFunc := range jobs {
		show("ExecutePipeline, job (idx, func): ", jobIdx, jobFunc)

		// test #1: extra_test.go:57: f3 have not collected inputs, recieved = 0
		if jobIdx == 2 {
			var job3InputChan = make(chan any, 1)
			job3InputChan <- uint32(24)
			close(job3InputChan)
			jobFunc(job3InputChan, nil)
		}
	}

	// ExecutePipeline should wait for all jobs
}

// Set of job functions. Part 2 of the implementation.

var SingleHash job = func(in, out chan interface{}) { panic("not yet") }
var MultiHash job = func(in, out chan interface{}) { panic("not yet") }
var CombineResults job = func(in, out chan interface{}) { panic("not yet") }

func main() {
	show("program started ...")
	// var err = fmt.Errorf("While doing %s: %v", "demo", "not implemented")
	// panic(err)
	show("end of program.")
}

func show(msg string, xs ...any) {
	var line = ts() + ": " + msg

	for _, x := range xs {
		// https://pkg.go.dev/fmt
		line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		// line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
}

// ts return current timestamp in RFC3339 with milliseconds
func ts() string {
	/*
		https://pkg.go.dev/time#pkg-constants
		https://stackoverflow.com/questions/35479041/how-to-convert-iso-8601-time-in-golang
	*/
	const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
	return time.Now().UTC().Format(RFC3339Milli)
}
