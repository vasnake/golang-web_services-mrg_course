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
	show("main, started ...")

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

	show("main, done.")
	//fmt.Scanln()
	/*
	   2023-11-24T10:05:19.583Z: main, started ...
	   2023-11-24T10:05:19.583Z: generator, before writing: int(0);
	   2023-11-24T10:05:19.583Z: generator, after writing: int(0);
	   2023-11-24T10:05:19.583Z: generator, before writing: int(1);
	   2023-11-24T10:05:19.583Z: generator, after writing: int(1);
	   2023-11-24T10:05:19.583Z: generator, before writing: int(2);
	   2023-11-24T10:05:19.583Z: main, readed from chan: int(0);
	   2023-11-24T10:05:19.583Z: main, readed from chan: int(1);
	   2023-11-24T10:05:19.583Z: main, readed from chan: int(2);
	   2023-11-24T10:05:19.583Z: generator, after writing: int(2);
	   2023-11-24T10:05:19.583Z: generator, before writing: int(3);
	   2023-11-24T10:05:19.583Z: generator, after writing: int(3);
	   2023-11-24T10:05:19.583Z: generator, before writing: int(4);
	   2023-11-24T10:05:19.583Z: generator, after writing: int(4);
	   2023-11-24T10:05:19.583Z: generator, before writing: int(5);
	   2023-11-24T10:05:19.583Z: main, readed from chan: int(3);
	   2023-11-24T10:05:19.583Z: main, readed from chan: int(4);
	   2023-11-24T10:05:19.583Z: main, readed from chan: int(5);
	   2023-11-24T10:05:19.583Z: generator, after writing: int(5);
	   2023-11-24T10:05:19.583Z: generator, finished.
	   2023-11-24T10:05:19.583Z: main, done.
	*/
}

func select_1() {
	ch1 := make(chan int, 1) // buffered rw channel
	// ch1 <- 1                 // if uncomment, case 1 will be selected

	ch2 := make(chan int) // un-buffered rw channel

	var x = 42
	select {
	case x := <-ch1: // read from chan if it has message
		show("got x from ch1, x: ", x)
	case ch2 <- x: // if upper cases are no-go: write to chan if it can receive a message
		show("put x to ch2, x: ", x)
	default: // w/o default you will get deadlock
		show("have no ready channels")
	}
	// `select` is not wrapped to loop, evaluate only once, than exit
}

func select_2() {
	show("program started ...")
	ch1 := make(chan int, 2) // buffered rw chan, fill it up
	ch1 <- 1
	ch1 <- 2

	ch2 := make(chan int, 2) // buffered rw chan, put just one message to it
	ch2 <- 3

	var i = 0
	var x int
LOOP: // labeled loop
	for {
		i += 1
		show("loop iteration: ", i)

		select { // check channels from top to bottom
		case x = <-ch1:
			show("got x from ch1, x: ", x)
		case x = <-ch2:
			show("got x from ch2, x: ", x)
		default: // if no channels are ready
			show("all channels are empty, bye.")
			break LOOP // exit from labeled loop
		}
		show("end of loop iteration: ", i)
	}
	show("end of program.")
	/*
	   2023-11-24T10:33:00.639Z: program started ...
	   2023-11-24T10:33:00.639Z: loop iteration: int(1);
	   2023-11-24T10:33:00.639Z: got x from ch1, x: int(1);
	   2023-11-24T10:33:00.639Z: end of loop iteration: int(1);
	   2023-11-24T10:33:00.639Z: loop iteration: int(2);
	   2023-11-24T10:33:00.639Z: got x from ch1, x: int(2);
	   2023-11-24T10:33:00.639Z: end of loop iteration: int(2);
	   2023-11-24T10:33:00.639Z: loop iteration: int(3);
	   2023-11-24T10:33:00.639Z: got x from ch2, x: int(3);
	   2023-11-24T10:33:00.639Z: end of loop iteration: int(3);
	   2023-11-24T10:33:00.639Z: loop iteration: int(4);
	   2023-11-24T10:33:00.639Z: all channels are empty, bye.
	   2023-11-24T10:33:00.639Z: end of program.
	*/
}

