package csrf_token

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	// "time"
	// "math/rand"

	_ "github.com/go-sql-driver/mysql"
)

func MainCsrf() {
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

	pageTemplates := NewTemplates()

	// tokenHandlers, err := NewHMACHashToken("golangcourse")
	// tokenHandlers, err := NewAesCryptHashToken("qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB")
	tokenHandlers, err := NewJwtToken("qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB") // toy example, don't do that
	if err != nil {
		log.Fatalf("cant init tokens: %v\n", err)
	}

	h := &PhotolistHandler{
		St:     NewDbStorage(db),
		Tmpl:   pageTemplates,
		Tokens: tokenHandlers,
	}

	u := &UserHandler{
		DB:   db,
		Tmpl: pageTemplates,
	}

	// route handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/photos/", h.List)
	mux.HandleFunc("/photos/upload", h.Upload)
	mux.HandleFunc("/photos/rate", h.Rate)
	mux.HandleFunc("/user/login", u.Login)
	mux.HandleFunc("/user/logout", u.Logout)
	mux.HandleFunc("/user/reg", u.Reg)
	mux.HandleFunc("/", Index)

	// http handlers
	http.Handle("/", AuthMiddleware(db, mux))
	// files, w/o auth
	http.Handle("/images/", http.StripPrefix(
		"/images/",
		http.FileServer(http.Dir("./images")),
	))
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
