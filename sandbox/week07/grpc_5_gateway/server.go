package grpc_5_gateway

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	// "../session"
	// "gws/7/microservices/gateway/session"
)

/*
http client example requests

curl -X POST -k http://localhost:8080/v1/session/create -H "Content-Type: text/plain" -d '{"login":"rvasily", "useragent": "chrome"}'
curl http://localhost:8080/v1/session/check/XVlBzgbaiC
curl -X POST -k http://localhost:8080/v1/session/delete -H "Content-Type: text/plain" -d '{"ID":"XVlBzgbaiC"}'
*/

func MainServer() {
	proxyAddr := ":8080"
	serviceAddr := "127.0.0.1:8081"

	go gRPCService(serviceAddr)
	HTTPProxy(proxyAddr, serviceAddr)
}

func gRPCService(serviceAddr string) {
	listener, err := net.Listen("tcp", serviceAddr)
	if err != nil {
		log.Fatalln("failed to listen TCP port", err)
	}

	server := grpc.NewServer()

	RegisterAuthCheckerServer(server, NewSessionManager())

	fmt.Println("starting gRPC server at " + serviceAddr)
	server.Serve(listener)
}

func HTTPProxy(proxyAddr, serviceAddr string) {
	grcpConn, err := grpc.Dial(
		serviceAddr,
		// grpc.WithInsecure(), deprecated
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("failed to connect to grpc", err)
	}
	defer grcpConn.Close()

	grpcGWMux := runtime.NewServeMux()

	err = RegisterAuthCheckerHandler(
		context.Background(),
		grpcGWMux,
		grcpConn,
	)
	if err != nil {
		log.Fatalln("failed to start HTTP server", err)
	}

	mux := http.NewServeMux()
	// отправляем в прокси только то что нужно
	mux.Handle("/v1/session/", grpcGWMux)
	mux.HandleFunc("/", helloWorld)

	fmt.Println("starting HTTP server at " + proxyAddr)
	log.Fatal(http.ListenAndServe(proxyAddr, mux))
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "URL:", r.URL.String())
}
