package main

import (
	"fmt"
	"strings"
)

const (
	iterationsNum = 7
	goroutinesNum = 5
)

func doSomeWork(in int) {
	for j := 0; j < iterationsNum; j++ {
		fmt.Printf(formatWork(in, j))
		// return to scheduler, allow to switch to other tasks; break own time-slice
		//runtime.Gosched()
		//time.Sleep(time.Millisecond)
	}
}

func main() {
	for i := 0; i < goroutinesNum; i++ {

		// sequential calls
		//doSomeWork(i)

		// concurrent calls, can't return result value here
		go doSomeWork(i)
	}

	// don't exit before all threads are finished
	fmt.Scanln()
}

func formatWork(in, j int) string {
	return fmt.Sprintln(
		strings.Repeat("  ", in),
		"█",
		strings.Repeat("  ", goroutinesNum-in),
		"th",
		in,
		"iter",
		j,
		strings.Repeat("■", j))
}
