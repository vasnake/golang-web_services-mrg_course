package main

import (
	"fmt"
	"time"
)

func main() {

	ch1 := make(chan int) // unbuffered channel, all put values must be taken from chan, or else: deadlock
	//ch1 := make(chan int, 1) // buffered channel, buffer size = 1 allow to put 1 unread message to chan.

	go func(in chan int) {
		fmt.Println("anon, Waiting for channel ...")
		v := <-in // read from channel, transfer control over data to me
		fmt.Println("anon, Value from channel:", v)
		fmt.Println("anon, done.")
	}(ch1)

	fmt.Println("parent, go to sleep ...")
	time.Sleep(3 * time.Second)

	fmt.Println("parent, Writing to channel ...")
	ch1 <- 42 // write to channel, transfer control over data to anon. routine
	// ch1 <- 100500 // if unbuffered: deadlock in program, fatal; or goroutine/mem leak

	fmt.Println("parent, done.")

	fmt.Scanln()
}
