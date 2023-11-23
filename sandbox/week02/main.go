package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

func goroutinesDemo() {
	const (
		iterationsNum = 7
		goroutinesNum = 5
	)

	var formatProgressBar = func(taskNum, iterNum int) string {
		return fmt.Sprintln(
			strings.Repeat("  ", taskNum),
			"█",
			strings.Repeat("  ", goroutinesNum-taskNum),
			"goroutine:",
			taskNum,
			"iter:",
			iterNum,
			strings.Repeat("■", iterNum))
	}

	var doSomeWork = func(taskNum int) {
		for iterNum := 0; iterNum < iterationsNum; iterNum++ {
			// Imagine some heavy processing there, could be blocked by inner infinite loop
			var line = formatProgressBar(taskNum, iterNum)
			fmt.Printf(line)

			// return to the scheduler, allow to switch to other tasks; break my own time-slice
			runtime.Gosched()
			time.Sleep(time.Millisecond)
		}
	}

	for taskNum := 0; taskNum < goroutinesNum; taskNum++ {
		// sequential calls, wait for each task
		//doSomeWork(taskNum)

		// concurrent calls, no waiting here
		go doSomeWork(taskNum)
	}

	// wait for user
	fmt.Scanln()
}

func chan_1() {
	ch1 := make(chan int) // unbuffered channel, all written values must be taken from chan, or else: deadlock
	//ch1 := make(chan int, 1) // buffered channel, buffer size = 1 allow to put 1 unread message to chan.

	// start async task
	go func(in chan int) {
		show("anon, Waiting for channel ...") // #2
		v := <-in                             // read from channel, transfer control over data to me
		show("anon, Value from channel: ", v) // #5
	}(ch1)

	// give some time to start reading from channel
	show("parent, doing some heavy stuff, wait for 3 seconds ...") // #1
	time.Sleep(3 * time.Second)

	show("parent, Writing to channel ...") // #3
	ch1 <- 42                              // write to channel, transfer control over data to anon. routine
	// ch1 <- 100500 // if unbuffered: deadlock in program, fatal; or goroutine/mem leak

	show("parent, hit enter please ...") // #4
	fmt.Scanln()                         // wait for user
	/*
	   2023-11-23T16:03:31.037Z: parent, doing some heavy stuff, wait for 3 seconds ...
	   2023-11-23T16:03:31.037Z: anon, Waiting for channel ...
	   2023-11-23T16:03:34.039Z: parent, Writing to channel ...
	   2023-11-23T16:03:34.039Z: parent, hit enter please ...
	   2023-11-23T16:03:34.039Z: anon, Value from channel: int(42);
	*/
}

func chan_2() {
	in := make(chan int, 1) // rw channel, buffered

	// async, start writing to channel
	go func(out chan<- int) { // n.b. signature: write-only channel
		for i := 0; i <= 5; i++ {
			show("generator, before writing: ", i)
			out <- i
			show("generator, after writing: ", i)
		}
		close(out) // signal for reader to stop iterations
		//out <- 12 // panic: send on closed channel: can' write to closed chan
		show("generator, finished.")
	}(in)

	for i := range in { // read from channel
		// i, isChanOpened := <-in
		show("main, readed from chan: ", i)
	}

	//fmt.Scanln()
}

func main() {
	// goroutinesDemo()
	// chan_1()
	chan_2()

	var err = fmt.Errorf("While doing `main`: %v", "not implemented")
	panic(err)
}

func show(msg string, xs ...any) {
	const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
	var ts = time.Now().UTC().Format(RFC3339Milli)

	var line string = fmt.Sprintf("%s: %s", ts, msg)

	for _, x := range xs {
		line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		// line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
}
