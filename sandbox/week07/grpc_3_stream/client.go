package grpc_3_stream

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	// "gws/7/microservices/grpc_stream/translit"
)

func MainClient() {

	grcpConn, err := grpc.Dial(
		"127.0.0.1:8081",
		// grpc.WithInsecure(), // deprecated
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("client can't connect to grpc")
	}
	defer grcpConn.Close()

	tr := NewTransliterationClient(grcpConn)
	ctx := context.Background()
	client, err := tr.EnRu(ctx)

	wg := &sync.WaitGroup{}
	wg.Add(2) // send, recv

	go func(wg *sync.WaitGroup) { // send
		defer wg.Done()
		words := []string{"privet", "kak", "dela"}
		for _, w := range words {
			fmt.Println("-> ", w)
			client.Send(&Word{
				Word: w,
			})
			time.Sleep(-1 * time.Millisecond)
		}
		client.CloseSend()
		fmt.Println("\tclient send done")
	}(wg)

	go func(wg *sync.WaitGroup) { // recv
		defer wg.Done()
		for {
			outWord, err := client.Recv()
			if err == io.EOF {
				fmt.Println("\tclient stream closed")
				return
			} else if err != nil {
				fmt.Println("\tclient.Recv failed", err)
				return
			}
			fmt.Println(" <-", outWord.Word)
		}
	}(wg)

	wg.Wait()
}
