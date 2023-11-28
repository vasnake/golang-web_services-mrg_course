package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ExecutePipeline: run set of jobs.
func ExecutePipeline(jobs ...job) {
	// implementation notes:
	// type job func(in, out chan interface{})
	// first in = nil; last out = nil; 'next in' = 'previous out'
	// close channels when the first job finished (all input data processed)
	show("ExecutePipeline, creating pipeline from x jobs, x: ", len(jobs))

	if len(jobs) < 1 {
		var err = fmt.Errorf("While doing %s: %v", "ExecutePipeline", "number of jobs < 1")
		panic(err)
	}

	var firstJob func()                          // run after launching all others
	var currInput, currOutput chan any           // current job pipes
	var pipes = make([]chan any, 0, len(jobs)-1) // all pipes

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

			// create-and-save first-job-closure
			var firstOutput = currOutput
			var firstJobFunc = jobFunc
			firstJob = func() { firstJobFunc(nil, firstOutput) }

		} else { // second and next and next ...
			currInput = currOutput

			if (jobIdx + 1) == len(jobs) {
				show("ExecutePipeline, creating last job ...")
				currOutput = nil
			} else {
				show("ExecutePipeline, creating intermediate job, idx: ", jobIdx)
				currOutput = createPipe()
			}

			show("ExecutePipeline, starting async job, idx: ", jobIdx)
			go jobFunc(currInput, currOutput)
		}
	}

	// ExecutePipeline should wait for the first job
	show("ExecutePipeline, start first job ...")
	firstJob()
	show("ExecutePipeline, first job done. Pipeline done.")
}

// Set of job functions. Part 2 of the implementation.

var MultiHash job = func(in, out chan interface{}) {
	for inVal := range in {
		out <- computeMultiHash(inVal.(string)) // type assertion: should be replaced with type switch or Sprintf
	}
}

var SingleHash job = func(in, out chan interface{}) {
	var err = fmt.Errorf("While doing %s: %v", "SingleHash", "not implemented")
	panic(err)
}

var CombineResults job = func(in, out chan interface{}) {
	var err = fmt.Errorf("While doing %s: %v", "CombineResults", "not implemented")
	panic(err)
}

func computeMultiHash(text string) string {
	/*
		MultiHash считает значение crc32(th+data)):
		(конкатенация цифры, приведённой к строке и строки),
		где th=0..5 ( т.е. 6 хешей на каждое входящее значение ),
		потом берёт конкатенацию результатов в порядке расчета (0..5),
		где data - то что пришло на вход (и ушло на выход из SingleHash)
	*/
	const partsCount = 6
	var parts = [partsCount]string{}
	var wait = &sync.WaitGroup{}

	var computePart = func(idx int, text string) {
		parts[idx] = DataSignerCrc32(strconv.Itoa(idx) + text)
		wait.Done()
	}

	wait.Add(partsCount)
	for idx := 0; idx < partsCount; idx++ {
		// async
		go computePart(idx, text)
	}
	wait.Wait()

	return strings.Join(parts[:], "")
}

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
