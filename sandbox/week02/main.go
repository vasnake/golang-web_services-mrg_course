package main

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// goroutinesDemo()
	// chan_1()
	// chan_12()
	// chan_2()
	// select_1()
	// select_2()
	// select_3()
	// timeoutDemo()
	// tickDemo()
	// afterfunc()
	// context_cancel()
	// context_timeout()
	// async_work()
	// workerpool()
	// waitgroupDemo()
	// ratelim()
	// race_1()
	// race_2()
	// atomic_1()
	atomic_2()
}

func goroutinesDemo() {
	// run goroutinesNum tasks, each task execute iterationsNum steps/iterations.
	// sequentially or in async/parallel (in gorutines)
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
		doSomeWork(taskNum)

		// concurrent calls, no waiting on next line
		// go doSomeWork(taskNum)
	}

	userInput("hit ENTER") // must wait for all goroutines ...
}

func chan_1() {
	show("chan_1: start reader, then write")

	ch1 := make(chan int) // unbuffered channel, all written values must be taken from chan, or else: deadlock
	//ch1 := make(chan int, 1) // buffered channel, buffer size = 1 allow to put 1 unread message to chan.

	go func(in chan int) {
		show("anon, reading from channel ...")
		v := <-in
		show("anon, value from channel: ", v)
	}(ch1)

	// give some time to start reading from channel
	show("main, doing some heavy stuff, wait for 3 seconds ...")
	time.Sleep(3 * time.Second)

	show("main, writing to channel ...")
	ch1 <- 42 // write to channel, transfer control over data to anon. routine
	// ch1 <- 100500                          // if unbuffered: deadlock in program, fatal; or goroutine/mem leak

	userInput("hit ENTER") // #4
	/*
		2024-04-09T05:25:23.726Z: chan_1: start reader, then write
		2024-04-09T05:25:23.726Z: main, doing some heavy stuff, wait for 3 seconds ...
		2024-04-09T05:25:23.726Z: anon, reading from channel ...
		2024-04-09T05:25:26.729Z: main, writing to channel ...
		2024-04-09T05:25:26.729Z: hit ENTER
		2024-04-09T05:25:26.729Z: anon, value from channel: int(42);
	*/
}

func chan_12() {
	show("chan_12: start writer, then read")

	ch1 := make(chan int)

	go func(c chan int) {
		show("anon, writing to channel ...")
		c <- 42
		show("anon, written")
	}(ch1)

	show("main, doing some heavy stuff, wait for 3 seconds ...")
	time.Sleep(3 * time.Second) // give some time to start other tasks

	show("main, reading from channel ...")
	v := <-ch1
	show("main, value from channel: ", v)
	// v = <-ch1 // fatal, deadlock

	userInput("hit ENTER")
	/*
		2024-04-09T05:11:20.854Z: chan_12: start writer, then read
		2024-04-09T05:11:20.854Z: main, doing some heavy stuff, wait for 3 seconds ...
		2024-04-09T05:11:20.854Z: anon, writing to channel ...
		2024-04-09T05:11:23.855Z: main, reading from channel ...
		2024-04-09T05:11:23.855Z: main, value from channel: int(42);
		2024-04-09T05:11:23.855Z: hit ENTER
		2024-04-09T05:11:23.855Z: anon, written
	*/
}

func chan_2() {
	show("chan_2: buffered chan, start writer iterations, then read")

	bufferedPipe := make(chan int, 1) // rw channel, buffered

	// async, start writing to channel
	go func(out chan<- int) { // n.b. signature: write-only channel
		show("generator, started ...")
		for i := 0; i <= 5; i++ {
			show("generator, before writing: ", i)
			out <- i
			show("generator, after writing: ", i)
		}

		show("generator, closing channel ...")
		close(out) // signal for reader to stop iterations
		// out <- 12  // panic: send on closed channel
		show("generator, finished.")
	}(bufferedPipe)

	show("main, doing some heavy stuff, wait for 3 seconds ...")
	time.Sleep(3 * time.Second) // give some time to start other tasks

	for i := range bufferedPipe { // read from channel
		// i, isChanOpened := <-in
		show("main, readed from chan: ", i)
	}
	// bufferedPipe <- 12 // panic: send on closed channel

	show("main, done.")
	/*
		2024-04-09T05:34:47.051Z: chan_2: buffered chan, start writer iterations, then read
		2024-04-09T05:34:47.052Z: main, doing some heavy stuff, wait for 3 seconds ...

		2024-04-09T05:34:47.052Z: generator, started ...
		2024-04-09T05:34:47.052Z: generator, before writing: int(0);
		2024-04-09T05:34:47.052Z: generator, after writing: int(0);
		2024-04-09T05:34:47.052Z: generator, before writing: int(1); // all this while main sleep

		2024-04-09T05:34:50.054Z: main, readed from chan: int(0);
		2024-04-09T05:34:50.054Z: main, readed from chan: int(1);
		2024-04-09T05:34:50.054Z: generator, after writing: int(1);
		2024-04-09T05:34:50.054Z: generator, before writing: int(2);
		2024-04-09T05:34:50.054Z: generator, after writing: int(2);
		2024-04-09T05:34:50.054Z: generator, before writing: int(3);
		2024-04-09T05:34:50.054Z: generator, after writing: int(3);
		2024-04-09T05:34:50.054Z: generator, before writing: int(4);
		2024-04-09T05:34:50.054Z: main, readed from chan: int(2);
		2024-04-09T05:34:50.054Z: main, readed from chan: int(3);
		2024-04-09T05:34:50.054Z: main, readed from chan: int(4);
		2024-04-09T05:34:50.054Z: generator, after writing: int(4);
		2024-04-09T05:34:50.054Z: generator, before writing: int(5);
		2024-04-09T05:34:50.054Z: generator, after writing: int(5);
		2024-04-09T05:34:50.054Z: generator, closing channel ...
		2024-04-09T05:34:50.054Z: generator, finished.
		2024-04-09T05:34:50.054Z: main, readed from chan: int(5);
		2024-04-09T05:34:50.054Z: main, done.
	*/
}

