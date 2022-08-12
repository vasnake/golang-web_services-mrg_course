package main

import (
	"fmt"
)

func main() {
	ch1 := make(chan int, 1)
	// ch1 <- 1 // if uncomment, case 1 will be selected

	ch2 := make(chan int)

	select {
	case val := <-ch1:
		fmt.Println("got ch1 val", val)
	case ch2 <- 1:
		fmt.Println("put `1` to ch2")
	default: // w/o default you will get deadlock
		fmt.Println("have no ready channels")
	}

}