func select_3() {
	show("program started ...")

	stopChan := make(chan bool) // commands chan, unbuffered rw
	dataChan := make(chan int)  // messages chan, unbuffered rw

	// run async messages sender
	go func(stopChannel <-chan bool, dataChannel chan<- int) {
		i, msg := 0, 0
		for { // loop until stop command recieved
			i += 1
			show("task, loop iteration: ", i)

			select {
			// N.B. no default branch, you have to be sure that someone give you the stop signal
			// Think about protocol here ...

			// check commands channel first, data chan here is for writing, so life is good (not so if data chan were for reading)
			case <-stopChannel: // you have to stop now
				show("task, got `stop` command, quitting ...")
				close(dataChannel) // this will help reader to stop iterations in `range chan``
				return

			case dataChannel <- msg: // write to data if someone reading
				show("task, msg was sent, msg: ", msg)
				msg++
			}
		}
	}(stopChan, dataChan)

	// read messages
	for msg := range dataChan {
		show("main, got msg from data chan, msg: ", msg)
		if msg > 2 {
			show("main, enough messages, stop ...")
			stopChan <- true
		}
	}
	show("end of program.")
	/*
	   2023-11-24T10:54:23.069Z: program started ...
	   2023-11-24T10:54:23.069Z: task, loop iteration: int(1);
	   2023-11-24T10:54:23.069Z: task, msg was sent, msg: int(0);
	   2023-11-24T10:54:23.069Z: task, loop iteration: int(2);
	   2023-11-24T10:54:23.069Z: main, got msg from data chan, msg: int(0);
	   2023-11-24T10:54:23.069Z: main, got msg from data chan, msg: int(1);
	   2023-11-24T10:54:23.069Z: task, msg was sent, msg: int(1);
	   2023-11-24T10:54:23.069Z: task, loop iteration: int(3);
	   2023-11-24T10:54:23.069Z: task, msg was sent, msg: int(2);
	   2023-11-24T10:54:23.069Z: task, loop iteration: int(4);
	   2023-11-24T10:54:23.069Z: main, got msg from data chan, msg: int(2);
	   2023-11-24T10:54:23.069Z: main, got msg from data chan, msg: int(3);
	   2023-11-24T10:54:23.069Z: main, enough messages, stop ...
	   2023-11-24T10:54:23.069Z: task, msg was sent, msg: int(3);
	   2023-11-24T10:54:23.069Z: task, loop iteration: int(5);
	   2023-11-24T10:54:23.069Z: task, got `stop` command, quitting ...
	   2023-11-24T10:54:23.069Z: end of program.
	*/
}

