package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Pipe chan any

// ExecutePipeline: run set of jobs.
func ExecutePipeline(jobs ...job) {
	// implementation notes:
	// type job func(in, out chan interface{})
	// first in = nil; last out = nil; 'next in' = 'previous out'
	// close channels when the first job finished (all input data processed)
	show("ExecutePipeline, creating pipeline from n jobs, n: ", len(jobs))

	if len(jobs) < 1 {
		var err = fmt.Errorf("While doing %s: %v", "ExecutePipeline", "number of jobs < 1")
		panic(err)
	}

	var createPipe = func() Pipe {
		var pipe = make(Pipe)
		show("ExecutePipeline, added new pipe: ", pipe)
		return pipe
	}

	var lastPipe Pipe = nil

	for jobIdx, jobFunc := range jobs {
		var inPipe Pipe = nil
		if jobIdx > 0 {
			inPipe = lastPipe
		}
		lastPipe = createPipe()

		show("ExecutePipeline, starting job, (idx, func, outPipe): ", jobIdx, jobFunc, lastPipe)

		// async
		go func(ipPipe, outPipe Pipe, jobFunc job, jobIdx int) {
			show("job id started, (id, pipe): ", jobIdx, outPipe)
			jobFunc(inPipe, outPipe)
			close(outPipe)
			show("job id done, (id, pipe): ", jobIdx, outPipe)
		}(inPipe, lastPipe, jobFunc, jobIdx)
	}

	show("ExecutePipeline, wait for the last pipe: ", lastPipe)
	for x := range lastPipe {
		show("ExecutePipeline, got pipeline output: ", x)
	}
	show("ExecutePipeline, pipeline done.")
}

// ExecutePipeline_invalid is NOT-working first attempt to create a pipeline executor.
// For educational purposes only.
func ExecutePipeline_invalid(jobs ...job) {
	// implementation notes:
	// type job func(in, out chan interface{})
	// first in = nil; last out = nil; 'next in' = 'previous out'
	// close channels when the first job finished (all input data processed)
	show("ExecutePipeline, creating pipeline from x jobs, x: ", len(jobs))

	if len(jobs) < 1 {
		var err = fmt.Errorf("While doing %s: %v", "ExecutePipeline", "number of jobs < 1")
		panic(err)
	}

	var firstJob func()                      // run after launching all others
	var pipes = make([]Pipe, 0, len(jobs)-1) // all pipes

	defer func() {
		// close all pipes after firstJob is finished
		for idx, pipe := range pipes {
			show("ExecutePipeline, closing pipe, (idx, pipe): ", idx, pipe)
			close(pipe)
		}
	}()

	var createPipe = func() Pipe {
		var pipe = make(Pipe)
		pipes = append(pipes, pipe)
		show("ExecutePipeline, added new pipe: ", pipe, pipes)
		return pipe
	}

	for jobIdx, jobFunc := range jobs {
		show("ExecutePipeline, processing job, (idx, func): ", jobIdx, jobFunc)
		var currInput, currOutput Pipe

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

/*
2023-11-28T16:24:25.272Z: ExecutePipeline, pipeline done.
    main_test.go:157: execition too long
        Got: 8.080864273s
        Expected: <3s
--- FAIL: TestSigner (8.08s)
*/

var MultiHash job = func(in, out chan interface{}) {
	show("MultiHash, started ...")
	for inVal := range in {
		show("MultiHash, in value: ", inVal)
		var outVal = computeMultiHash(inVal.(string)) // type assertion: should be replaced with type switch or Sprintf
		show("MultiHash, out value: ", outVal)
		out <- outVal
	}
	show("MultiHash, done.")
}

var SingleHash job = func(in, out chan interface{}) {
	show("SingleHash, started ...")
	for inVal := range in {
		show("SingleHash, in value: ", inVal)
		var outVal = computeSingleHash(strconv.Itoa(inVal.(int))) // type assertion: should be replaced with type switch or Sprintf
		show("SingleHash, out value: ", outVal)
		out <- outVal
	}
	show("SingleHash, done.")
}

var CombineResults job = func(in, out chan interface{}) {
	show("CombineResults, started ...")
	var messages = make([]string, 0, 64)
	for inVal := range in {
		show("CombineResults, in value: ", inVal)
		messages = append(messages, inVal.(string)) // type assertion: should be replaced with type switch or Sprintf
	}
	show("CombineResults, collected messages: ", messages)
	var outVal = computeCombineResults(messages)
	show("CombineResults, out value: ", outVal)
	out <- outVal
	show("CombineResults, done.")
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
	var crc32 = DataSignerCrc32
	var parts = [partsCount]string{}
	var wait = &sync.WaitGroup{}

	var computePart = func(idx int, text string) {
		parts[idx] = crc32(strconv.Itoa(idx) + text)
		wait.Done()
	}

	// async
	wait.Add(partsCount)
	for idx := 0; idx < partsCount; idx++ {
		go computePart(idx, text)
	}
	wait.Wait()

	return strings.Join(parts[:], "")
}

func computeSingleHash(text string) string {
	/*
	   SingleHash считает значение crc32(data)+"~"+crc32(md5(data))
	   ( конкатенация двух строк через ~),
	   где data - то что пришло на вход (по сути - числа из первой функции)
	*/
	var crc32 = DataSignerCrc32
	var md5 = DataSignerMd5
	var firstPart, secondPart string
	var wait = &sync.WaitGroup{}

	var computeFirstPart = func() {
		firstPart = crc32(text)
		wait.Done()
	}

	var computeSecondPart = func() {
		secondPart = crc32(md5(text))
		wait.Done()
	}

	// async
	wait.Add(2)
	go computeFirstPart()
	go computeSecondPart()
	wait.Wait()

	return firstPart + "~" + secondPart
}

func computeCombineResults(lines []string) string {
	/*
	   CombineResults получает все результаты,
	   сортирует (https://golang.org/pkg/sort/),
	   объединяет отсортированный результат через _ (символ подчеркивания) в одну строку
	*/
	slices.Sort(lines)
	return strings.Join(lines, "_")
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
