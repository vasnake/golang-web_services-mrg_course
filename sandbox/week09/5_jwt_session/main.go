package jwt_session

import (
	"database/sql"
	"fmt"
	"log"
	// "math/rand"
	"net/http"
	// "time"

	_ "github.com/go-sql-driver/mysql"
)

// photolist demo: jwt csrf token, jwt session cookie
func MainJwtSession() {
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

	tmpls := NewTemplates()

	// stateless (almost) session
	// works fast (no db) but w/o 'destroy session(s)' functionality
	sessionManagerH := NewSessionsJWTVer("golangcourseSessionSecret", db) // session service, SessionManager using jwt cookie
	// completely stateless session
	// sessionManagerH := NewSessionsJWT("golangcourseSessionSecret") // w/o user version, even faster

	// stateful session // good solution, but slow
	// sessionManagerH := NewSessionsDB(db)

	userH := &UserHandler{
		DB:       db,
		Tmpl:     tmpls,
		Sessions: sessionManagerH, // session service
	}

	// CSRF tokens using JWT
	// tokens, err := NewHMACHashToken("golangcourseCsrfSecret")
	// tokens, err := NewAesCryptHashToken("qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB")
	tokens, err := NewJwtToken("qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB")
	if err != nil {
		log.Fatalf("cant init tokens: %v\n", err)
	}

	svcH := &PhotolistHandler{
		St:     NewDbStorage(db), // storage service
		Tmpl:   tmpls,
		Tokens: tokens, // token service, csrf
	}

	// route handlers

	mux := http.NewServeMux()
	mux.HandleFunc("/photos/", svcH.List)
	mux.HandleFunc("/photos/upload", svcH.Upload)
	mux.HandleFunc("/photos/rate", svcH.Rate)

	mux.HandleFunc("/user/login", userH.Login)
	mux.HandleFunc("/user/logout", userH.Logout)
	mux.HandleFunc("/user/reg", userH.Reg)
	mux.HandleFunc("/user/change_pass", userH.ChangePassword)

	mux.HandleFunc("/", Index)

	// http handlers

	http.Handle("/", AuthMiddleware(sessionManagerH, mux)) // use jwt session cookie (check user version against db)
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