func timeoutDemo() {
	show("program started ...")

	var longSQLQuery = func() (doneChan chan bool) {
		// create result chan, work 2 sec and than put `true` to result chan

		show("query, creating new 'done' channel ...")
		doneChan = make(chan bool, 1) // buffered rw chan

		// async work, when done send true to chan
		go func() {
			show("task, started working ...")
			time.Sleep(2 * time.Second)
			show("task, work is done, sending 'done' signal ...")
			doneChan <- true
			show("task, signal 'done' sent.")
		}()

		// don't wait, return immediately
		show("query, async task started (check result chan for `done` signal), exit.")
		return doneChan
	}

	// query takes 2 seconds, so:
	// guard timer: при x=1 выполнится таймаут, при x=3 - выполнится query
	// if x=2 it depends (nondeterministic)
	const x = 2
	timer := time.NewTimer(x * time.Second)

	// one pass, who will be ready first?
	show("main, wait for first signal ...")
	select {
	case t := <-timer.C: // `C` means 'channel'?
		show("main, got signal from timer: ", t) // after x sec
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
	   2023-11-24T11:40:30.947Z: program started ...
	   2023-11-24T11:40:30.947Z: main, wait for first signal ...
	   2023-11-24T11:40:30.947Z: query, creating new 'done' channel ...
	   2023-11-24T11:40:30.947Z: query, async task started (check result chan for `done` signal), exit.
	   2023-11-24T11:40:30.947Z: task, started working ...
	   2023-11-24T11:40:33.001Z: task, work is done, sending 'done' signal ...
	   2023-11-24T11:40:33.001Z: task, signal 'done' sent.
	   2023-11-24T11:40:33.001Z: main, query is done: bool(true);
	   2023-11-24T11:40:33.001Z: main, timer stopped, empty timer chan ...
	   2023-11-24T11:40:33.001Z: end of program.
	*/
}

func tickDemo() {
	show("program started ...")

	// stoppable ticker, preferrable way to create a ticker
	// n.b. type: `*time.Ticker`
	show("good ticker demo ...")
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
	show("zombie ticker demo ...")
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
	   2023-11-24T12:01:56.867Z: program started ...
	   2023-11-24T12:01:56.867Z: good ticker demo ...
	   2023-11-24T12:01:56.968Z: got next tick: number; time: int(1); time.Time(2023-11-24 14:01:56.968544395 +0200 EET m=+0.101224668);
	   2023-11-24T12:01:57.073Z: got next tick: number; time: int(2); time.Time(2023-11-24 14:01:57.073396686 +0200 EET m=+0.206076961);
	   2023-11-24T12:01:57.174Z: got next tick: number; time: int(3); time.Time(2023-11-24 14:01:57.174613062 +0200 EET m=+0.307293345);
	   2023-11-24T12:01:57.174Z: ticks processed: int(3);
	   2023-11-24T12:01:57.174Z: zombie ticker demo ...
	   2023-11-24T12:01:57.276Z: got next tick: number; time: int(1); time.Time(2023-11-24 14:01:57.276123575 +0200 EET m=+0.408803839);
	   2023-11-24T12:01:57.382Z: got next tick: number; time: int(2); time.Time(2023-11-24 14:01:57.382690281 +0200 EET m=+0.515370576);
	   2023-11-24T12:01:57.482Z: got next tick: number; time: int(3); time.Time(2023-11-24 14:01:57.482853534 +0200 EET m=+0.615533785);
	   2023-11-24T12:01:57.482Z: ticks processed: int(3);
	   2023-11-24T12:01:57.482Z: end of program.
	*/
}

func afterfunc() {
	show("program started ...")

	var sayHello = func() {
		show("task, Hello World")
	}

	show("creating after-func timer ...")
	timer := time.AfterFunc(1*time.Second, sayHello)

	show("Hit enter when ready (don't want to see 'Hello World'? Hit it before 1-sec-timer) ...")
	fmt.Scanln()

	timer.Stop() // if you press a button before timer is ready, you won't see a hello-world message

	show("end of program.")
	/*
	   2023-11-24T12:32:39.612Z: program started ...
	   2023-11-24T12:32:39.612Z: creating after-func timer ...
	   2023-11-24T12:32:39.612Z: Hit enter when ready (don't want to see 'Hello World'? Hit it before 1-sec-timer) ...
	   2023-11-24T12:32:40.687Z: task, Hello World
	   2023-11-24T12:32:43.934Z: end of program.
	*/
}

func context_cancel() {
	show("program started ...")

	// n.b. non-ascii func name
	var студент = func(ctx context.Context, studentId int, result chan<- int) {
		// work for some time but stop working if ctx.Done is fired

		// random think time
		thinkTime := time.Duration(rand.Intn(100)+10) * time.Millisecond
		show("task, студент x думает y time, x; y: ", studentId, thinkTime)

		// wait for 'ctx.done' or 'think-time-is-out' signal, whichever are first
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
	resultChan := make(chan int, 1)                             // buffered rw chan

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
		2023-11-24T13:25:05.910Z: program started ...
		2023-11-24T13:25:05.910Z: main, create cancellable context ...
		2023-11-24T13:25:05.910Z: main, start 3 workers ...
		2023-11-24T13:25:05.910Z: main, wait for first answer ...
		2023-11-24T13:25:05.910Z: task, студент x думает y time, x; y: int(2); time.Duration(77ms);
		2023-11-24T13:25:05.910Z: task, студент x думает y time, x; y: int(0); time.Duration(91ms);
		2023-11-24T13:25:05.910Z: task, студент x думает y time, x; y: int(1); time.Duration(89ms);
		2023-11-24T13:25:05.989Z: task, student x is ready, x: int(2);
		2023-11-24T13:25:05.989Z: task, student x is out, x: int(2);
		2023-11-24T13:25:05.989Z: main, fastest student: int(2);
		2023-11-24T13:25:05.989Z: main, send cancel message via context to all workers ...
		2023-11-24T13:25:05.989Z: main, give them some time to exit ...
		2023-11-24T13:25:05.989Z: task, got 'stop' signal, student: int(0); *errors.errorString(context canceled);
		2023-11-24T13:25:05.989Z: task, got 'stop' signal, student: int(1); *errors.errorString(context canceled);
		2023-11-24T13:25:06.096Z: end of program
	*/
}

func context_timeout() {
	// see context_cancel() for details
	show("program started ...")

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
	defer cancelFunc() // free resources

	result := make(chan int, 1) // buffered rw chan

	show("main, starting async work ...")
	for id := 0; id < 3; id++ {
		go doWork(ctx, id, result)
	}

	readyStudentsCount := 0

	show("main, waiting for students ...")
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
	   2023-11-24T13:52:58.417Z: program started ...
	   2023-11-24T13:52:58.417Z: main, starting async work ...
	   2023-11-24T13:52:58.417Z: main, waiting for students ...
	   2023-11-24T13:52:58.417Z: main, wait iteration: int(1);
	   2023-11-24T13:52:58.417Z: task, student x think y time, x; y: int(2); time.Duration(66ms);
	   2023-11-24T13:52:58.417Z: task, student x think y time, x; y: int(1); time.Duration(48ms);
	   2023-11-24T13:52:58.417Z: task, student x think y time, x; y: int(0); time.Duration(74ms);
	   2023-11-24T13:52:58.466Z: task, student x is ready, x: int(1);
	   2023-11-24T13:52:58.466Z: task, student x is out, x: int(1);
	   2023-11-24T13:52:58.466Z: main, student x is ready, x: int(1);
	   2023-11-24T13:52:58.466Z: main, wait iteration: int(2);
	   2023-11-24T13:52:58.477Z: task, got 'stop' signal, student: int(2); context.deadlineExceededError(context deadline exceeded);
	   2023-11-24T13:52:58.477Z: task, got 'stop' signal, student: int(0); context.deadlineExceededError(context deadline exceeded);
	   2023-11-24T13:52:58.477Z: main, got 'stop' signal, no more waiting. context.deadlineExceededError(context deadline exceeded);
	   2023-11-24T13:52:58.477Z: main, results, ready students count: int(1);
	   2023-11-24T13:52:58.585Z: end of program.
	*/
}

func async_work() {
	// get 3 pages and for each page (async) comments
	show("program started ...")

	var getComments = func() chan string {
		// make a channel and start a work that send a message to that channel.

		// буферизированный канал, in case no one wants to read it
		result := make(chan string, 1)

		// async task
		go func(out chan<- string) {
			// getting page comments imitation
			time.Sleep(200 * time.Millisecond)
			var comments = "32 комментария"
			show("getComments, async operation finished, return comments: ", comments)
			out <- comments
		}(result)

		return result
	}

	// get page data and (async) comments for that page
	var getPage = func() {
		commentsChan := getComments() // start comments retrieval, async

		// getting page imitation
		time.Sleep(100 * time.Millisecond) // I'm working, don't touch me ...
		show("getPage, got page data: ", "foo")

		// return // here could be the case where unbuffered channel would cause a deadlock

		// wait for getComments finish
		commentsData := <-commentsChan // got comments
		show("getPage, page comments: ", commentsData)
	}

	// for 3 pages get comments in async mode
	for i := 0; i < 3; i++ {
		getPage()
	}

	show("end of program.")
}

func workerpool() {
	show("program started ...")

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
		// start nothing, just read from channel and print recieved messages
		for job := range jobsChan {
			fmt.Printf(formatWork(workerNo, job))
			runtime.Gosched() // cooperate
		}
		// channel closed
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
	   2023-11-26T16:04:09.230Z: program started ...
	   2023-11-26T16:04:09.230Z: main, start async workers ...
	   2023-11-26T16:04:09.230Z: worker x started, x: int(2);
	   2023-11-26T16:04:09.230Z: worker x started, x: int(0);
	   2023-11-26T16:04:09.230Z: worker x started, x: int(1);
	   2023-11-26T16:04:09.331Z: main, sending data to workers ...
	      █      worker 1 recieved Март
	        █    worker 2 recieved Январь
	    █        worker 0 recieved Февраль
	      █      worker 1 recieved Апрель
	        █    worker 2 recieved Май
	    █        worker 0 recieved Июнь
	      █      worker 1 recieved Июль
	        █    worker 2 recieved Август
	    █        worker 0 recieved Сентябрь
	      █      worker 1 recieved Октябрь
	        █    worker 2 recieved Ноябрь
	    █        worker 0 recieved Декабрь
	   2023-11-26T16:04:09.331Z: main, all work is done, closing jobs queue ...
	        █    worker 2 finished
	    █        worker 0 finished
	      █      worker 1 finished
	   2023-11-26T16:04:09.431Z: end of program.
	*/
}

func waitgroupDemo() {
	show("program started ...")

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
		show("worker x, starting, x: ", workerNo)

		// on exit decrement counter
		defer wg.Done() // wg.Add(-1)

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
			// wg.Add before starting goroutine
			waitGroup.Add(1)
			go startWorker(i, waitGroup)
		}

		// time.Sleep(time.Millisecond)
		show("main, wait for workers ...")
		waitGroup.Wait() // for counter == 0
		show("main, done.")
	}

	main()
	show("end of program.")
	/*
	   2023-11-27T10:50:49.404Z: program started ...
	   2023-11-27T10:50:49.404Z: main, start workers ...
	   2023-11-27T10:50:49.404Z: main, wait for workers ...
	   2023-11-27T10:50:49.404Z: worker x, starting, x: int(2);
	        █    worker 2 iter 0
	   2023-11-27T10:50:49.404Z: worker x, starting, x: int(0);
	    █        worker 0 iter 0
	   2023-11-27T10:50:49.405Z: worker x, starting, x: int(1);
	      █      worker 1 iter 0
	        █    worker 2 iter 1 ■
	    █        worker 0 iter 1 ■
	      █      worker 1 iter 1 ■
	      █      worker 1 iter 2 ■■
	    █        worker 0 iter 2 ■■
	        █    worker 2 iter 2 ■■
	        █    worker 2 iter 3 ■■■
	    █        worker 0 iter 3 ■■■
	      █      worker 1 iter 3 ■■■
	      █      worker 1 iter 4 ■■■■
	    █        worker 0 iter 4 ■■■■
	        █    worker 2 iter 4 ■■■■
	      █      worker 1 iter 5 ■■■■■
	        █    worker 2 iter 5 ■■■■■
	    █        worker 0 iter 5 ■■■■■
	      █      worker 1 iter 6 ■■■■■■
	    █        worker 0 iter 6 ■■■■■■
	        █    worker 2 iter 6 ■■■■■■
	   2023-11-27T10:50:50.103Z: worker x, done. x: int(1);
	   2023-11-27T10:50:50.103Z: worker x, done. x: int(2);
	   2023-11-27T10:50:50.103Z: worker x, done. x: int(0);
	   2023-11-27T10:50:50.103Z: main, done.
	   2023-11-27T10:50:50.103Z: end of program.
	*/
}

