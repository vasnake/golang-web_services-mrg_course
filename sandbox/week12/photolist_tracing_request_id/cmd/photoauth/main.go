package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"week12/photolist_tracing_request_id/pkg/config"
	"week12/photolist_tracing_request_id/pkg/middleware"
	"week12/photolist_tracing_request_id/pkg/session"
	"week12/photolist_tracing_request_id/pkg/user"

	_ "github.com/go-sql-driver/mysql"
)

var (
	appName   string = "photoauth"
	buildHash string = "unknown"
	buildTime string = "unknown"
)

func main() {
	log.Printf("[startup] %s, commit %s, build %s", appName, buildHash, buildTime)

	cfg := &config.Config{}
	viperSvc, err := config.Read(appName, config.Defaults, cfg)
	if err != nil {
		log.Fatalf("can't read config, err: %v\n", err)
	}

	env1 := viperSvc.GetString("example.env1")
	env2 := viperSvc.GetString("example.env2")
	log.Printf("[startup] cfg.HTTP.Port %#v, example.env1 %#v, example.env2 %#v", cfg.HTTP.Port, env1, env2)

	grpcAddr := viperSvc.GetString("session.grpc_addr")
	show("downstream svc session.grpc_addr from config: ", grpcAddr)

	listenAddr := viperSvc.GetString("http.port")
	show("listen http.port from config: ", listenAddr)

	// основные настройки к базе
	// dsn := "root:love@tcp(host.docker.internal:3306)/photolist?charset=utf8&interpolateParams=true"
	dsn := "%s:%s@tcp(%s)/%s?charset=utf8&interpolateParams=true"
	dsn = fmt.Sprintf(dsn, cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Database)
	show("sql.Open mysql DSN: ", dsn)
	db, err := sql.Open("mysql", dsn)
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalf("[startup] can't connect to db, err: %v\n", err)
	}
	defer db.Close()

	sessMan, err := session.NewSessionsGRPC(grpcAddr)
	if err != nil {
		log.Fatalf("[startup] can't connect to session grpc, err: %v\n", err)
	}

	usersRepo := user.NewUsersRepository(db)

	userSvc := &user.UserHandler{
		Tmpl:      nil,
		Sessions:  sessMan,
		UsersRepo: usersRepo,
	}

	httpHandler := middleware.AccessLog(http.HandlerFunc(userSvc.InternalImagesAuth))
	httpHandler = middleware.RequestIDMiddleware(httpHandler)
	http.Handle("/api/v1/internal/images/auth", httpHandler)

	log.Printf("[startup] listening server at %s", listenAddr)
	http.ListenAndServe(listenAddr, nil)
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
