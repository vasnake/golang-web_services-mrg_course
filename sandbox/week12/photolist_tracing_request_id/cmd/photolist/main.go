package main_program

import (
	"database/sql"
	"flag"
	"fmt"
	ioutil "io" // io/ioutil"
	"log"
	"net/http"
	"time"

	"week12/photolist_tracing_request_id/pkg/assets"
	"week12/photolist_tracing_request_id/pkg/blobstorage"
	"week12/photolist_tracing_request_id/pkg/config"
	"week12/photolist_tracing_request_id/pkg/graphql"
	"week12/photolist_tracing_request_id/pkg/middleware"
	"week12/photolist_tracing_request_id/pkg/photos"
	"week12/photolist_tracing_request_id/pkg/session"
	"week12/photolist_tracing_request_id/pkg/templates"
	"week12/photolist_tracing_request_id/pkg/token"
	"week12/photolist_tracing_request_id/pkg/user"

	gqlgen_handler "github.com/99designs/gqlgen/handler"
	_ "github.com/go-sql-driver/mysql"
)

var (
	appName   string = "photolist"
	buildHash string = "unknown"
	buildTime string = "unknown"
)

func MainDemo() { main() }

func main() {
	log.Printf("[startup] %s, commit %s, build %s", appName, buildHash, buildTime)

	flag.StringVar(&user.APP_ID, "appid", "foo?", "oauth app id (client id) from github registered app")
	flag.StringVar(&user.APP_SECRET, "appsecret", "bar?", "oauth app secret (client key) from github registered app")
	flag.Parse()
	show("you must not show this! appid, appsecret: ", user.APP_ID, user.APP_SECRET)

	cfg := &config.Config{}
	viperSvc, err := config.Read(appName, config.Defaults, cfg)
	if err != nil {
		log.Fatalf("can't read config, err: %v\n", err)
	}
	listenAddr := viperSvc.GetString("http.port")
	show("listen http.port from config: ", listenAddr)

	// основные настройки к базе
	// dsn := "root:love@tcp(127.0.0.1:3306)/photolist?charset=utf8&interpolateParams=true"
	dsn := "%s:%s@tcp(%s)/%s?charset=utf8&interpolateParams=true"
	dsn = fmt.Sprintf(dsn, cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Database)
	show("sql.Open mysql DSN: ", dsn)
	db, err := sql.Open("mysql", dsn)
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalf("can't connect to db, err: %v\n", err)
	}
	defer db.Close()

	// csrfTokens, err := token.NewHMACHashToken("golangcourse")
	// csrfTokens, err := token.NewAesCryptHashToken("qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB")
	csrfTokens, err := token.NewJwtToken(cfg.Token.Secret)
	if err != nil {
		log.Fatalf("can't init tokens: %v\n", err)
	}

	tmpls := templates.NewTemplates(assets.Assets, csrfTokens)
	if err != nil {
		log.Fatalf("can't init templates: %v\n", err)
	}

	// storage := blobstorage.NewFSStorage("./images/")
	storage, err := blobstorage.NewS3Storage(cfg.S3.Host, cfg.S3.Access, cfg.S3.Secret, cfg.S3.Bucket)
	if err != nil {
		log.Fatalln("can't create s3 blobstorage", err)
	}

	photosRepo := photos.NewPhotosRepository(db)
	usersRepo := user.NewUsersRepository(db)

	appSvc := &photos.PhotolistHandler{
		UsersRepo:   usersRepo,
		PhotosRepo:  photosRepo,
		Tmpl:        tmpls,
		BlobStorage: storage,
	}

	// sessMan := session.NewSessionsJWT("golangcourseSessionSecret")
	// sessMan := session.NewSessionsJWTVer(cfg.Session.Secret, db)
	sessMan := session.NewSessionsDB(db)

	userSvc := &user.UserHandler{
		Tmpl:      tmpls,
		Sessions:  sessMan,
		UsersRepo: usersRepo,
	}

	mux := http.NewServeMux()

	// mux.HandleFunc("/photos/", h.List)
	mux.HandleFunc("/photos/", appSvc.ListGQL)

	mux.HandleFunc("/api/v1/photos/list", appSvc.ListAPI)
	mux.HandleFunc("/api/v1/photos/upload", appSvc.UploadAPI)
	mux.HandleFunc("/api/v1/photos/rate", appSvc.RateAPI)

	mux.HandleFunc("/user/login", userSvc.Login)
	mux.HandleFunc("/user/login_oauth", userSvc.LoginOauth)
	mux.HandleFunc("/user/logout", userSvc.Logout)
	mux.HandleFunc("/user/reg", userSvc.Reg)
	mux.HandleFunc("/user/change_pass", userSvc.ChangePassword)

	mux.HandleFunc("/api/v1/user/follow", userSvc.FollowAPI)
	mux.HandleFunc("/api/v1/user/following", userSvc.FollowingAPI)
	mux.HandleFunc("/api/v1/user/recomends", userSvc.RecomendsAPI)

	mux.HandleFunc("/", Index)

	{ // START gqlgen part
		gqlResolver := &graphql.Resolver{
			PhotosRepo:  photosRepo,
			UsersRepo:   usersRepo,
			BlobStorage: storage,
		}
		gqlCfg := graphql.Config{
			Resolvers: gqlResolver,
		}
		gqlHandler := gqlgen_handler.GraphQL(
			graphql.NewExecutableSchema(gqlCfg),
			gqlgen_handler.ComplexityLimit(500),
			gqlgen_handler.RequestMiddleware(graphql.RequestMiddleware),   // каждый запрос после парсинга
			gqlgen_handler.ResolverMiddleware(graphql.ResolverMiddleware), // каждый вызов ресолвера
		)
		gqlHttpHandler := graphql.UserLoaderMiddleware(gqlResolver, gqlHandler)

		mux.Handle("/graphql", gqlHttpHandler)
		mux.Handle("/graphql/", gqlHttpHandler)
		// TODO enable csrf for graphql after playground done
		mux.HandleFunc("/playground", gqlgen_handler.Playground("GraphQL playground", "/graphql"))
	} // END gqlgen part

	// отрабатывают в обратном добавлению порядке, те AuthMiddleware будет 1-м
	photolistHttpHandlers := token.CsrfTokenMiddleware(csrfTokens, mux)
	photolistHttpHandlers = session.AuthMiddleware(sessMan, photolistHttpHandlers)
	photolistHttpHandlers = middleware.AccessLog(photolistHttpHandlers)
	photolistHttpHandlers = middleware.RequestIDMiddleware(photolistHttpHandlers)

	http.Handle("/", photolistHttpHandlers)

	http.HandleFunc("/api/v1/internal/images/auth", userSvc.InternalImagesAuth)

	imagesHandler := http.StripPrefix(
		"/images/",
		http.FileServer(http.Dir("./images")),
	)
	http.Handle("/images/", imagesHandler)

	http.Handle("/static/", http.FileServer(assets.Assets))

	f, err := assets.Assets.Open("/static/favicon.ico")
	if err == nil {
		favicon, err := ioutil.ReadAll(f)
		f.Close()
		if err == nil {
			http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
				w.Write(favicon)
			})
		} else {
			show("ReadAll favicon failed: ", err)
		}
	} else {
		show("assets.open favicon failed: ", err)
	}

	log.Printf("[startup] listening server at %s", listenAddr)
	http.ListenAndServe(listenAddr, nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	_, err := session.SessionFromContext(r.Context())
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
