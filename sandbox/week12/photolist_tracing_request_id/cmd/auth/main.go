package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"week12/photolist_tracing_request_id/pkg/config"
	"week12/photolist_tracing_request_id/pkg/session"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	appName   string = "auth"
	buildHash string = "unknown"
	buildTime string = "unknown"
)

func rpcInterceptorAccessLog(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	rpcHandler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	reply, err := rpcHandler(ctx, req)

	var requestID string
	md, mdExists := metadata.FromIncomingContext(ctx)
	if mdExists {
		requestID = md["x-request-id"][0]
	} else {
		requestID = "unknown"
	}

	log.Printf("[access] %s %s %s '%v'", requestID, time.Since(start), info.FullMethod, err)
	return reply, err
}

func main() {
	log.Printf("[startup] %s, commit %s, build %s", appName, buildHash, buildTime)

	cfg := &config.Config{}
	viperSvc, err := config.Read(appName, config.Defaults, cfg)
	if err != nil {
		log.Fatalf("[startup] can't read config, err: %v\n", err)
	}

	listenAddr := viperSvc.GetString("service.port")
	show("service.port from config: ", listenAddr)

	// основные настройки к базе
	dsn := "%s:%s@tcp(%s)/%s?charset=utf8&interpolateParams=true"
	dsn = fmt.Sprintf(dsn, cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Database)
	show("sql.Open mysql DSN: ", dsn)
	db, err := sql.Open("mysql", dsn)
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalf("[startup] can't connect to db, err: %v\n", err)
	}
	defer db.Close()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(rpcInterceptorAccessLog),
	)

	authSvc := &session.AuthService{
		DB: db,
	}

	session.RegisterAuthServer(grpcServer, authSvc)

	lsnr, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalln("[startup] can't listen port", err)
	}

	log.Printf("[startup] listening server at %s", listenAddr)
	grpcServer.Serve(lsnr)
}

// ts returns current timestamp in RFC3339 with milliseconds
func ts() string {
	/*
		https://pkg.go.dev/time#pkg-constants
		https://stackoverflow.com/questions/35479041/how-to-convert-iso-8601-time-in-golang
	*/
	const (
		RFC3339      = "2006-01-02T15:04:05Z07:00"
		RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
	)
	return time.Now().UTC().Format(RFC3339Milli)
}

// show writes message to standard output. Message combined from prefix msg and slice of arbitrary arguments
func show(msg string, xs ...any) {
	var line = ts() + ": " + msg

	for _, x := range xs {
		// https://pkg.go.dev/fmt
		// line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
}
