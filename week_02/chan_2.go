package main

import (
	"fmt"
)

func main() {
	in := make(chan int, 1) // rw channel

	go func(out chan<- int) { // wo channel
		for i := 0; i <= 10; i++ {
			fmt.Println("before writing", i)
			out <- i
			fmt.Println("after writing", i)
		}
		close(out) // signal for reader to stop iterations
		//out <- 12 // panic: send on closed channel, can' write to closed chan
		fmt.Println("generator finished")
	}(in)

	for i := range in { // read from channel
		// i, isOpened := <-in
		// if !isOpened {
		// 	break
		// }
		fmt.Println("\tgot", i)
	}

	//fmt.Scanln()
}
