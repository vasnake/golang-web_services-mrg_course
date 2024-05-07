package grpc_2

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	// "gws/7/microservices/grpc/session"
	"week07/grpc_1"
)

// grpc.Dial(..., grpc.WithUnaryInterceptor(timingInterceptor), ...)
func timingInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()

	// do work
	err := invoker(ctx, method, req, reply, cc, opts...)

	// do context stuff
	fmt.Printf(`--
	call=%v
	req=%#v
	reply=%#v
	time=%v
	err=%v
`, method, req, reply, time.Since(start), err)

	return err
}

// -----

// grpc.Dial(..., grpc.WithPerRPCCredentials(&tokenAuth{"100500"}), ...)
type tokenAuth struct {
	Token string
}

func (t *tokenAuth) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"access-token": t.Token,
	}, nil
}

func (c *tokenAuth) RequireTransportSecurity() bool {
	return false
}

// -----

func MainClient() {

	grcpConn, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithUnaryInterceptor(timingInterceptor),
		grpc.WithPerRPCCredentials(&tokenAuth{"100500"}),
		// grpc.WithInsecure(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("can't connect to grpc")
	}
	defer grcpConn.Close()

	sessManager := grpc_1.NewAuthCheckerClient(grcpConn)

	// context metadata demo
	md := metadata.Pairs(
		"api-req-id", "123",
		"subsystem", "cli",
	)
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)

	// ----------------------------------------------------

	// prefix, suffix demo
	var header, trailer metadata.MD

	// создаем сессию
	sessId, err := sessManager.Create(ctx,
		&grpc_1.Session{
			Login:     "rvasily",
			Useragent: "chrome",
		},
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	fmt.Println("sessId", sessId, err)
	fmt.Println("header", header)
	fmt.Println("trailer", trailer)

	// проеряем сессию
	sess, err := sessManager.Check(ctx,
		&grpc_1.SessionID{
			ID: sessId.ID,
		})
	fmt.Println("sess", sess, err)

	// удаляем сессию
	_, err = sessManager.Delete(ctx,
		&grpc_1.SessionID{
			ID: sessId.ID,
		})

	// проверяем еще раз
	sess, err = sessManager.Check(ctx,
		&grpc_1.SessionID{
			ID: sessId.ID,
		})
	fmt.Println("sess", sess, err)
}
