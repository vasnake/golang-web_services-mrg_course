package main

import (
	"database/sql"
	"io/ioutil"
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

	// tokens, err := NewHMACHashToken("golangcourse")
	// tokens, err := NewAesCryptHashToken("qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB")
	tokens, err := NewJwtToken("qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB")

	tmpls := NewTemplates(Assets, tokens)
	if err != nil {
		log.Fatalf("cant init tokens: %v\n", err)
	}

	h := &PhotolistHandler{
		St:     NewDbStorage(db),
		Tmpl:   tmpls,
		Tokens: tokens,
		UserDB: db,
	}

	// sm := NewSessionsDB(db)
	// sm := NewSessionsJWT("golangcourseSessionSecret")
	sm := NewSessionsJWTVer("golangcourseSessionSecret", db) // persistent session

	u := &UserHandler{
		DB:       db,
		Tmpl:     tmpls,
		Sessions: sm,
	}

	mux := http.NewServeMux()

	// API
	mux.HandleFunc("/api/v1/photos/list", h.ListAPI)
	mux.HandleFunc("/api/v1/photos/upload", h.UploadAPI)
	mux.HandleFunc("/api/v1/photos/rate", h.RateAPI)

	mux.HandleFunc("/api/v1/user/follow", u.FollowAPI)
	mux.HandleFunc("/api/v1/user/following", u.FollowingAPI)
	mux.HandleFunc("/api/v1/user/recomends", u.RecomendsAPI)

	// pages

	mux.HandleFunc("/photos/", h.List)

	mux.HandleFunc("/user/login", u.Login)
	mux.HandleFunc("/user/login_oauth", u.LoginOauth)
	mux.HandleFunc("/user/logout", u.Logout)
	mux.HandleFunc("/user/reg", u.Reg)
	mux.HandleFunc("/user/change_pass", u.ChangePassword)

	mux.HandleFunc("/", Index)

	// middleware

	// отрабатывают в обратном добавлению порядке, те AuthMiddleware будет 1-м
	handlers := CsrfTokenMiddleware(tokens, mux)
	handlers = AuthMiddleware(sm, handlers)
	http.Handle("/", handlers)

	// static

	imagesHandler := http.StripPrefix(
		"/images/",
		http.FileServer(http.Dir("./images")),
	)
	http.Handle("/images/", imagesHandler)

	http.Handle("/static/", http.FileServer(Assets))

	f, _ := Assets.Open("/static/favicon.ico")
	defer f.Close()
	favicon, _ := ioutil.ReadAll(f)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Write(favicon)
	})

	// run

	log.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
