package grpc_4_balanced

import (
	"flag"
	"fmt"
	"time"

	// "gws/7/microservices/grpc/session"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"

	consulapi "github.com/hashicorp/consul/api"
)

var (
	grpcPort   = flag.Int("grpc", 8081, "listen addr")
	consulAddr = flag.String("consul", "127.0.0.1:8500", "consul addr (8500 in original consul)")
)

/*
	go run *.go --grpc="8081" --consul="127.0.0.1:8500"
	go run *.go --grpc="8082" --consul="127.0.0.1:8500"
*/

func MainServer_1() {
	flag.Parse()
	port := strconv.Itoa(*grpcPort)
	port = "8081"
	MainServer(port)
}
func MainServer_2() {
	flag.Parse()
	port := strconv.Itoa(*grpcPort)
	port = "8082"
	MainServer(port)
}

func MainServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalln("server can't listen port", err)
	}

	server := grpc.NewServer()

	RegisterAuthCheckerServer(server, NewSessionManager(port)) // run session manager grpc server

	config := consulapi.DefaultConfig()
	config.Address = *consulAddr

	serviceID := "SAPI_127.0.0.1:" + port
	consul, err := consulapi.NewClient(config)

	var deReg = func(serviceID string) {
		err := consul.Agent().ServiceDeregister(serviceID)
		if err != nil {
			fmt.Println("server can't remove service from consul", err)
			return
		}
		fmt.Println("server deregistered in consul", serviceID)
	}
	// deReg(serviceID) // dirty hack

	err = consul.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
		ID:      serviceID,
		Name:    "session-api",
		Port:    strconvAtoi(port),
		Address: "127.0.0.1",
	})
	if err != nil {
		fmt.Println("server can't add service to consul", err)
		return
	}
	fmt.Println("server registered in consul", serviceID)

	// defer deReg(serviceID)
	time.AfterFunc(5*time.Second, func() { deReg(serviceID); server.Stop() }) // dirty hack

	fmt.Println("starting server at " + port)
	server.Serve(listener)

	// go server.Serve(listener)
	// fmt.Println("Press any key to exit")
	// fmt.Scanln()
}

func strconvAtoi(n string) int {
	i, err := strconv.Atoi(n)
	if err != nil {
		fmt.Println("strconvAtoi failed on", n)
	}
	return i
}
