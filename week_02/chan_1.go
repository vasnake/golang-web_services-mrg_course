package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)

	go func(in chan int) {
		fmt.Println("Waiting for channel ...")
		v := <-in // read from channel
		fmt.Println("Value from channel:", v)
		fmt.Println("Next line ...")
	}(ch1)

	fmt.Println("Writing to channel ...")
	time.Sleep(3 * time.Second)
	ch1 <- 42 // write to channel
	fmt.Println("Written to channel")

	fmt.Scanln()
}
