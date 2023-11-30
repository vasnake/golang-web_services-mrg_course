package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Pipe = chan any

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
		var pipe = make(Pipe, 100) // let's try buffered pipes
		show("ExecutePipeline, added new pipe: ", pipe)
		return pipe
	}

	var lastPipe Pipe = nil

	for jobIdx, jobFunc := range jobs {
		var inPipe Pipe = nil // current job input pipe, nil for first job, prev. job output as current input for non-first jobs
		if jobIdx > 0 {
			inPipe = lastPipe
		}
		lastPipe = createPipe() // current job output pipe

		show("ExecutePipeline, starting job (idx, func, outPipe): ", jobIdx, jobFunc, lastPipe)

		// start async job
		go func(inPipe, outPipe Pipe, jobFunc job, jobIdx int) {
			show("job id started, (id, in-pipe, out-pipe, func): ", jobIdx, inPipe, outPipe, jobFunc)
			jobFunc(inPipe, outPipe)
			close(outPipe) // signal: job done
			show("job id done, id: ", jobIdx)
		}(inPipe, lastPipe, jobFunc, jobIdx)
	}

	show("ExecutePipeline, wait for the last pipe: ", lastPipe)
	for x := range lastPipe {
		show("ExecutePipeline, got pipeline output: ", x)
	}
	show("ExecutePipeline, pipeline done.")
}

// Set of job functions. Part 2 of the implementation.

var SingleHash job = func(in, out Pipe) {
	var inConv = func(x any) string {
		return strconv.Itoa(x.(int)) // type assertion: should be replaced with type switch or Sprintf
	}

	selectedHash(in, out, inConv, computeSingleHash, "SingleHash")
}

var MultiHash job = func(in, out Pipe) {
	var inConv = func(x any) string {
		return x.(string) // type assertion: should be replaced with type switch or Sprintf
	}

	selectedHash(in, out, inConv, computeMultiHash, "MultiHash")
}

var selectedHash = func(in, out Pipe, inConv func(any) string, computeHash func(string) string, funcTag string) {
	show(funcTag + ", started ...")

	var wait = &sync.WaitGroup{}

	var doCompute = func(inVal string, out Pipe) {
		var outVal = computeHash(inVal)
		show(funcTag+" async, computed (in => out): ", inVal, outVal)
		out <- outVal
		show(funcTag + " async, done.")
		wait.Done()
	}

	for inVal := range in {
		show(funcTag+", in value: ", inVal)
		wait.Add(1)
		// async
		go doCompute(inConv(inVal), out)
	}

	show(funcTag + ", waiting ...")
	wait.Wait()
	show(funcTag + ", done.")
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

func computeSingleHash(text string) string {
	/*
	   SingleHash считает значение crc32(data)+"~"+crc32(md5(data))
	   ( конкатенация двух строк через ~),
	   где data - то что пришло на вход (по сути - числа из первой функции)
	*/
	var crc32 = DataSignerCrc32
	var firstPart, secondPart string
	var wait = &sync.WaitGroup{}

	var computeFirstPart = func() {
		firstPart = crc32(text)
		wait.Done()
	}

	var computeSecondPart = func() {
		secondPart = crc32(md5WithMutex(text))
		wait.Done()
	}

	// async
	wait.Add(2)
	go computeFirstPart()
	go computeSecondPart()
	wait.Wait()

	return firstPart + "~" + secondPart
}

var md5Mutex = &sync.Mutex{}
var md5WithMutex = func(text string) string {
	md5Mutex.Lock()
	var result = DataSignerMd5(text)
	md5Mutex.Unlock()
	return result
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
