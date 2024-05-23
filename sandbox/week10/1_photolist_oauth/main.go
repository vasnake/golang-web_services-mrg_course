package v_oauth

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"time"

	// "math/rand"
	"net/http"
	// "time"

	_ "github.com/go-sql-driver/mysql"
)

func MainDemo() {
	// rand.Seed(time.Now().UnixNano())
	flag.StringVar(&APP_ID, "appid", "foo?", "app id (client id) from github registered app")
	flag.StringVar(&APP_SECRET, "appsecret", "bar?", "app secret (client key) from github registered app")
	flag.Parse()
	show("you mustn't but: appid, appsecret: ", APP_ID, APP_SECRET)

	// основные настройки к базе
	// dsn := "root:love@tcp(127.0.0.1:3306)/photolist?charset=utf8&interpolateParams=true"
	dsn := "root@tcp(127.0.0.1:3306)/photolist?charset=utf8&interpolateParams=true"
	db, err := sql.Open("mysql", dsn)
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalf("cant connect to db, err: %v\n", err)
	}
	defer db.Close()

	tmpls := NewTemplates()
	if err != nil {
		log.Fatalf("cant init tokens: %v\n", err)
	}

	// csrfTokens, err := NewHMACHashToken("golangcourse")
	// csrfTokens, err := NewAesCryptHashToken("qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB")
	csrfTokens, err := NewJwtToken("qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB")

	h := &PhotolistHandler{
		St:     NewDbStorage(db),
		Tmpl:   tmpls,
		Tokens: csrfTokens,
	}

	// sm := NewSessionsDB(db)
	// sm := NewSessionsJWT("golangcourseSessionSecret")
	sm := NewSessionsJWTVer("golangcourseSessionSecret", db)

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
	mux.HandleFunc("/user/login_oauth", u.LoginOauth)
	mux.HandleFunc("/user/logout", u.Logout)
	mux.HandleFunc("/user/reg", u.Reg)
	mux.HandleFunc("/user/change_pass", u.ChangePassword)
	mux.HandleFunc("/", Index)

	// http handlers with auth
	http.Handle("/", AuthMiddleware(sm, mux))

	// http handlers w/o auth
	staticHandler := http.StripPrefix(
		"/images/",
		http.FileServer(http.Dir("./images")),
	)
	http.Handle("/images/", staticHandler)

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./favicon.ico")
	})

	fmt.Println("starting server at :8080")
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
