package main

import (
	"fmt"
	"strings"
	// "runtime"
)

const (
	iterationsNum = 7
	goroutinesNum = 5
)

func doSomeWork(goNum int) {
	for j := 0; j < iterationsNum; j++ {
		// could be blocked by inner infinite loop
		fmt.Printf(formatWork(goNum, j))

		// return to scheduler, allow to switch to other tasks; break my own time-slice
		// runtime.Gosched()
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

func formatWork(goNum, iterNum int) string {
	return fmt.Sprintln(
		strings.Repeat("  ", goNum),
		"█",
		strings.Repeat("  ", goroutinesNum-goNum),
		"goroutine:",
		goNum,
		"iter:",
		iterNum,
		strings.Repeat("■", iterNum))
}

/*
package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

const (
	iterationsNum = 6
	goroutinesNum = 6
)

func doWork(th int) {
	for j := 0; j < iterationsNum; j++ {
		fmt.Printf(formatWork(th, j))
		// time.Sleep(time.Millisecond)
	}
}

func main() {
	for i := 0; i < goroutinesNum; i++ {
		go doWork(i)
	}
	fmt.Scanln()
}

func formatWork(in, j int) string {
	return fmt.Sprintln(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"th", in,
		"iter", j, strings.Repeat("■", j))
}

func imports() {
	fmt.Println(time.Millisecond, runtime.NumCPU())
}

*/