func ratelim() {
	// It's not a rate limit demo, it is a concurrency limit demo
	show("program started ...")

	type ZeroWidthMessage struct{}
	var signal ZeroWidthMessage
	const (
		iterationsCount = 5
		goroutinesCount = 4
		quotaLimit      = 2
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
		defer wg.Done() // decrement workers counter

		show("worker x, waiting for quota, x: ", workerNo)
		quotaChan <- signal                // take slot; buffered chan, size of the buffer = number of concurrent tasks
		defer func() { _ = <-quotaChan }() // free slot
		show("worker x, start working, x: ", workerNo)

		for j := 0; j < iterationsCount; j++ {
			fmt.Printf(formatWork(workerNo, j))

			// share resources, don't be greedy
			if (j+1)%3 == 0 {
				show("worker x, releasing slot, x: ", workerNo)
				_ = <-quotaChan // ratelim.go, возвращаем слот
				// runtime.Gosched() // даём поработать другим горутинам
				time.Sleep(10 * time.Millisecond)
				show("worker x, waiting for quota, x: ", workerNo)
				quotaChan <- signal // ratelim.go, wait for quota
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
	/*
	   2023-11-27T11:55:35.523Z: program started ...
	   2023-11-27T11:55:35.523Z: main, starting workers ...
	   2023-11-27T11:55:35.523Z: main, waiting for workers ...
	   2023-11-27T11:55:35.523Z: worker x, waiting for quota, x: int(3);
	   2023-11-27T11:55:35.523Z: worker x, start working, x: int(3);
	          █    worker 3 iter 0
	          █    worker 3 iter 1 ■
	          █    worker 3 iter 2 ■■
	   2023-11-27T11:55:35.523Z: worker x, releasing slot, x: int(3);
	   2023-11-27T11:55:35.523Z: worker x, waiting for quota, x: int(0);
	   2023-11-27T11:55:35.523Z: worker x, start working, x: int(0);
	    █          worker 0 iter 0
	    █          worker 0 iter 1 ■
	    █          worker 0 iter 2 ■■
	   2023-11-27T11:55:35.523Z: worker x, releasing slot, x: int(0);
	   2023-11-27T11:55:35.523Z: worker x, waiting for quota, x: int(1);
	   2023-11-27T11:55:35.523Z: worker x, start working, x: int(1);
	      █        worker 1 iter 0
	      █        worker 1 iter 1 ■
	      █        worker 1 iter 2 ■■
	   2023-11-27T11:55:35.523Z: worker x, releasing slot, x: int(1);
	   2023-11-27T11:55:35.523Z: worker x, waiting for quota, x: int(2);
	   2023-11-27T11:55:35.523Z: worker x, start working, x: int(2);
	        █      worker 2 iter 0
	        █      worker 2 iter 1 ■
	        █      worker 2 iter 2 ■■
	   2023-11-27T11:55:35.523Z: worker x, releasing slot, x: int(2);
	   2023-11-27T11:55:35.533Z: worker x, waiting for quota, x: int(3);
	          █    worker 3 iter 3 ■■■
	          █    worker 3 iter 4 ■■■■
	   2023-11-27T11:55:35.533Z: worker x, stop working, x: int(3);
	   2023-11-27T11:55:35.533Z: worker x, waiting for quota, x: int(2);
	   2023-11-27T11:55:35.533Z: worker x, waiting for quota, x: int(1);
	      █        worker 1 iter 3 ■■■
	        █      worker 2 iter 3 ■■■
	        █      worker 2 iter 4 ■■■■
	   2023-11-27T11:55:35.533Z: worker x, stop working, x: int(2);
	      █        worker 1 iter 4 ■■■■
	   2023-11-27T11:55:35.533Z: worker x, waiting for quota, x: int(0);
	   2023-11-27T11:55:35.533Z: worker x, stop working, x: int(1);
	    █          worker 0 iter 3 ■■■
	    █          worker 0 iter 4 ■■■■
	   2023-11-27T11:55:35.534Z: worker x, stop working, x: int(0);
	   2023-11-27T11:55:35.534Z: main, all workers are finished.
	   2023-11-27T11:55:35.534Z: end of program.
	*/
}

func race_1() {
	// go run -race week02 # program should fail
	show("program started ...")
	const workersCount = 33

	var counters = map[int]int{}

	wg := &sync.WaitGroup{}
	wg.Add(workersCount)

	show("main, starting workers ...")
	for i := 0; i < workersCount; i++ {
		// start async workers, each updating map values
		go func(workerNo int) {
			show("worker x started, x: ", workerNo)
			defer wg.Done()
			for j := 0; j < 55; j++ {
				var key = workerNo*10 + j
				// write to global map (not a good idea)
				counters[key]++
			}
		}(i)
	}

	show("main, wait for workers ...")
	wg.Wait()
	show("main, workers finished, result: ", counters)

	show("end of program.")
	/*
	   2023-11-27T13:51:34.536Z: program started ...
	   2023-11-27T13:51:34.536Z: main, starting workers ...
	   2023-11-27T13:51:34.536Z: main, wait for workers ...
	   2023-11-27T13:51:34.536Z: worker x started, x: int(32);
	   2023-11-27T13:51:34.536Z: worker x started, x: int(15);
	   2023-11-27T13:51:34.536Z: worker x started, x: int(26);
	   fatal error: concurrent map writes
	*/
}

func race_2() {
	show("program started ...")

	var counters = map[int]int{}
	countersMutex := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	const workersCount = 33

	wg.Add(workersCount)
	for i := 0; i < workersCount; i++ {
		// async
		go func(counters map[int]int, workerNo int, mutex *sync.Mutex) {
			defer wg.Done()
			for j := 0; j < 55; j++ {
				// sequential write
				mutex.Lock()
				counters[workerNo*10+j]++
				mutex.Unlock()
			}
		}(counters, i, countersMutex)
	}

	wg.Wait()
	show("result: ", counters)

	show("end of program.")
}

func atomic_1() {
	show("program started ...")

	var totalOperations int32 = 0
	var inc = func() {
		totalOperations++
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
	show("program started ...")

	var totalOperations int32
	atomicTotalOperations := new(atomic.Int32)

	var inc = func() {
		atomic.AddInt32(&totalOperations, 1)
		atomicTotalOperations.Add(1)
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

func main() {
	// goroutinesDemo()
	// chan_1()
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

// ts return current timestamp in RFC3339 with milliseconds
func ts() string {
	/*
		https://pkg.go.dev/time#pkg-constants
		https://stackoverflow.com/questions/35479041/how-to-convert-iso-8601-time-in-golang
	*/
	const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
	return time.Now().UTC().Format(RFC3339Milli)
}
