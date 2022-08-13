package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	iterationsNum = 7
	goroutinesNum = 5
)

func startWorker(in int, wg *sync.WaitGroup) {
	// work in goroutine

	// on exit pong waitgroup
	defer wg.Done() // wg.Substract(1)

	// working ...
	for j := 0; j < iterationsNum; j++ {
		fmt.Printf(formatWork(in, j))
		time.Sleep(time.Millisecond) // попробуйте убрать этот sleep
	}
}

func main() {
	wg := &sync.WaitGroup{} // wait_2.go инициализируем группу

	for i := 0; i < goroutinesNum; i++ {
		// wg.Add надо вызывать в той горутине, которая порождает воркеров
		// иначе другая горутина может не успеть запуститься и выполнится Wait
		wg.Add(1) // ping waitgroup
		go startWorker(i, wg)
	}

	time.Sleep(time.Millisecond)
	wg.Wait() // for waitgroup == 0

	fmt.Println(11111)
}

func formatWork(in, j int) string {
	return fmt.Sprintln(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"th", in,
		"iter", j, strings.Repeat("■", j))
}
