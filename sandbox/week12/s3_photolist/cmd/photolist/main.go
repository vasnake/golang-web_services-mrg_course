package s3_photolist_main

import (
	"database/sql"
	"flag"
	"fmt"
	ioutil "io" // "io/ioutil"
	"log"
	"net/http"
	"time"

	gqlgenHandler "github.com/99designs/gqlgen/handler"
	_ "github.com/go-sql-driver/mysql"

	"week12/s3_photolist/pkg/assets"
	"week12/s3_photolist/pkg/blobstorage"
	"week12/s3_photolist/pkg/graphql"
	"week12/s3_photolist/pkg/photos"
	"week12/s3_photolist/pkg/session"
	"week12/s3_photolist/pkg/templates"
	"week12/s3_photolist/pkg/token"
	"week12/s3_photolist/pkg/user"
)

var (
	// go build -ldflags "-X main.buildHash=${COMMIT} -X main.buildTime=${BUILD_TIME}" -o ./bin/photolist
	appName   string = "photolist"
	buildHash string = "_dev"
	buildTime string = "_dev"
)

func MainDemo() {
	flag.StringVar(&user.APP_ID, "appid", "foo?", "oauth app id (client id) from github registered app")
	flag.StringVar(&user.APP_SECRET, "appsecret", "bar?", "oauth app secret (client key) from github registered app")
	flag.Parse()
	show("you must not show this! appid, appsecret: ", user.APP_ID, user.APP_SECRET)

	log.Printf("starting %s, commit %s, build %s", appName, buildHash, buildTime)

	// rand.Seed(time.Now().UnixNano())

	// основные настройки к базе
	// sqlDsn := "root:love@tcp(127.0.0.1:3306)/photolist?charset=utf8&interpolateParams=true"
	// sqlDsn := "root:love@tcp(host.docker.internal:3306)/photolist?charset=utf8&interpolateParams=true"
	sqlDsn := "root@tcp(127.0.0.1:3306)/photolist?charset=utf8&interpolateParams=true"
	db, err := sql.Open("mysql", sqlDsn)
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalf("can't connect to db, err: %v\n", err)
	}
	defer db.Close()

	// csrfTokens, err := token.NewHMACHashToken("golangcourse")
	// csrfTokens, err := token.NewAesCryptHashToken("qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB")
	csrfTokens, err := token.NewJwtToken("qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB")

	tmplates := templates.NewTemplates(assets.Assets, csrfTokens)
	if err != nil {
		log.Fatalf("can't init tokens: %v\n", err)
	}

	// blobStorage := blobstorage.NewFSStorage("./images/")
	// blobStorage, err := blobstorage.NewS3Storage("http://minio:9000",
	blobStorage, err := blobstorage.NewS3Storage(
		"http://127.0.0.1:9000",
		"access_123",
		"secret_123",
		"photolist",
	)
	if err != nil {
		log.Fatalln("can't creat s3 blobstorage", err)
	}

	photosRepo := photos.NewPhotosRepository(db)
	usersRepo := user.NewUsersRepository(db)

	h := &photos.PhotolistHandler{
		UsersRepo:   usersRepo,
		PhotosRepo:  photosRepo,
		Tmpl:        tmplates,
		BlobStorage: blobStorage,
	}

	// userSessions := session.NewSessionsDB(db)
	// userSessions := session.NewSessionsJWT("golangcourseSessionSecret")
	userSessions := session.NewSessionsJWTVer("golangcourseSessionSecret", db)

	u := &user.UserHandler{
		Tmpl:      tmplates,
		Sessions:  userSessions,
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

	mux.HandleFunc("/", Index)

	// START gqlgen part
	gqlResolver := &graphql.Resolver{
		PhotosRepo:  photosRepo,
		UsersRepo:   usersRepo,
		BlobStorage: blobStorage,
	}
	cfg := graphql.Config{
		Resolvers: gqlResolver,
	}
	gqlHandler := gqlgenHandler.GraphQL(
		graphql.NewExecutableSchema(cfg),
		gqlgenHandler.ComplexityLimit(500),
	)
	myGqlHandler := graphql.UserLoaderMiddleware(gqlResolver, gqlHandler)

	mux.Handle("/graphql", myGqlHandler)
	// TODO enable csrf for graphql after playground done
	mux.HandleFunc("/playground", gqlgenHandler.Playground("GraphQL playground", "/graphql"))
	// END gqlgen part

	// отрабатывают в обратном добавлению порядке, те AuthMiddleware будет 1-м
	httpHandlers := token.CsrfTokenMiddleware(csrfTokens, mux)
	httpHandlers = session.AuthMiddleware(userSessions, httpHandlers)

	http.Handle("/", httpHandlers)

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

	listenAddr := ":8080"
	log.Printf("starting listening server at %s", listenAddr)
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