func select_1() {
	show("select_1: check and process two channels, one operation in one pass")

	ch1 := make(chan int, 1) // buffered rw channel
	ch1 <- 1                 // if uncomment, case 1 will be selected
	// ch2 := make(chan int) // un-buffered rw channel
	ch2 := make(chan int, 1)
	// try different combinations: buffered, unbuffered, empty chan, not empty chan, ...

	var x = 42

	select {
	case x := <-ch1: // read from chan if it has message
		show("got from ch1, x: ", x)
	case ch2 <- x: // write to chan if it can receive a message, and other cases not ready
		show("put to ch2, x: ", x)
	default: // w/o default you will get deadlock
		show("default: have no ready channels")
	}
	// `select` is not wrapped to a loop, evaluate only once, than exit
	show("done.")
	/*
		2024-04-09T06:17:34.974Z: select_1: check and process two channels, one operation in one pass
		2024-04-09T06:17:34.974Z: put to ch2, x: int(42);
		2024-04-09T06:17:34.974Z: done.
	*/
}

func select_2() {
	show("select_2: check two channels in loop, until both are empty")

	ch1 := make(chan int, 2) // buffered rw chan, fill it up
	ch1 <- 1
	ch1 <- 2

	ch2 := make(chan int, 2) // buffered rw chan, put just one message to it
	ch2 <- 3

	var i = 0
	var x int
LOOP: // labeled block
	for {
		i += 1
		show("loop iteration: ", i)

		select { // check channels in arbitrary order
		case x = <-ch1:
			show("got from ch1, x: ", x)
		case x = <-ch2:
			show("got from ch2, x: ", x)
		default: // if no channels are ready
			show("all channels are empty, bye.")
			break LOOP // exit from labeled loop
		}

		show("end of loop iteration: ", i)
	}
	show("end of program.")
	/*
		2024-04-09T06:22:27.818Z: select_2: check two channels in loop, until both empty
		2024-04-09T06:22:27.818Z: loop iteration: int(1);
		2024-04-09T06:22:27.818Z: got from ch2, x: int(3);
		2024-04-09T06:22:27.818Z: end of loop iteration: int(1);
		2024-04-09T06:22:27.818Z: loop iteration: int(2);
		2024-04-09T06:22:27.818Z: got from ch1, x: int(1);
		2024-04-09T06:22:27.818Z: end of loop iteration: int(2);
		2024-04-09T06:22:27.818Z: loop iteration: int(3);
		2024-04-09T06:22:27.818Z: got from ch1, x: int(2);
		2024-04-09T06:22:27.818Z: end of loop iteration: int(3);
		2024-04-09T06:22:27.818Z: loop iteration: int(4);
		2024-04-09T06:22:27.818Z: all channels are empty, bye.
		2024-04-09T06:22:27.818Z: end of program.
	*/
}

func select_3() {
	show("select_3: produce messages until 'stop' signal")

	stopChan := make(chan bool) // commands chan, unbuffered rw
	dataChan := make(chan int)  // messages chan, unbuffered rw

	// run async messages sender
	go func(stopChannel <-chan bool, dataChannel chan<- int) {
		i, msg := 0, 11
		for { // loop until stop command recieved
			i += 1
			show("task, loop iteration: ", i)

			select {
			// N.B. no default branch, you have to be sure that someone give you the stop signal
			// Think about protocol here ...

			case <-stopChannel: // check commands channel first (not really)
				show("task, got `stop` command, quitting ...")
				close(dataChannel) // writer must close his channel, we are the writer for data chan
				return             // not a good pattern, should break the loop and close on exit, probably in defer

			case dataChannel <- msg: // write to data if someone reading
				show("task, msg was sent, msg: ", msg)
				msg++
			}
		}
	}(stopChan, dataChan)

	// read messages
	for msg := range dataChan {
		show("main, got msg from data chan, msg: ", msg)
		if msg >= 13 {
			show("main, enough messages, stop ...")
			stopChan <- true
		}
	}
	show("end of program.")
	/*
		2024-04-09T06:32:00.413Z: select_3: produce messages until 'stop' signal
		2024-04-09T06:32:00.413Z: task, loop iteration: int(1);
		2024-04-09T06:32:00.413Z: task, msg was sent, msg: int(11);
		2024-04-09T06:32:00.413Z: task, loop iteration: int(2);
		2024-04-09T06:32:00.413Z: main, got msg from data chan, msg: int(11);
		2024-04-09T06:32:00.413Z: main, got msg from data chan, msg: int(12);
		2024-04-09T06:32:00.413Z: task, msg was sent, msg: int(12);
		2024-04-09T06:32:00.413Z: task, loop iteration: int(3);
		2024-04-09T06:32:00.413Z: task, msg was sent, msg: int(13);
		2024-04-09T06:32:00.413Z: task, loop iteration: int(4);
		2024-04-09T06:32:00.413Z: main, got msg from data chan, msg: int(13);
		2024-04-09T06:32:00.413Z: main, enough messages, stop ...
		2024-04-09T06:32:00.414Z: task, got `stop` command, quitting ...
		2024-04-09T06:32:00.414Z: end of program.
	*/
}

