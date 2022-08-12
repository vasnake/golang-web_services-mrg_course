package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func студент(ctx context.Context, workerNum int, out chan<- int) {
	waitTime := time.Duration(rand.Intn(100)+10) * time.Millisecond // random think time
	fmt.Println(workerNum, "студент думает", waitTime)

	select {
	case <-ctx.Done():
		return

	case <-time.After(waitTime):
		fmt.Println("студент", workerNum, "придумал")
		out <- workerNum
	}
}

func main() {
	ctx, finish := context.WithCancel(context.Background()) // context with cancelling channel
	result := make(chan int, 1)

	for i := 0; i <= 10; i++ { // start 10 workers
		go студент(ctx, i, result)
	}

	foundBy := <-result // get first message
	fmt.Println("вопрос был задан студентом", foundBy)
	finish() // and send cancel message via context to all workers

	time.Sleep(time.Second)
}
