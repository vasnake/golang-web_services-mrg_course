package auth

import (
	"database/sql"
	"fmt"
	"log"
	// "math/rand"
	// "time"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func MainAuth() {
	// rand.Seed(time.Now().UnixNano()) // deprecated

	// основные настройки к базе
	// dsn := "root:love@tcp(127.0.0.1:3306)/photolist?charset=utf8&interpolateParams=true"
	dsn := "root@tcp(127.0.0.1:3306)/photolist?charset=utf8&interpolateParams=true"
	db, err := sql.Open("mysql", dsn)
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalf("can't connect to db, err: %v\n", err)
	}
	defer db.Close()

	srvH := &PhotolistHandler{
		St:   NewDbStorage(db),
		Tmpl: NewTemplates(),
	}

	userH := &UserHandler{
		DB:   db,
		Tmpl: NewUserTemplates(),
	}

	mux := http.NewServeMux() // dispatch route handlers

	mux.HandleFunc("/photos/", srvH.List)
	mux.HandleFunc("/photos/upload", srvH.Upload)
	mux.HandleFunc("/user/login", userH.Login)
	mux.HandleFunc("/user/logout", userH.Logout)
	mux.HandleFunc("/user/reg", userH.Reg)
	mux.HandleFunc("/", Index)

	// add http hooks
	http.Handle("/", AuthMiddleware(db, mux))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))

	// run
	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

// Index: redirect to login or pics feed
func Index(w http.ResponseWriter, r *http.Request) {
	_, err := SessionFromContext(r.Context())

	// no session, go to login page
	if err != nil {
		http.Redirect(w, r, "/user/login", http.StatusFound)
		return
	}

	// have session, rock'n'roll
	http.Redirect(w, r, "/photos/", http.StatusFound)
}
