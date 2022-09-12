package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// основные настройки к базе
	dsn := "root:love@tcp(127.0.0.1:3306)/photolist?charset=utf8&interpolateParams=true"
	db, err := sql.Open("mysql", dsn)

	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalf("cant connect to db, err: %v\n", err)
	}

	h := &PhotolistHandler{
		St:   NewDbStorage(db),
		Tmpl: NewTemplates(),
	}

	// aded user session support
	u := &UserHandler{
		DB:   db,
		Tmpl: NewUserTemplates(),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/photos/", h.List)
	mux.HandleFunc("/photos/upload", h.Upload)

	mux.HandleFunc("/user/login", u.Login)
	mux.HandleFunc("/user/logout", u.Logout)
	mux.HandleFunc("/user/reg", u.Reg)

	mux.HandleFunc("/", Index)

	// wrap all to auth
	http.Handle("/", AuthMiddleware(db, mux))

	staticHandler := http.StripPrefix(
		"/images/",
		http.FileServer(http.Dir("./images")),
	)
	// files w/o auth
	http.Handle("/images/", staticHandler)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