func timeoutDemo() {
	show("timeout: program started ...")

	var longSQLQuery = func() (doneChan chan bool) {
		// create result chan, work 2 sec and than put `true` to result chan

		show("query, creating new 'done' channel ...")
		doneChan = make(chan bool, 1) // buffered rw chan

		// start async work, when done send true to chan
		go func() {
			show("task, started working ...")
			time.Sleep(2 * time.Second)
			show("task, work is done, sending 'done' signal ...")
			doneChan <- true
			show("task, signal 'done' sent.")
		}()

		// don't wait, return immediately
		show("query, async task started (check result chan for `done` signal)")
		return doneChan
	}

	// query takes 2 seconds, so:
	// guard timer: при x=1 выполнится таймаут, при x=3 - выполнится query
	// if x=2 it depends (nondeterministic)
	const x = 2
	timer := time.NewTimer(x * time.Second)

	// one pass, who will be ready first?
	show("main, wait for first signal (timeout or query) ...")
	select {
	case t := <-timer.C: // `C` means 'channel'?
		show("main, got signal from timer: ", t) // after t sec
	case t := <-time.After(5 * time.Second): // n.b. пока не выстрелит - не соберётся сборщиком мусора, read docstring
		show("main, got `time.After` signal: ", t)
	case isDone := <-longSQLQuery():
		show("main, query is done: ", isDone) // after 2 sec of query execution

		// empty and stop timer
		if !timer.Stop() {
			show("main, timer stopped, empty timer chan ...")
			<-timer.C
		}
	}

	show("end of program.")
	/*
		2024-04-09T06:55:18.324Z: timeout: program started ...
		2024-04-09T06:55:18.324Z: main, wait for first signal (timeout or query) ...
		2024-04-09T06:55:18.324Z: query, creating new 'done' channel ...
		2024-04-09T06:55:18.324Z: query, async task started (check result chan for `done` signal)
		2024-04-09T06:55:18.324Z: task, started working ...
		2024-04-09T06:55:20.325Z: task, work is done, sending 'done' signal ...
		2024-04-09T06:55:20.325Z: task, signal 'done' sent.
		2024-04-09T06:55:20.325Z: main, got signal from timer: time.Time(2024-04-09 09:55:20.325651518 +0300 EEST m=+2.001236148);
		2024-04-09T06:55:20.325Z: end of program.
	*/
}

func tickDemo() {
	show("tick: program started ...")

	// stoppable ticker, preferred way to create a ticker
	// n.b. type: `*time.Ticker`
	show("good ticker, demo ...")
	var ticker *time.Ticker = time.NewTicker(100 * time.Millisecond) // ticker chan with backpressure support

	i := 0
	// each second ...
	for tick := range ticker.C {
		i++
		show("got next tick: number; time: ", i, tick)

		// enough
		if i > 2 {
			// stop before exit, or else (resources leak)
			ticker.Stop()
			break
		}
	}
	show("ticks processed: ", i)

	// end of first demo

	// unstoppable ticker, use it if you know what you are doing
	// не может быть остановлен и собран сборщиком мусора
	// n.b. type: `<-chan time.Time`
	show("zombie ticker, demo ...")
	var zombieTickerChan <-chan time.Time = time.Tick(100 * time.Millisecond)
	i = 0
	for tick := range zombieTickerChan {
		i++
		show("got next tick: number; time: ", i, tick)

		// enough
		if i > 2 {
			// n.b. no stopping the ticker here, you can't, it is unstoppable
			break
		}
	}
	show("ticks processed: ", i)

	show("end of program.")
	/*
		2024-04-09T07:00:15.176Z: tick: program started ...
		2024-04-09T07:00:15.176Z: good ticker, demo ...
		2024-04-09T07:00:15.277Z: got next tick: number; time: int(1); time.Time(2024-04-09 10:00:15.276993152 +0300 EEST m=+0.100415578);
		2024-04-09T07:00:15.377Z: got next tick: number; time: int(2); time.Time(2024-04-09 10:00:15.37761897 +0300 EEST m=+0.201041398);
		2024-04-09T07:00:15.477Z: got next tick: number; time: int(3); time.Time(2024-04-09 10:00:15.477240768 +0300 EEST m=+0.300663201);
		2024-04-09T07:00:15.477Z: ticks processed: int(3);
		2024-04-09T07:00:15.477Z: zombie ticker, demo ...
		2024-04-09T07:00:15.578Z: got next tick: number; time: int(1); time.Time(2024-04-09 10:00:15.578048904 +0300 EEST m=+0.401471336);
		2024-04-09T07:00:15.677Z: got next tick: number; time: int(2); time.Time(2024-04-09 10:00:15.67763894 +0300 EEST m=+0.501061372);
		2024-04-09T07:00:15.778Z: got next tick: number; time: int(3); time.Time(2024-04-09 10:00:15.778321762 +0300 EEST m=+0.601744197);
		2024-04-09T07:00:15.778Z: ticks processed: int(3);
		2024-04-09T07:00:15.778Z: end of program.
	*/
}

