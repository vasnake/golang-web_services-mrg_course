package ajax3

import (
	"database/sql"
	"flag"
	"fmt"
	ioutil "io" // "io/ioutil"
	"log"
	"time"

	// "math/rand"
	"net/http"
	// "time"

	_ "github.com/go-sql-driver/mysql"
)

func MainDemo() {
	flag.StringVar(&APP_ID, "appid", "foo?", "oauth app id (client id) from github registered app")
	flag.StringVar(&APP_SECRET, "appsecret", "bar?", "oauth app secret (client key) from github registered app")
	flag.Parse()
	show("you must not show this! appid, appsecret: ", APP_ID, APP_SECRET)

	// rand.Seed(time.Now().UnixNano())

	// основные настройки к базе
	// dsn := "root:love@tcp(127.0.0.1:3306)/photolist?charset=utf8&interpolateParams=true"
	dsn := "root@tcp(127.0.0.1:3306)/photolist?charset=utf8&interpolateParams=true"
	db, err := sql.Open("mysql", dsn)
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalf("cant connect to db, err: %v\n", err)
	}
	defer db.Close()

	// csrfTokens, err := NewHMACHashToken("golangcourse")
	// csrfTokens, err := NewAesCryptHashToken("qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB")
	csrfTokens, err := NewJwtToken("qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB")

	templates := NewTemplates(Assets, csrfTokens)
	if err != nil {
		log.Fatalf("cant init templates: %v\n", err)
	}

	photosHandlers := &PhotolistHandler{
		St:     NewDbStorage(db),
		Tmpl:   templates,
		Tokens: csrfTokens,
		UserDB: db,
	}

	// sessions := NewSessionsDB(db)
	// sessions := NewSessionsJWT("golangcourseSessionSecret")
	sessions := NewSessionsJWTVer("golangcourseSessionSecret", db)

	usersHandlers := &UserHandler{
		DB:       db,
		Tmpl:     templates,
		Sessions: sessions,
	}

	// routes
	mux := http.NewServeMux()
	mux.HandleFunc("/photos/", photosHandlers.List) // photos/upload lands here

	mux.HandleFunc("/api/v1/photos/list", photosHandlers.ListAPI)
	mux.HandleFunc("/api/v1/photos/upload", photosHandlers.UploadAPI)
	mux.HandleFunc("/api/v1/photos/rate", photosHandlers.RateAPI)

	mux.HandleFunc("/user/login", usersHandlers.Login)
	mux.HandleFunc("/user/login_oauth", usersHandlers.LoginOauth)
	mux.HandleFunc("/user/logout", usersHandlers.Logout)
	mux.HandleFunc("/user/reg", usersHandlers.Reg)
	mux.HandleFunc("/user/change_pass", usersHandlers.ChangePassword)

	mux.HandleFunc("/api/v1/user/follow", usersHandlers.FollowAPI)
	mux.HandleFunc("/api/v1/user/following", usersHandlers.FollowingAPI)
	mux.HandleFunc("/api/v1/user/recomends", usersHandlers.RecomendsAPI)

	mux.HandleFunc("/", Index)

	// middlware stack (FIFO)
	handlers := AuthMiddleware(
		sessions,
		CsrfTokenMiddleware(csrfTokens, mux),
	)

	// http handlers

	http.Handle("/", handlers)

	imagesHandler := http.StripPrefix(
		"/images/",
		http.FileServer(http.Dir("./images")),
	)
	http.Handle("/images/", imagesHandler)

	http.Handle("/static/", http.FileServer(Assets))

	f, _ := Assets.Open("/static/favicon.ico")
	favicon, _ := ioutil.ReadAll(f)
	f.Close()
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Write(favicon)
	})

	// rock'n'roll

	log.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	_, err := SessionFromContext(r.Context())
	if err != nil {
		http.Redirect(w, r, "/user/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/photos/", http.StatusFound)
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
