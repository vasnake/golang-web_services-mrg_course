package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"photolist/pkg/assets"
	"photolist/pkg/blobstorage"
	"photolist/pkg/config"
	"photolist/pkg/graphql"
	"photolist/pkg/index"
	"photolist/pkg/middleware"
	"photolist/pkg/photos"
	"photolist/pkg/session"
	"photolist/pkg/templates"
	"photolist/pkg/token"
	"photolist/pkg/user"

	// "github.com/99designs/gqlgen-contrib/gqlopentracing"
	gqlgenHandler "github.com/99designs/gqlgen/handler"
	_ "github.com/go-sql-driver/mysql"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

	// jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

var (
	appName   string = "photolist"
	buildHash string = "unknown"
	buildTime string = "unknown"
)

func main() {
	log.Printf("[startup] %s, commit %s, build %s", appName, buildHash, buildTime)

	// from command line, using flag package
	// export OAUTH_APP_ID=Ov2***gJF
	// export OAUTH_APP_SECRET=ada***860
	// go run main.go -appid ${OAUTH_APP_ID:-foo} -appsecret ${OAUTH_APP_SECRET:-bar}
	flag.StringVar(&user.APP_ID, "appid", "unknown", "oauth app id (client id) from github registered app")
	flag.StringVar(&user.APP_SECRET, "appsecret", "unknown", "oauth app secret (client key) from github registered app")
	flag.Parse()
	show("you must not show this! appid, appsecret from flag package: ", user.APP_ID, user.APP_SECRET)

	cfg := &config.Config{}
	viperSvc, err := config.Read(appName, config.Defaults, cfg)
	if err != nil {
		log.Fatalf("can't read config: %v\n", err)
	}
	log.Printf("config, example.env2: %s", viperSvc.GetString("example.env2"))

	// from secrets.env file
	oauthAppId := viperSvc.GetString("OAUTH_APP_ID")
	oauthAppSecret := viperSvc.GetString("OAUTH_APP_SECRET")
	if oauthAppId != "" && user.APP_ID == "unknown" {
		user.APP_ID = oauthAppId
	}
	if oauthAppSecret != "" && user.APP_SECRET == "unknown" {
		user.APP_SECRET = oauthAppSecret
	}
	show("you must not show this! appid, appsecret from viper package: ", user.APP_ID, user.APP_SECRET)

	if user.APP_ID == "" || user.APP_ID == "unknown" || user.APP_SECRET == "" || user.APP_SECRET == "unknown" {
		show("OAuth functionality broken, OAuth tokens not provided")
	}

	listenAddr := viperSvc.GetString("http.port")
	log.Printf("tcp listen addr (config http.port): '%s'", listenAddr)

	db, closeDBFunc := openDBConnect(cfg)
	defer closeDBFunc()

	closeTracerFunc := setTracer(appName, buildHash, buildTime, viperSvc.GetString("JAEGER_AGENT_ADDR"))
	defer closeTracerFunc()

	csrfTokens, err := token.NewJwtToken(cfg.Token.Secret)
	if err != nil {
		log.Fatalf("can't init tokens: %v\n", err)
	}

	pageTemplates := templates.NewTemplates(assets.Assets, csrfTokens)

	blobStorage, err := blobstorage.NewS3Storage(cfg.S3.Host, cfg.S3.Access, cfg.S3.Secret, cfg.S3.Bucket)
	if err != nil {
		log.Fatalf("can't create s3 blobstorage %v\n", err)
	}

	photosRepo := photos.NewPhotosRepository(db)
	usersRepo := user.NewUsersRepository(db)
	sessionSvc := session.NewSessionsDB(db)

	appSvc := &photos.PhotolistHandler{
		UsersRepo:   usersRepo,
		PhotosRepo:  photosRepo,
		Tmpl:        pageTemplates,
		BlobStorage: blobStorage,
	}

	userSvc := &user.UserHandler{
		Tmpl:      pageTemplates,
		Sessions:  sessionSvc,
		UsersRepo: usersRepo,
	}

	gqlResolver := &graphql.Resolver{
		PhotosRepo:  photosRepo,
		UsersRepo:   usersRepo,
		BlobStorage: blobStorage,
	}
	gqlHttpHandlerFunc := gqlgenHandler.GraphQL(
		graphql.NewExecutableSchema(graphql.Config{Resolvers: gqlResolver}),
		gqlgenHandler.ComplexityLimit(500),
		gqlgenHandler.RequestMiddleware(graphql.RequestMiddleware),
		gqlgenHandler.ResolverMiddleware(graphql.ResolverMiddleware),
		gqlgenHandler.Tracer(graphql.NewTracer()),
	)
	gqlHttpHandler := graphql.UserLoaderMiddleware(gqlResolver, gqlHttpHandlerFunc)

	mux := setupRoutesMultiplexer(appSvc, userSvc, gqlHttpHandler)

	// middleware
	// отрабатывают в обратном добавлению порядке, те RequestIDMiddleware будет 1-м
	httpRootHandler := token.CsrfTokenMiddleware(csrfTokens, mux)

	httpRootHandler = session.AuthMiddleware(sessionSvc, httpRootHandler)
	httpRootHandler = middleware.AccessLog(httpRootHandler)
	httpRootHandler = middleware.RequestIDMiddleware(httpRootHandler)

	http.Handle("/", httpRootHandler)
	http.Handle("/static/", http.FileServer(assets.Assets))
	handleFavicon(assets.Assets)

	// don't need this, separate svc exists: photoauth
	http.HandleFunc("/api/v1/internal/images/auth", userSvc.InternalImagesAuth)

	// don't need this, nginx works fine
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))

	log.Printf("[startup] listening server at %s", listenAddr)
	http.ListenAndServe(listenAddr, nil)
}