func afterfunc() {
	show("AfterFunc: program started ...")

	var sayHello = func() {
		show("task, Hello World")
	}

	show("creating after-func timer ...")
	timer := time.AfterFunc(1*time.Second, sayHello)

	userInput("Hit enter when ready (don't want to see 'Hello World'? Hit it before 1-sec-timer) ...")

	timer.Stop() // if you press a button before timer is ready, you won't see a hello-world message

	show("end of program.")
	/*
		2024-04-09T07:02:53.441Z: AfterFunc: program started ...
		2024-04-09T07:02:53.441Z: creating after-func timer ...
		2024-04-09T07:02:53.441Z: Hit enter when ready (don't want to see 'Hello World'? Hit it before 1-sec-timer) ...
		2024-04-09T07:02:54.441Z: task, Hello World
		2024-04-09T07:02:56.474Z: end of program.
	*/
}

func context_cancel() {
	show("context_cancel: program started ...")

	// n.b. non-ascii func name. Why? Just for fun
	var студент = func(ctx context.Context, studentId int, result chan<- int) {
		// work for some time but stop working if ctx.Done is fired

		// random think time
		thinkTime := time.Duration(rand.Intn(100)+10) * time.Millisecond
		show("task, студент x думает y time, x; y: ", studentId, thinkTime)

		// wait for 'ctx.done' or 'think-time-is-out' signal, whichever is first
		select {
		// check 'stop' signal
		case <-ctx.Done():
			show("task, got 'stop' signal, student: ", studentId, ctx.Err())
			return

		// check if work is done
		case <-time.After(thinkTime): // n.b. not GC recoverable if not fired
			show("task, student x is ready, x: ", studentId)
			result <- studentId
		}

		show("task, student x is out, x: ", studentId)
	}

	show("main, create cancellable context ...")
	ctx, cancelFunc := context.WithCancel(context.Background()) // create copy of context with a new Done chan
	resultChan := make(chan int, 3)                             // buffered rw chan

	show("main, start 3 workers ...")
	for id := 0; id < 3; id++ { // start 3 workers
		go студент(ctx, id, resultChan)
	}

	show("main, wait for first answer ...")
	firstReadyStudent := <-resultChan
	show("main, fastest student: ", firstReadyStudent)

	show("main, send cancel message via context to all workers ...")
	cancelFunc()

	show("main, give them some time to exit ...")
	time.Sleep(100 * time.Millisecond)

	show("end of program")
	/*
		2024-04-09T07:31:36.178Z: context_cancel: program started ...
		2024-04-09T07:31:36.178Z: main, create cancellable context ...
		2024-04-09T07:31:36.178Z: main, start 3 workers ...
		2024-04-09T07:31:36.178Z: main, wait for first answer ...
		2024-04-09T07:31:36.178Z: task, студент x думает y time, x; y: int(2); time.Duration(63ms);
		2024-04-09T07:31:36.178Z: task, студент x думает y time, x; y: int(0); time.Duration(10ms);
		2024-04-09T07:31:36.178Z: task, студент x думает y time, x; y: int(1); time.Duration(101ms);
		2024-04-09T07:31:36.188Z: task, student x is ready, x: int(0);
		2024-04-09T07:31:36.188Z: task, student x is out, x: int(0);
		2024-04-09T07:31:36.188Z: main, fastest student: int(0);
		2024-04-09T07:31:36.189Z: main, send cancel message via context to all workers ...
		2024-04-09T07:31:36.189Z: main, give them some time to exit ...
		2024-04-09T07:31:36.189Z: task, got 'stop' signal, student: int(2); *errors.errorString(context canceled);
		2024-04-09T07:31:36.189Z: task, got 'stop' signal, student: int(1); *errors.errorString(context canceled);
		2024-04-09T07:31:36.289Z: end of program
	*/
}

func context_timeout() {
	// see context_cancel() for details
	show("context_timeout: program started ...")

	var doWork = func(ctx context.Context, workerId int, resultChan chan<- int) {
		thinkTime := time.Duration(rand.Intn(100)+10) * time.Millisecond
		show("task, student x think y time, x; y: ", workerId, thinkTime)

		select {
		case <-ctx.Done():
			show("task, got 'stop' signal, student: ", workerId, ctx.Err())
			return

		case <-time.After(thinkTime):
			show("task, student x is ready, x: ", workerId)
			resultChan <- workerId
		}

		show("task, student x is out, x: ", workerId)
	}

	timeLimit := 60 * time.Millisecond
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeLimit)
	defer func() {
		show("main, deferred cleanup")
		cancelFunc() // free resources
	}()

	result := make(chan int, 3) // buffered rw chan

	show("main, starting async work ...")
	for id := 0; id < 3; id++ {
		go doWork(ctx, id, result)
	}

	readyStudentsCount := 0

	show("main, waiting for students (or timeout) ...")
	var i = 0
