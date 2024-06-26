package main

import (
	"coursera/microservices/grpc/session"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"

	consulapi "github.com/hashicorp/consul/api"
)

var (
	// command line parameters, consul address; our service address
	grpcPort   = flag.Int("grpc", 8081, "listen addr")
	consulAddr = flag.String("consul", "192.168.99.100:32769", "consul addr (8500 in original consul)")
)

/*
	go run *.go --grpc="8081" --consul="192.168.99.100:32769"
	go run *.go --grpc="8082" --consul="192.168.99.100:32769"
*/

func main() {
	flag.Parse() // commandline params
	port := strconv.Itoa(*grpcPort)
	serviceID := "SAPI_127.0.0.1:" + port // our service ID string

	// setup service

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalln("cant listen port", err)
	}

	server := grpc.NewServer()
	session.RegisterAuthCheckerServer(server, NewSessionManager(port))

	fmt.Println("starting server at " + port)
	go server.Serve(lis)

	// ask consul to register service

	config := consulapi.DefaultConfig()
	config.Address = *consulAddr
	consul, err := consulapi.NewClient(config)

	err = consul.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
		ID:      serviceID,
		Name:    "session-api",
		Port:    *grpcPort,
		Address: "127.0.0.1",
	})
	if err != nil {
		fmt.Println("cant add service to consul", err)
		return
	}
	fmt.Println("registered in consul", serviceID)

	// on exit: remove record from consul
	defer func() {
		err := consul.Agent().ServiceDeregister(serviceID)
		if err != nil {
			fmt.Println("cant add service to consul", err)
			return
		}
		fmt.Println("deregistered in consul", serviceID)
	}()

	fmt.Println("Press any key to exit")
	fmt.Scanln()
}
