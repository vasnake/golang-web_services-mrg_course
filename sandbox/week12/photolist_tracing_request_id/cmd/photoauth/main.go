package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"photolist/pkg/config"
	"photolist/pkg/middleware"
	"photolist/pkg/session"
	"photolist/pkg/user"

	_ "github.com/go-sql-driver/mysql"
)

var (
	appName   string = "photoauth"
	buildHash string = "_dev"
	buildTime string = "_dev"
)

func main() {
	log.Printf("[startup] %s, commit %s, build %s", appName, buildHash, buildTime)
	rand.Seed(time.Now().UnixNano())

	cfg := &config.Config{}
	v1, err := config.Read(appName, config.Defaults, cfg)
	if err != nil {
		log.Fatalf("cant read config, err: %v\n", err)
	}

	env1 := v1.GetString("example.env1")
	env2 := v1.GetString("example.env2")
	log.Printf("[startup] cfg: %#v, env1 %#v, env2 %#v", cfg.HTTP.Port, env1, env2)

	// основные настройки к базе
	// dsn := "root:love@tcp(host.docker.internal:3306)/photolist?charset=utf8&interpolateParams=true"
	dsn := "%s:%s@tcp(%s)/%s?charset=utf8&interpolateParams=true"
	dsn = fmt.Sprintf(dsn, cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Database)
	db, err := sql.Open("mysql", dsn)
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalf("[startup] cant connect to db, err: %v\n", err)
	}

	usersRepo := user.NewUsersRepository(db)
	log.Println("sess grpc addr:", v1.GetString("session.grpc_addr"))
	sm, err := session.NewSessionsGRPC(v1.GetString("session.grpc_addr"))
	if err != nil {
		log.Fatalf("[startup] cant connect to session grpc, err: %v\n", err)
	}

	u := &user.UserHandler{
		Tmpl:      nil,
		Sessions:  sm,
		UsersRepo: usersRepo,
	}

	handlers := middleware.AccessLog(http.HandlerFunc(u.InternalImagesAuth))
	handlers = middleware.RequestIDMiddleware(handlers)

	http.Handle("/api/v1/internal/images/auth", handlers)

	listenAddr := ":" + v1.GetString("http.port")
	log.Printf("[startup] listening server at %s", listenAddr)
	http.ListenAndServe(listenAddr, nil)
}
