package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	iterationsNum = 6
	goroutinesNum = 5
	quotaLimit    = 2
)

func startWorker(in int, wg *sync.WaitGroup, quotaCh chan struct{}) {
	quotaCh <- struct{}{} // берём свободный слот

	defer wg.Done() // on_exit, report to waitgroup

	for j := 0; j < iterationsNum; j++ {
		fmt.Printf(formatWork(in, j))

		if j%2 == 0 { // make a fair play
			<-quotaCh             // ratelim.go, возвращаем слот
			quotaCh <- struct{}{} // ratelim.go, берём слот
		}

		runtime.Gosched() // даём поработать другим горутинам
	}

	<-quotaCh // возвращаем слот // why not in on_exit (defer)?
}

func main() {
	wg := &sync.WaitGroup{}
	quotaCh := make(chan struct{}, quotaLimit)

	for i := 0; i < goroutinesNum; i++ {
		wg.Add(1)
		go startWorker(i, wg, quotaCh)
	}

	time.Sleep(time.Millisecond) // why?
	wg.Wait()
}

func formatWork(in, j int) string {
	return fmt.Sprintln(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"th", in,
		"iter", j, strings.Repeat("■", j))
}
