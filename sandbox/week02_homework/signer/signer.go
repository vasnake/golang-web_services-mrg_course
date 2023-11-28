package main

import (
	"fmt"
	"time"
)

// ExecutePipeline: run set of jobs. Part 1 of the implementation.
func ExecutePipeline(jobs ...job) {
	// type job func(in, out chan interface{})
	// first in = nil; last out = nil; 'next in' = 'previous out'
	// close channels when the first job finished (all input data processed)

	if len(jobs) < 1 {
		var err = fmt.Errorf("While doing %s: %v", "ExecutePipeline", "number of jobs < 1")
		panic(err)
	}

	var firstJob func() // run after launching all others
	var currInput, currOutput chan any
	var pipes = make([]chan any, 0, len(jobs)-1)

	defer func() {
		// close all pipes after firstJob is finished
		for idx, pipe := range pipes {
			show("ExecutePipeline, closing pipe, idx: ", idx)
			close(pipe)
		}
	}()

	var createPipe = func() chan any {
		var pipe = make(chan any)
		pipes = append(pipes, pipe)
		return pipe
	}

	for jobIdx, jobFunc := range jobs {
		show("ExecutePipeline, processing job, (idx, func): ", jobIdx, jobFunc)

		if jobIdx == 0 {
			show("ExecutePipeline, creating first job ...")
			currInput = nil
			currOutput = createPipe()
			var firstOutput = currOutput
			var firstJobFunc = jobFunc
			firstJob = func() { firstJobFunc(nil, firstOutput) }
		} else { // second and next and next ...
			currInput = currOutput
			if (jobIdx + 1) == len(jobs) {
				show("ExecutePipeline, creating last job ...")
				currOutput = nil
			} else {
				show("ExecutePipeline, creating job, idx: ", jobIdx)
				currOutput = createPipe()
			}

			show("ExecutePipeline, starting async job, idx: ", jobIdx)
			go jobFunc(currInput, currOutput)
		}
	}

	// ExecutePipeline should wait for first job
	show("ExecutePipeline, start first job ...")
	firstJob()
	show("ExecutePipeline, first job done.")
}

// Set of job functions. Part 2 of the implementation.

var MultiHash job = func(in, out chan interface{}) { panic("not yet") }
var SingleHash job = func(in, out chan interface{}) { panic("not yet") }
var CombineResults job = func(in, out chan interface{}) { panic("not yet") }

/*
panic: not yet

goroutine 22 [running]:
signer.glob..func6(0x0?, 0x0?)
        /mnt/c/Users/valik/data/github/golang-web_services-mrg_course/sandbox/week02_homework/signer/signer.go:71 +0x25
created by signer.ExecutePipeline in goroutine 20
        /mnt/c/Users/valik/data/github/golang-web_services-mrg_course/sandbox/week02_homework/signer/signer.go:58 +0x345

*/

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
