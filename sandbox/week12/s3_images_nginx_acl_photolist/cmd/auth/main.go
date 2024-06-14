package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"week12/s3_images_nginx_acl_photolist/pkg/config"
	"week12/s3_images_nginx_acl_photolist/pkg/session"

	"github.com/carlmjohnson/versioninfo"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

var (
	appName string = "auth"
	// COMMIT?=$(shell git rev-parse --short HEAD) BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S') go build -ldflags "-X main.buildHash=${COMMIT} -X main.buildTime=${BUILD_TIME}" ...
	// modern way: https://blog.carlana.net/post/2023/golang-git-hash-how-to/
	buildHash string = "unknown"
	buildTime string = "unknown"
)

func main() {
	log.Printf("[startup 1] service '%s', -ldflags info: buildHash '%s', buildTime '%s'", appName, buildHash, buildTime)

	log.Printf("[startup 2] service '%s', 'go build -buildvcs' info: Version '%s', Revision '%s', DirtyBuild '%s', LastCommit '%s', ShortInfo '%s'",
		appName,
		versioninfo.Version, versioninfo.Revision, versioninfo.DirtyBuild, versioninfo.LastCommit, versioninfo.Short())

	cfg := &config.Config{}
	viperSvc, err := config.Read(appName, config.Defaults, cfg)
	if err != nil {
		log.Fatalf("[startup 3] can't read config, err: %v\n", err)
	}

	servicePort := viperSvc.GetString("service.port")
	show("service.port from config: ", servicePort)

	dsn := "%s:%s@tcp(%s)/%s?charset=utf8&interpolateParams=true"
	dsn = fmt.Sprintf(dsn, cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Database)
	show("sql.Open mysql DSN: ", dsn)

	db, err := sql.Open("mysql", dsn)
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalf("[startup 4] can't connect to db, err: %v\n", err)
	}
	defer db.Close()

	grpcSrv := grpc.NewServer()
	userSessions := &session.AuthService{
		DB: db,
	}

	session.RegisterAuthServer(grpcSrv, userSessions)

	lis, err := net.Listen("tcp", servicePort)
	if err != nil {
		log.Fatalln("[startup 5] can't listen port", err)
	}

	log.Printf("[startup 6] grpc serve at tcp port '%s'", servicePort)
	grpcSrv.Serve(lis)
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
