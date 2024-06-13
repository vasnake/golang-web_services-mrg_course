package main

import (
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
)

var (
	appName   string = "auth"
	buildHash string = "_dev"
	buildTime string = "_dev"
)

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

	server := grpc.NewServer()
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
