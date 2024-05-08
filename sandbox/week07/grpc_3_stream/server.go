package grpc_3_stream

import (
	"fmt"
	// "gws/7/microservices/grpc_stream/translit"
	"log"
	"net"

	"google.golang.org/grpc"
)

func MainServer() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("server can't listen port", err)
	}
	server := grpc.NewServer()

	RegisterTransliterationServer(server, NewTr())

	fmt.Println("starting server at :8081")
	server.Serve(listener)
}
