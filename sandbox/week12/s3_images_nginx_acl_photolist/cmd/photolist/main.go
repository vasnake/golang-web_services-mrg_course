package main_program

import (
	"database/sql"
	"flag"
	"fmt"
	ioutil "io" // "io/ioutil"
	"log"
	"net/http"
	"time"

	"week12/s3_images_nginx_acl_photolist/pkg/assets"
	"week12/s3_images_nginx_acl_photolist/pkg/blobstorage"
	"week12/s3_images_nginx_acl_photolist/pkg/config"
	"week12/s3_images_nginx_acl_photolist/pkg/graphql"
	"week12/s3_images_nginx_acl_photolist/pkg/photos"
	"week12/s3_images_nginx_acl_photolist/pkg/session"
	"week12/s3_images_nginx_acl_photolist/pkg/templates"
	"week12/s3_images_nginx_acl_photolist/pkg/token"
	"week12/s3_images_nginx_acl_photolist/pkg/user"

	gqlgenHandler "github.com/99designs/gqlgen/handler"
	_ "github.com/go-sql-driver/mysql"
)

var (
	appName string = "photolist"
	// COMMIT?=$(shell git rev-parse --short HEAD) BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S') go build -ldflags "-X main.buildHash=${COMMIT} -X main.buildTime=${BUILD_TIME}" ...
	buildHash string = "unknown"
	buildTime string = "unknown"
)

func MainDemo() {
	main()
}

func main() {
	flag.StringVar(&user.APP_ID, "appid", "foo?", "oauth app id (client id) from github registered app")
	flag.StringVar(&user.APP_SECRET, "appsecret", "bar?", "oauth app secret (client key) from github registered app")
	flag.Parse()
	log.Printf("[startup] %s, commit %s, build %s", appName, buildHash, buildTime)
	show("you must not show this! appid, appsecret: ", user.APP_ID, user.APP_SECRET)

	cfg := &config.Config{}
	viperSvc, err := config.Read(appName, config.Defaults, cfg)
	// log.Println(viperSvc, cfg, err, viperSvc.GetString("example.env2"))
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

	// s3Storage := blobstorage.NewFSStorage("./images/")
	s3Storage, err := blobstorage.NewS3Storage(cfg.S3.Host, cfg.S3.Access, cfg.S3.Secret, cfg.S3.Bucket)
	if err != nil {
		log.Fatalln("can't create s3 blobstorage", err)
	}

	// csrfTokens, err := token.NewHMACHashToken("golangcourse")
	// csrfTokens, err := token.NewAesCryptHashToken("qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB")
	csrfTokens, err := token.NewJwtToken(cfg.Token.Secret)
	pageTemplates := templates.NewTemplates(assets.Assets, csrfTokens)
	if err != nil {
		log.Fatalf("can't init tokens: %v\n", err)
	}

	photosRepo := photos.NewPhotosRepository(db)
	usersRepo := user.NewUsersRepository(db)

	appHandlers := &photos.PhotolistHandler{
		UsersRepo:   usersRepo,
		PhotosRepo:  photosRepo,
		Tmpl:        pageTemplates,
		BlobStorage: s3Storage,
	}

	// sm := session.NewSessionsJWT("golangcourseSessionSecret")
	// sm := session.NewSessionsJWTVer(cfg.Session.Secret, db)
	sessions := session.NewSessionsDB(db)

	userHandlers := &user.UserHandler{
		Tmpl:      pageTemplates,
		Sessions:  sessions,
		UsersRepo: usersRepo,
	}

	mux := http.NewServeMux()

	// mux.HandleFunc("/photos/", h.List)
	mux.HandleFunc("/photos/", appHandlers.ListGQL)

	mux.HandleFunc("/api/v1/photos/list", appHandlers.ListAPI)
	mux.HandleFunc("/api/v1/photos/upload", appHandlers.UploadAPI)
	mux.HandleFunc("/api/v1/photos/rate", appHandlers.RateAPI)

	mux.HandleFunc("/user/login", userHandlers.Login)
	mux.HandleFunc("/user/login_oauth", userHandlers.LoginOauth)
	mux.HandleFunc("/user/logout", userHandlers.Logout)
	mux.HandleFunc("/user/reg", userHandlers.Reg)
	mux.HandleFunc("/user/change_pass", userHandlers.ChangePassword)

	mux.HandleFunc("/api/v1/user/follow", userHandlers.FollowAPI)
	mux.HandleFunc("/api/v1/user/following", userHandlers.FollowingAPI)
	mux.HandleFunc("/api/v1/user/recomends", userHandlers.RecomendsAPI)

	mux.HandleFunc("/", Index)

	// START gqlgen part -------------------------------------------------------

	gqlResolver := &graphql.Resolver{
		PhotosRepo:  photosRepo,
		UsersRepo:   usersRepo,
		BlobStorage: s3Storage,
	}
	gqlCfg := graphql.Config{
		Resolvers: gqlResolver,
	}
	gqlHandler := gqlgenHandler.GraphQL(
		graphql.NewExecutableSchema(gqlCfg),
		gqlgenHandler.ComplexityLimit(500),
	)
	myGqlHandler := graphql.UserLoaderMiddleware(gqlResolver, gqlHandler)

	mux.Handle("/graphql", myGqlHandler)
	// TODO enable csrf for graphql after playground done
	mux.HandleFunc("/playground", gqlgenHandler.Playground("GraphQL playground", "/graphql"))

	// END gqlgen part ----------------------------------------------------------

	// отрабатывают в обратном добавлению порядке, те AuthMiddleware будет 1-м
	appHttpHandlers := token.CsrfTokenMiddleware(csrfTokens, mux)
	appHttpHandlers = session.AuthMiddleware(sessions, appHttpHandlers)
	http.Handle("/", appHttpHandlers)

	http.HandleFunc("/api/v1/internal/images/auth", userHandlers.InternalImagesAuth)

	imagesHandler := http.StripPrefix(
		"/images/",
		http.FileServer(http.Dir("./images")),
	)
	http.Handle("/images/", imagesHandler)

	http.Handle("/static/", http.FileServer(assets.Assets))

	f, _ := assets.Assets.Open("/static/favicon.ico")
	favicon, _ := ioutil.ReadAll(f)
	f.Close()
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Write(favicon)
	})

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
