package grpc_1

import (
	"fmt"
	// "gws/7/microservices/grpc/session"
	"log"
	"net"

	"google.golang.org/grpc"
)

func MainServer() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("can't listen port 8081", err)
	}

	server := grpc.NewServer()

	RegisterAuthCheckerServer(server, NewSessionManager())

	fmt.Println("starting server at :8081")
	server.Serve(listener)
}