LOOP: // labeled infinite loop (if context is broken)
	for {
		i += 1
		show("main, wait iteration: ", i)

		select {
		case <-ctx.Done():
			show("main, got 'stop' signal, no more waiting. ", ctx.Err())
			break LOOP

		case readyStudentId := <-result:
			readyStudentsCount++
			show("main, student x is ready, x: ", readyStudentId)
		}
	}
	show("main, results, ready students count: ", readyStudentsCount)

	time.Sleep(100 * time.Millisecond)
	show("end of program.")
	/*
		2024-04-09T07:39:59.928Z: context_timeout: program started ...
		2024-04-09T07:39:59.928Z: main, starting async work ...
		2024-04-09T07:39:59.928Z: main, waiting for students (or timeout) ...
		2024-04-09T07:39:59.928Z: main, wait iteration: int(1);
		2024-04-09T07:39:59.928Z: task, student x think y time, x; y: int(0); time.Duration(62ms);
		2024-04-09T07:39:59.928Z: task, student x think y time, x; y: int(2); time.Duration(17ms);
		2024-04-09T07:39:59.928Z: task, student x think y time, x; y: int(1); time.Duration(41ms);
		2024-04-09T07:39:59.945Z: task, student x is ready, x: int(2);
		2024-04-09T07:39:59.945Z: task, student x is out, x: int(2);
		2024-04-09T07:39:59.945Z: main, student x is ready, x: int(2);
		2024-04-09T07:39:59.945Z: main, wait iteration: int(2);
		2024-04-09T07:39:59.970Z: task, student x is ready, x: int(1);
		2024-04-09T07:39:59.970Z: task, student x is out, x: int(1);
		2024-04-09T07:39:59.970Z: main, student x is ready, x: int(1);
		2024-04-09T07:39:59.970Z: main, wait iteration: int(3);
		2024-04-09T07:39:59.989Z: task, got 'stop' signal, student: int(0); context.deadlineExceededError(context deadline exceeded);
		2024-04-09T07:39:59.989Z: main, got 'stop' signal, no more waiting. context.deadlineExceededError(context deadline exceeded);
		2024-04-09T07:39:59.989Z: main, results, ready students count: int(2);
		2024-04-09T07:40:00.090Z: end of program.
		2024-04-09T07:40:00.090Z: main, deferred cleanup
	*/
}

func async_work() {
	// get 3 pages and for each page get (async) comments
	show("async_work: program started ...")

	var getComments = func() chan string {
		// make a channel and start a work that send a message to that channel.
		// буферизированный канал, in case no one wants to read it
		result := make(chan string, 1)

		// async task
		show("getComments, start async task ...")
		go func(out chan<- string) {
			show("task, started ...")
			// getting page comments imitation
			time.Sleep(200 * time.Millisecond)
			var comments = "32 комментария"
			show("task, write-to-chan comments: ", comments)
			out <- comments
			show("task, done")
			// WTF? writer should close(out)
		}(result)

		show("getComments, done.")
		return result
	}

	// get page data and (async) comments for that page
	var getPage = func(id int) {
		show("getPage, started ...", id)
		commentsChan := getComments() // start comments retrieval, async

		// getting page imitation
		time.Sleep(100 * time.Millisecond) // I'm working, don't touch me ...
		show("getPage, got page data: pageID, data: ", id, "foo")

		// if error ... return // here could be the case where unbuffered channel would cause a deadlock

		// wait for comments
		commentsData := <-commentsChan // got comments
		show("getPage, got comments: pageID, comments:  ", id, commentsData)
	}

	// for 3 pages get comments in async mode
	// show("main, get pages sequentially ...")
	// for i := 0; i < 3; i++ {
	// 	getPage(i)
	// }
	show("main, get pages, concurrent mode ...")
	for i := 0; i < 3; i++ {
		go getPage(i)
	}

	userInput("hit ENTER")
	show("end of program.")
	/*
		2024-04-09T08:02:09.590Z: async_work: program started ...
		2024-04-09T08:02:09.590Z: main, get pages concurrent ...
		2024-04-09T08:02:09.590Z: hit ENTER
		2024-04-09T08:02:09.590Z: getPage, started ...int(0);
		2024-04-09T08:02:09.590Z: getComments, start async task ...
		2024-04-09T08:02:09.590Z: getComments, done.
		2024-04-09T08:02:09.590Z: getPage, started ...int(1);
		2024-04-09T08:02:09.590Z: getPage, started ...int(2);
		2024-04-09T08:02:09.590Z: getComments, start async task ...
		2024-04-09T08:02:09.590Z: getComments, done.
		2024-04-09T08:02:09.590Z: task, started ...
		2024-04-09T08:02:09.590Z: getComments, start async task ...
		2024-04-09T08:02:09.590Z: task, started ...
		2024-04-09T08:02:09.590Z: getComments, done.
		2024-04-09T08:02:09.590Z: task, started ...
		2024-04-09T08:02:09.691Z: getPage, got page data: pageID, data: int(0); string(foo);
		2024-04-09T08:02:09.691Z: getPage, got page data: pageID, data: int(2); string(foo);
		2024-04-09T08:02:09.691Z: getPage, got page data: pageID, data: int(1); string(foo);
		2024-04-09T08:02:09.791Z: task, write-to-chan comments: string(32 комментария);
		2024-04-09T08:02:09.791Z: task, done
		2024-04-09T08:02:09.791Z: getPage, got comments: pageID, comments:  int(1); string(32 комментария);
		2024-04-09T08:02:09.791Z: task, write-to-chan comments: string(32 комментария);
		2024-04-09T08:02:09.791Z: task, done
		2024-04-09T08:02:09.791Z: getPage, got comments: pageID, comments:  int(2); string(32 комментария);
		2024-04-09T08:02:09.791Z: task, write-to-chan comments: string(32 комментария);
		2024-04-09T08:02:09.791Z: task, done
		2024-04-09T08:02:09.791Z: getPage, got comments: pageID, comments:  int(0); string(32 комментария);
		2024-04-09T08:02:16.414Z: end of program.
	*/
}