var setupRoutesMultiplexer = func(appSvc *photos.PhotolistHandler, userSvc *user.UserHandler, gqlHttpHandler http.Handler) *http.ServeMux {
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

	mux.HandleFunc("/", index.Index)

	{ // START gqlgen part
		mux.Handle("/graphql", gqlHttpHandler)
		mux.Handle("/graphql/", gqlHttpHandler)
		// TODO enable csrf for graphql after playground done
		mux.HandleFunc("/playground", gqlgenHandler.Playground("GraphQL playground", "/graphql"))
	} // END gqlgen part

	return mux
}

var handleFavicon = func(fs http.FileSystem) {
	file, err := fs.Open("/static/favicon.ico")
	if err == nil {
		faviconBytes, err := io.ReadAll(file)
		file.Close()
		if err == nil {
			http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
				w.Write(faviconBytes)
			})
		} else {
			log.Printf("favicon, io ReadAll failed: %w", err)
		}
	} else {
		log.Printf("favicon, Assets Open failed: %w", err)
	}

	return
}

var setTracer = func(appName, buildHash, buildTime, agentAddr string) func() {
	cfg := jaegercfg.Configuration{
		ServiceName: appName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: agentAddr,
		},
		Tags: []opentracing.Tag{
			{Key: "buildHash", Value: buildHash},
			{Key: "buildTime", Value: buildTime},
		},
	}

	tracer, closer, err := cfg.NewTracer(
		// jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
	if err != nil {
		log.Fatalf("jaeger NewTracer failed: %v\n", err)
	}

	opentracing.SetGlobalTracer(tracer)

	return func() { closer.Close() }
}

var openDBConnect = func(cfg *config.Config) (*sql.DB, func()) {
	// основные настройки к базе
	// mysqlDsn := "root:love@tcp(127.0.0.1:3306)/photolist?charset=utf8&interpolateParams=true"
	mysqlDsn := "%s:%s@tcp(%s)/%s?charset=utf8&interpolateParams=true"
	mysqlDsn = fmt.Sprintf(mysqlDsn, cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Database)
	log.Printf("mysql DSN: %s", mysqlDsn)
	db, err := sql.Open("mysql", mysqlDsn)
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalf("can't connect to db: %v\n", err)
	}

	return db, func() { db.Close() }
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
