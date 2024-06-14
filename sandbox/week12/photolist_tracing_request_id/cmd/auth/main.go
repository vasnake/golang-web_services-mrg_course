package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"photolist/pkg/config"
	"photolist/pkg/session"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	appName   string = "auth"
	buildHash string = "_dev"
	buildTime string = "_dev"
)

func AccessLogInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	reply, err := handler(ctx, req)
	md, mdExists := metadata.FromIncomingContext(ctx)
	var requestID string
	if mdExists {
		requestID = md["x-request-id"][0]
	} else {
		requestID = "-"
	}
	log.Printf("[access] %s %s %s '%v'", requestID, time.Since(start), info.FullMethod, err)
	return reply, err
}

func main() {
	log.Printf("[startup] %s, commit %s, build %s", appName, buildHash, buildTime)
	rand.Seed(time.Now().UnixNano())

	cfg := &config.Config{}
	v1, err := config.Read(appName, config.Defaults, cfg)
	if err != nil {
		log.Fatalf("[startup] cant read config, err: %v\n", err)
	}

	// основные настройки к базе
	dsn := "%s:%s@tcp(%s)/%s?charset=utf8&interpolateParams=true"
	dsn = fmt.Sprintf(dsn, cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Database)
	db, err := sql.Open("mysql", dsn)
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalf("[startup] cant connect to db, err: %v\n", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(AccessLogInterceptor),
	)
	svc := &session.AuthService{
		DB: db,
	}
	session.RegisterAuthServer(server, svc)

	listenAddr := v1.GetString("service.port")
	lis, err := net.Listen("tcp", ":"+listenAddr)
	if err != nil {
		log.Fatalln("[startup] cant listen port", err)
	}
	log.Printf("[startup] listening server at %s", listenAddr)
	server.Serve(lis)
}