func workerpool() {
	show("workerpool, program started ...")

	const workersPoolSize = 3

	var messages = []string{
		"Январь", "Февраль", "Март",
		"Апрель", "Май", "Июнь",
		"Июль", "Август", "Сентябрь",
		"Октябрь", "Ноябрь", "Декабрь",
	}

	runtime.GOMAXPROCS(0) // попробуйте: 0 = все доступные; 1 = just one

	var formatWork = func(workerNo int, input string) string {
		return fmt.Sprintln(
			strings.Repeat("  ", workerNo), "█",
			strings.Repeat("  ", workersPoolSize-workerNo), "worker",
			workerNo,
			"recieved", input,
		)
	}

	var printFinishWork = func(workerNo int) {
		fmt.Println(
			strings.Repeat("  ", workerNo), "█",
			strings.Repeat("  ", workersPoolSize-workerNo), "worker",
			workerNo, "finished",
		)
	}

	var startWorker = func(workerNo int, jobsChan <-chan string) {
		show("worker x started, x: ", workerNo)
		// just read from channel and print recieved messages
		for job := range jobsChan {
			fmt.Printf(formatWork(workerNo, job))
			runtime.Gosched() // cooperate
		}
		// channel was closed
		printFinishWork(workerNo)
	}

	var jobsChan = make(chan string, 0) // попробуйте увеличить размер канала

	show("main, start async workers ...")
	for i := 0; i < workersPoolSize; i++ {
		go startWorker(i, jobsChan)
	}
	time.Sleep(100 * time.Millisecond) // give them some startup time

	show("main, sending data to workers ...")
	for _, msg := range messages {
		jobsChan <- msg
	}

	show("main, all work is done, closing jobs queue ...")
	close(jobsChan) // попробуйте закомментировать: worker don't stop, memory leak or deadlock

	time.Sleep(100 * time.Millisecond)
	show("end of program.")
	/*
		2024-04-09T08:15:04.273Z: workerpool, program started ...
		2024-04-09T08:15:04.273Z: main, start async workers ...
		2024-04-09T08:15:04.273Z: worker x started, x: int(2);
		2024-04-09T08:15:04.273Z: worker x started, x: int(0);
		2024-04-09T08:15:04.273Z: worker x started, x: int(1);
		2024-04-09T08:15:04.374Z: main, sending data to workers ...
		   █      worker 1 recieved Март
		     █    worker 2 recieved Январь
		 █        worker 0 recieved Февраль
		   █      worker 1 recieved Апрель
		     █    worker 2 recieved Май
		 █        worker 0 recieved Июнь
		     █    worker 2 recieved Июль
		 █        worker 0 recieved Август
		     █    worker 2 recieved Сентябрь
		 █        worker 0 recieved Октябрь
		     █    worker 2 recieved Ноябрь
		 █        worker 0 recieved Декабрь
		2024-04-09T08:15:04.374Z: main, all work is done, closing jobs queue ...
		   █      worker 1 finished
		 █        worker 0 finished
		     █    worker 2 finished
		2024-04-09T08:15:04.475Z: end of program.
	*/
}

func waitgroupDemo() {
	show("waitgroup, program started ...")
	// we have seen this code, just replace Scanln with WaitGroup

	const (
		iterationsCount = 7
		goroutinesCount = 3
	)

	var formatWork = func(workerNo, iterNo int) string {
		return fmt.Sprintln(
			strings.Repeat("  ", workerNo), "█",
			strings.Repeat("  ", goroutinesCount-workerNo),
			"worker", workerNo,
			"iter", iterNo, strings.Repeat("■", iterNo),
		)
	}

	var startWorker = func(workerNo int, wg *sync.WaitGroup) {
		// on exit decrement counter
		defer wg.Done() // wg.Add(-1)

		show("worker x, starting, x: ", workerNo)
		// working ...
		for j := 0; j < iterationsCount; j++ {
			fmt.Printf(formatWork(workerNo, j))
			time.Sleep(99 * time.Millisecond) // попробуйте убрать этот sleep
		}
		show("worker x, done. x: ", workerNo)
	}

	var main = func() {
		waitGroup := &sync.WaitGroup{} // just a counter wrapper. no copy!

		show("main, start workers ...")
		for i := 0; i < goroutinesCount; i++ {
			waitGroup.Add(1) // wg.Add before starting goroutine
			go startWorker(i, waitGroup)
		}

		show("main, wait for workers ...")
		waitGroup.Wait() // for counter == 0
		show("main, done.")
	}

	main()
	show("end of program.")
	/*
		2024-04-09T08:34:46.764Z: waitgroup, program started ...
		2024-04-09T08:34:46.764Z: main, start workers ...
		2024-04-09T08:34:46.764Z: main, wait for workers ...
		2024-04-09T08:34:46.764Z: worker x, starting, x: int(2);
		     █    worker 2 iter 0
		2024-04-09T08:34:46.764Z: worker x, starting, x: int(1);
		   █      worker 1 iter 0
		2024-04-09T08:34:46.764Z: worker x, starting, x: int(0);
		 █        worker 0 iter 0
		 █        worker 0 iter 1 ■
		     █    worker 2 iter 1 ■
		   █      worker 1 iter 1 ■
		   █      worker 1 iter 2 ■■
		     █    worker 2 iter 2 ■■
		 █        worker 0 iter 2 ■■
		   █      worker 1 iter 3 ■■■
		 █        worker 0 iter 3 ■■■
		     █    worker 2 iter 3 ■■■
		   █      worker 1 iter 4 ■■■■
		     █    worker 2 iter 4 ■■■■
		 █        worker 0 iter 4 ■■■■
		 █        worker 0 iter 5 ■■■■■
		   █      worker 1 iter 5 ■■■■■
		     █    worker 2 iter 5 ■■■■■
		   █      worker 1 iter 6 ■■■■■■
		 █        worker 0 iter 6 ■■■■■■
		     █    worker 2 iter 6 ■■■■■■
		2024-04-09T08:34:47.462Z: worker x, done. x: int(2);
		2024-04-09T08:34:47.462Z: worker x, done. x: int(0);
		2024-04-09T08:34:47.462Z: worker x, done. x: int(1);
		2024-04-09T08:34:47.462Z: main, done.
		2024-04-09T08:34:47.462Z: end of program.
	*/
}

