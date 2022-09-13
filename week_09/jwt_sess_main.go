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

	tmpls := NewTemplates()

	// tokens, err := NewHMACHashToken("golangcourseCsrfSecret")
	// tokens, err := NewAesCryptHashToken("qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB")
	tokens, err := NewJwtToken("qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB")

	if err != nil {
		log.Fatalf("cant init tokens: %v\n", err)
	}

	h := &PhotolistHandler{
		St:     NewDbStorage(db),
		Tmpl:   tmpls,
		Tokens: tokens,
	}

	// added session manager, interface for different implementations support
	// sm := NewSessionsDB(db)
	// sm := NewSessionsJWT("golangcourseSessionSecret") // secret for JWT sign/verify procedures
	sm := NewSessionsJWTVer("golangcourseSessionSecret", db) // db with user profile, versioned

	// sessions in user handlers
	u := &UserHandler{
		DB:       db,
		Tmpl:     tmpls,
		Sessions: sm,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/photos/", h.List)
	mux.HandleFunc("/photos/upload", h.Upload)
	mux.HandleFunc("/photos/rate", h.Rate)

	mux.HandleFunc("/user/login", u.Login)
	mux.HandleFunc("/user/logout", u.Logout)
	mux.HandleFunc("/user/reg", u.Reg)
	mux.HandleFunc("/user/change_pass", u.ChangePassword)

	mux.HandleFunc("/", Index)

	// sessions in auth
	http.Handle("/", AuthMiddleware(sm, mux))

	staticHandler := http.StripPrefix(
		"/images/",
		http.FileServer(http.Dir("./images")),
	)
	http.Handle("/images/", staticHandler)

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./favicon.ico")
	})

	// http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./favicon.ico")
	// })

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
