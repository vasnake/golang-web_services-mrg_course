package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"

	"coursera/microservices/grpc_stream/translit"
)

func main() {

	// open connection
	grcpConn, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpConn.Close()

	// get grpc endpoint on that connection
	tr := translit.NewTransliterationClient(grcpConn)

	ctx := context.Background()
	// get service
	client, err := tr.EnRu(ctx)

	// wait, two goroutine
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// first, sending words to stream
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		words := []string{"privet", "kak", "dela"}

		for _, w := range words {
			fmt.Println("-> ", w)
			client.Send(&translit.Word{
				Word: w,
			})
			time.Sleep(time.Millisecond)
		}

		client.CloseSend() // stream closed
		fmt.Println("\tsend done")
	}(wg)

	// second, recieving words from stream
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		for {
			outWord, err := client.Recv()
			if err == io.EOF {
				fmt.Println("\tstream closed")
				return
			} else if err != nil {
				fmt.Println("\terror happed", err)
				return
			}

			fmt.Println(" <-", outWord.Word)
		}
	}(wg)

	wg.Wait()

}
