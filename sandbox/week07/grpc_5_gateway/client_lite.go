package grpc_5_gateway

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	// "../session"
	// "gws/7/microservices/gateway/session"
)

// grpc client demo
func MainClient() {

	grcpConn, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("client can't connect to grpc")
	}
	defer grcpConn.Close()

	sessManager := NewAuthCheckerClient(grcpConn)
	ctx := context.Background()

	// создаем сессию
	sessId, err := sessManager.Create(ctx,
		&Session{
			Login:     "rvasily",
			Useragent: "chrome",
		})
	fmt.Println("sessId", sessId, err)

	// проеряем сессию
	sess, err := sessManager.Check(ctx,
		&SessionID{
			ID: sessId.ID,
		})
	fmt.Println("sess", sess, err)

	// удаляем сессию
	_, err = sessManager.Delete(ctx,
		&SessionID{
			ID: sessId.ID,
		})

	// проверяем еще раз
	sess, err = sessManager.Check(ctx,
		&SessionID{
			ID: sessId.ID,
		})
	fmt.Println("sess", sess, err)
}
