package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"photolist/pkg/assets"
	"photolist/pkg/graphql"
	"photolist/pkg/index"
	"photolist/pkg/photos"
	"photolist/pkg/session"
	"photolist/pkg/templates"
	"photolist/pkg/token"
	"photolist/pkg/user"

	gqlgenHandler "github.com/99designs/gqlgen/handler"
	_ "github.com/go-sql-driver/mysql"
)

var (
	appName   string = "photolist"
	buildHash string = "_dev"
	buildTime string = "_dev"
)

func main() {
	log.Printf("starting %s, commit %s, build %s", appName, buildHash, buildTime)

	rand.Seed(time.Now().UnixNano())

	// основные настройки к базе
	dsn := "root:love@tcp(127.0.0.1:3306)/photolist?charset=utf8&interpolateParams=true"
	db, err := sql.Open("mysql", dsn)

	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalf("cant connect to db, err: %v\n", err)
	}

	// tokens, err := token.NewHMACHashToken("golangcourse")
	// tokens, err := token.NewAesCryptHashToken("qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB")
	tokens, err := token.NewJwtToken("qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB")

	tmpls := templates.NewTemplates(assets.Assets, tokens)
	if err != nil {
		log.Fatalf("cant init tokens: %v\n", err)
	}

	photosRepo := photos.NewDbStorage(db)
	usersRepo := user.NewUsersRepository(db)

	h := &photos.PhotolistHandler{
		UsersRepo: usersRepo,
		St:        photosRepo,
		Tmpl:      tmpls,
	}

	// sm := session.NewSessionsDB(db)
	// sm := session.NewSessionsJWT("golangcourseSessionSecret")
	sm := session.NewSessionsJWTVer("golangcourseSessionSecret", db)

	u := &user.UserHandler{
		Tmpl:      tmpls,
		Sessions:  sm,
		UsersRepo: usersRepo,
	}

	mux := http.NewServeMux()

	// mux.HandleFunc("/photos/", h.List)
	mux.HandleFunc("/photos/", h.ListGQL)

	mux.HandleFunc("/api/v1/photos/list", h.ListAPI)
	mux.HandleFunc("/api/v1/photos/upload", h.UploadAPI)
	mux.HandleFunc("/api/v1/photos/rate", h.RateAPI)

	mux.HandleFunc("/user/login", u.Login)
	mux.HandleFunc("/user/login_oauth", u.LoginOauth)
	mux.HandleFunc("/user/logout", u.Logout)
	mux.HandleFunc("/user/reg", u.Reg)
	mux.HandleFunc("/user/change_pass", u.ChangePassword)

	mux.HandleFunc("/api/v1/user/follow", u.FollowAPI)
	mux.HandleFunc("/api/v1/user/following", u.FollowingAPI)
	mux.HandleFunc("/api/v1/user/recomends", u.RecomendsAPI)

	mux.HandleFunc("/", index.Index)

	// START gqlgen part
	resolver := &graphql.Resolver{
		PhotosRepo: photosRepo,
		UsersRepo:  usersRepo,
	}
	cfg := graphql.Config{
		Resolvers: resolver,
	}
	gqlHandler := gqlgenHandler.GraphQL(
		graphql.NewExecutableSchema(cfg),
		gqlgenHandler.ComplexityLimit(500),
	)
	myGqlHandler := graphql.UserLoaderMiddleware(resolver, gqlHandler)

	mux.Handle("/graphql", myGqlHandler)
	// TODO enable csrf for graphql after playground done
	mux.HandleFunc("/playground", gqlgenHandler.Playground("GraphQL playground", "/graphql"))
	// END gqlgen part

	// отрабатывают в обратном добавлению порядке, те AuthMiddleware будет 1-м
	handlers := token.CsrfTokenMiddleware(tokens, mux)
	handlers = session.AuthMiddleware(sm, handlers)

	http.Handle("/", handlers)

	imagesHandler := http.StripPrefix(
		"/images/",
		http.FileServer(http.Dir("./images")),
	)
	http.Handle("/images/", imagesHandler)

	http.Handle("/static/", http.FileServer(assets.Assets))

	f, _ := assets.Assets.Open("/static/favicon.ico")
	defer f.Close()
	favicon, _ := ioutil.ReadAll(f)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Write(favicon)
	})

	listenAddr := ":8080"
	log.Printf("starting listening server at %s", listenAddr)
	http.ListenAndServe(listenAddr, nil)
}
