package main

import (
	"fmt"
	"time"
)

func getComments() chan string {
	// make a channel and stark a work that send a message to that channel.

	// надо использовать буферизированный канал, in case no one wants to read it
	result := make(chan string, 1)

	go func(out chan<- string) {
		time.Sleep(2 * time.Second)
		fmt.Println("async operation ready, return comments")
		out <- "32 комментария"
	}(result)

	return result
}

func getPage() {
	resultCh := getComments() // start comments retrieval, async

	time.Sleep(1 * time.Second) // I'm working, don't touch me ...
	fmt.Println("got related articles")

	// return // here could be the case where unbuffered channel would cause a deadlock

	commentsData := <-resultCh // got comments
	fmt.Println("main goroutine:", commentsData)
}

func main() {
	// for 3 pages get comments in async mode
	for i := 0; i < 3; i++ {
		getPage()
	}
}
