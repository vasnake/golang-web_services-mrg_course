package grpc_2

import (
	"fmt"
	// "gws/7/microservices/grpc/session"
	"log"
	"net"
	"time"
	"week07/grpc_1"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/tap"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

// grpc.NewServer(grpc.UnaryInterceptor(authInterceptor), ...)
func authInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)

	// do work
	reply, err := handler(ctx, req)

	// process context
	fmt.Printf(`--
	after incoming call=%v
	req=%#v
	reply=%#v
	time=%v
	md=%v
	err=%v
`, info.FullMethod, req, reply, time.Since(start), md, err)
	return reply, err
}

// -----

// grpc.NewServer(..., grpc.InTapHandle(rateLimiter), ...)
func rateLimiter(ctx context.Context, info *tap.Info) (context.Context, error) {
	// before heavy processing
	fmt.Printf("--\ncheck ratelim for %s\n", info.FullMethodName)
	return ctx, nil
}

// -----

func MainServer() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("can't listen port 8081", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor),
		grpc.InTapHandle(rateLimiter),
	)

	grpc_1.RegisterAuthCheckerServer(server, grpc_1.NewSessionManager())

	fmt.Println("starting server at :8081")
	server.Serve(listener)
}
