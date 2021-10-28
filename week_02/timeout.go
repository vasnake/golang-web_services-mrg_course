package main

import (
	"fmt"
	"time"
)

func longSQLQuery() chan bool { // sleep 2 sec and put `true` to chan
	ch := make(chan bool, 1)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- true
	}()

	return ch
}

func main() {
	// при 1 выполнится таймаут, при 3 - выполнится операция
	timer := time.NewTimer(3 * time.Second)

	// one pass, who first will be ready?
	select {
	case <-timer.C:
		fmt.Println("timer.C timeout happened") // after 3 sec
	case <-time.After(time.Minute):
		// n.b. пока не выстрелит - не соберётся сборщиком мусора
		fmt.Println("time.After timeout happened") // after 1 min
	case result := <-longSQLQuery():
		// освобождет ресурс
		if !timer.Stop() {
			<-timer.C
		}
		fmt.Println("operation result:", result) // after 2 sec
	}
}
