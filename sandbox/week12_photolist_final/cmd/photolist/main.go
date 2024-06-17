package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
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
	buildHash string = "_dev"
	buildTime string = "_dev"
)

func main() {
	log.Printf("[startup] %s, commit %s, build %s", appName, buildHash, buildTime)
	rand.Seed(time.Now().UnixNano())

	cfg := &config.Config{}
	v1, err := config.Read(appName, config.Defaults, cfg)
	log.Println(v1, cfg, err, v1.GetString("example.env2"))

	// основные настройки к базе
	// dsn := "root:love@tcp(127.0.0.1:3306)/photolist?charset=utf8&interpolateParams=true"
	dsn := "%s:%s@tcp(%s)/%s?charset=utf8&interpolateParams=true"
	dsn = fmt.Sprintf(dsn, cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Database)
	db, err := sql.Open("mysql", dsn)
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalf("cant connect to db, err: %v\n", err)
	}

	// log.Println("JAEGER_AGENT_HOST", )

	// start tracing cfg
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: appName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: v1.GetString("JAEGER_AGENT_ADDR"),
		},
		Tags: []opentracing.Tag{
			{Key: "buildHash", Value: buildHash},
			{Key: "buildTime", Value: buildTime},
		},
	}

	tracer, closer, err := jaegerCfgInstance.NewTracer(
		// jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	// end tracing cfg

	// tokens, err := token.NewHMACHashToken("golangcourse")
	// tokens, err := token.NewAesCryptHashToken("qsRY2e4hcM5T7X984E9WQ5uZ8Nty7fxB")
	tokens, err := token.NewJwtToken(cfg.Token.Secret)
	tmpls := templates.NewTemplates(assets.Assets, tokens)
	if err != nil {
		log.Fatalf("cant init tokens: %v\n", err)
	}

	// storage := blobstorage.NewFSStorage("./images/")
	storage, err := blobstorage.NewS3Storage(cfg.S3.Host,
		cfg.S3.Access, cfg.S3.Secret,
		cfg.S3.Bucket)
	if err != nil {
		log.Fatalln("cant creat s3 blobstorage", err)
	}

	photosRepo := photos.NewPhotosRepository(db)
	usersRepo := user.NewUsersRepository(db)

	h := &photos.PhotolistHandler{
		UsersRepo:   usersRepo,
		PhotosRepo:  photosRepo,
		Tmpl:        tmpls,
		BlobStorage: storage,
	}

	sm := session.NewSessionsDB(db)
	// sm := session.NewSessionsJWT("golangcourseSessionSecret")
	// sm := session.NewSessionsJWTVer(cfg.Session.Secret, db)

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

	{ // START gqlgen part
		resolver := &graphql.Resolver{
			PhotosRepo:  photosRepo,
			UsersRepo:   usersRepo,
			BlobStorage: storage,
		}
		gqlCfg := graphql.Config{
			Resolvers: resolver,
		}
		gqlHandler := gqlgenHandler.GraphQL(
			graphql.NewExecutableSchema(gqlCfg),
			gqlgenHandler.ComplexityLimit(500),
			gqlgenHandler.RequestMiddleware(graphql.RequestMiddleware),   // каждый запрос после парсинга
			gqlgenHandler.ResolverMiddleware(graphql.ResolverMiddleware), // каждый вызлв ресолвера
			gqlgenHandler.Tracer(graphql.NewTracer()),
		)
		myGqlHandler := graphql.UserLoaderMiddleware(resolver, gqlHandler)

		mux.Handle("/graphql", myGqlHandler)
		mux.Handle("/graphql/", myGqlHandler)
		// TODO enable csrf for graphql after playground done
		mux.HandleFunc("/playground", gqlgenHandler.Playground("GraphQL playground", "/graphql"))
	} // END gqlgen part

	// отрабатывают в обратном добавлению порядке, те AuthMiddleware будет 1-м
	handlers := token.CsrfTokenMiddleware(tokens, mux)
	handlers = session.AuthMiddleware(sm, handlers)
	handlers = middleware.AccessLog(handlers)
	handlers = middleware.RequestIDMiddleware(handlers)

	http.Handle("/", handlers)

	http.HandleFunc("/api/v1/internal/images/auth", u.InternalImagesAuth)

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

	listenAddr := ":" + v1.GetString("http.port")
	log.Printf("[startup] listening server at %s", listenAddr)
	http.ListenAndServe(listenAddr, nil)
}