func ratelim() {
	// It's not a rate limit demo, it is a concurrency limit demo
	show("ratelim: program started ...")

	type ZeroWidthMessage struct{}
	var signal ZeroWidthMessage
	const (
		iterationsCount = 5
		goroutinesCount = 4
		// play with this handles
		quotaLimit            = 2
		yieldOnEachNIteration = 3
	)

	var formatWork = func(workerNo, iterNo int) string {
		return fmt.Sprintln(
			strings.Repeat("  ", workerNo), "█",
			strings.Repeat("  ", goroutinesCount-workerNo),
			"worker", workerNo,
			"iter", iterNo, strings.Repeat("■", iterNo),
		)
	}

	var doWork = func(workerNo int, wg *sync.WaitGroup, quotaChan chan ZeroWidthMessage) {
		defer wg.Done() // decrement workers counter on exit

		show("worker x, waiting for quota, x: ", workerNo)
		quotaChan <- signal                // take slot; buffered chan, size of the buffer = number of concurrent tasks
		defer func() { _ = <-quotaChan }() // release slot on exit

		show("worker x, start working, x: ", workerNo)

		for j := 0; j < iterationsCount; j++ {
			fmt.Printf(formatWork(workerNo, j))
			// даём поработать другим горутинам
			time.Sleep(10 * time.Millisecond)
			runtime.Gosched()

			// share resources, don't be greedy: re-take slot
			if (j+1)%yieldOnEachNIteration == 0 {
				show("worker x, releasing slot, x: ", workerNo)
				_ = <-quotaChan
				time.Sleep(10 * time.Millisecond)
				show("worker x, waiting for quota, x: ", workerNo)
				quotaChan <- signal
			}
		}

		show("worker x, stop working, x: ", workerNo)
	}

	var main = func() {
		waitGroup := &sync.WaitGroup{}
		quotaChan := make(chan ZeroWidthMessage, quotaLimit)
		defer close(quotaChan)

		show("main, starting workers ...")
		for i := 0; i < goroutinesCount; i++ {
			waitGroup.Add(1)
			go doWork(i, waitGroup, quotaChan)
		}

		// time.Sleep(time.Millisecond) // why?
		show("main, waiting for workers ...")
		waitGroup.Wait()
		show("main, all workers are finished.")
	}

	main()
	show("end of program.")
	/* // two workers at a time:
	2024-04-09T08:55:23.016Z: ratelim: program started ...
	2024-04-09T08:55:23.016Z: main, starting workers ...
	2024-04-09T08:55:23.016Z: main, waiting for workers ...
	2024-04-09T08:55:23.016Z: worker x, waiting for quota, x: int(3);
	2024-04-09T08:55:23.016Z: worker x, start working, x: int(3);
	       █    worker 3 iter 0
	2024-04-09T08:55:23.016Z: worker x, waiting for quota, x: int(2);
	2024-04-09T08:55:23.016Z: worker x, start working, x: int(2);
	     █      worker 2 iter 0
	       █    worker 3 iter 1 ■
	     █      worker 2 iter 1 ■
	2024-04-09T08:55:23.016Z: worker x, releasing slot, x: int(3);
	2024-04-09T08:55:23.016Z: worker x, waiting for quota, x: int(0);
	2024-04-09T08:55:23.016Z: worker x, start working, x: int(0);
	 █          worker 0 iter 0
	2024-04-09T08:55:23.016Z: worker x, waiting for quota, x: int(1);
	 █          worker 0 iter 1 ■
	2024-04-09T08:55:23.016Z: worker x, releasing slot, x: int(2);
	2024-04-09T08:55:23.016Z: worker x, start working, x: int(1);
	   █        worker 1 iter 0
	   █        worker 1 iter 1 ■
	2024-04-09T08:55:23.016Z: worker x, releasing slot, x: int(1);
	2024-04-09T08:55:23.016Z: worker x, releasing slot, x: int(0);
	2024-04-09T08:55:23.026Z: worker x, waiting for quota, x: int(0);
	 █          worker 0 iter 2 ■■
	2024-04-09T08:55:23.026Z: worker x, waiting for quota, x: int(2);
	2024-04-09T08:55:23.026Z: worker x, waiting for quota, x: int(3);
	 █          worker 0 iter 3 ■■■
	2024-04-09T08:55:23.026Z: worker x, waiting for quota, x: int(1);
	     █      worker 2 iter 2 ■■
	2024-04-09T08:55:23.026Z: worker x, releasing slot, x: int(0);
	       █    worker 3 iter 2 ■■
	     █      worker 2 iter 3 ■■■
	       █    worker 3 iter 3 ■■■
	2024-04-09T08:55:23.026Z: worker x, releasing slot, x: int(2);
	   █        worker 1 iter 2 ■■
	2024-04-09T08:55:23.026Z: worker x, releasing slot, x: int(3);
	   █        worker 1 iter 3 ■■■
	2024-04-09T08:55:23.026Z: worker x, releasing slot, x: int(1);
	2024-04-09T08:55:23.037Z: worker x, waiting for quota, x: int(1);
	   █        worker 1 iter 4 ■■■■
	2024-04-09T08:55:23.037Z: worker x, waiting for quota, x: int(3);
	       █    worker 3 iter 4 ■■■■
	2024-04-09T08:55:23.037Z: worker x, stop working, x: int(1);
	2024-04-09T08:55:23.037Z: worker x, stop working, x: int(3);
	2024-04-09T08:55:23.037Z: worker x, waiting for quota, x: int(0);
	 █          worker 0 iter 4 ■■■■
	2024-04-09T08:55:23.037Z: worker x, stop working, x: int(0);
	2024-04-09T08:55:23.037Z: worker x, waiting for quota, x: int(2);
	     █      worker 2 iter 4 ■■■■
	2024-04-09T08:55:23.037Z: worker x, stop working, x: int(2);
	2024-04-09T08:55:23.037Z: main, all workers are finished.
	2024-04-09T08:55:23.037Z: end of program.
	*/
}

func race_1() {
	// go run -race week02 # program should fail
	show("race_1: program started ...")
	const (
		// sometimes race detected, sometimes not, heisenbug
		workersCount    = 21
		iterationsCount = 27
		mulFactor       = 37
	)

	var counters = map[int]int{}

	wg := &sync.WaitGroup{}
	wg.Add(workersCount)

	show("main, starting workers ...")
	for i := 0; i < workersCount; i++ {
		// start async workers, each updating map values
		go func(workerNo int) {
			show("worker x started, x: ", workerNo)
			defer wg.Done()
			for j := 0; j < iterationsCount; j++ {
				var key = workerNo*mulFactor + j // all keys are different
				// write to global map (not a good idea to begin with)
				counters[key]++
			}
		}(i)
	}

	show("main, wait for workers ...")
	wg.Wait()
	show("main, workers finished, result: ", counters)

	show("end of program.")
	/*
		2024-04-10T08:38:59.801Z: race_1: program started ...
		2024-04-10T08:38:59.801Z: main, starting workers ...
		2024-04-10T08:38:59.801Z: main, wait for workers ...
		2024-04-10T08:38:59.801Z: worker x started, x: int(20);
		2024-04-10T08:38:59.802Z: worker x started, x: int(9);
		2024-04-10T08:38:59.802Z: worker x started, x: int(19);
		2024-04-10T08:38:59.802Z: worker x started, x: int(0);
		2024-04-10T08:38:59.802Z: worker x started, x: int(5);
		fatal error: concurrent map writes

		goroutine 23 [running]:
		main.race_1.func1(0x5)
		        sandbox/week02/main.go:1032 +0xcd
		created by main.race_1 in goroutine 1
		        sandbox/week02/main.go:1026 +0x85
		...
	*/
}

func race_2() {
	show("rece_2: program started ...")

	var counters = map[int]int{}
	countersMtx := &sync.Mutex{} // n.b.: ref
	wg := &sync.WaitGroup{}
	const workersCount = 33

	wg.Add(workersCount)
	for i := 0; i < workersCount; i++ {
		// async
		go func(counters map[int]int, workerNo int, mtx *sync.Mutex) {
			defer wg.Done()
			for j := 0; j < 55; j++ {
				// sequential write
				mtx.Lock()
				counters[workerNo*10+j]++
				mtx.Unlock()
			}
		}(counters, i, countersMtx)
	}

	wg.Wait()
	{ // go run -race ... pill
		countersMtx.Lock()
		show("result: ", counters)
		countersMtx.Unlock()
	}

	show("end of program.")
}

func atomic_1() {
	show("atomic_1: program started ...")

	var totalOperations int32 = 0 // global counter
	var inc = func() {
		totalOperations++ // not synchronized operation
	}

	var main = func() {
		// async increment global counter (not a good idea)
		for i := 0; i < 1000; i++ {
			go inc()
		}

		time.Sleep(333 * time.Millisecond)
		show("expected 1000, got: ", totalOperations)
	}

	main()
	show("end of program.")
	/*
	   2023-11-27T14:18:22.894Z: program started ...
	   2023-11-27T14:18:23.249Z: expected 1000, got: int32(964);
	   2023-11-27T14:18:23.249Z: end of program.
	*/
}

func atomic_2() {
	show("atomic_2: program started ...")

	var totalOperations int32 = 0
	atomicTotalOperations := new(atomic.Int32) // reference
	atomicTotalOperations.Store(0)

	var inc = func() { // synchronized ops
		atomic.AddInt32(&totalOperations, 1) // one way
		atomicTotalOperations.Add(1)         // another way
	}

	var main = func() {
		for i := 0; i < 1000; i++ {
			go inc()
		}

		time.Sleep(333 * time.Millisecond)
		show("expected 1000, got: ", totalOperations, atomicTotalOperations)
	}

	main()
	show("end of program.")
	/*
	   2023-11-27T14:30:13.111Z: program started ...
	   2023-11-27T14:30:13.446Z: expected 1000, got: int32(1000); *atomic.Int32(&{{} 1000});
	   2023-11-27T14:30:13.446Z: end of program.
	*/
}

func demoTemplate() {
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

func userInput(msg string) (res string, err error) {
	show(msg)
	if n, e := fmt.Scanln(&res); n != 1 || e != nil {
		return "", e
	}
	return res, nil
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
