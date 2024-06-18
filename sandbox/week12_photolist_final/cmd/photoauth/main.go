package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"photolist/pkg/config"
	"photolist/pkg/middleware"
	"photolist/pkg/session"
	"photolist/pkg/user"

	_ "github.com/go-sql-driver/mysql"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	// jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

var (
	appName   string = "photoauth"
	buildHash string = "unknown"
	buildTime string = "unknown"
)

func main() {
	log.Printf("[startup] %s, commit %s, build %s", appName, buildHash, buildTime)

	cfg := &config.Config{}
	viperSvc, err := config.Read(appName, config.Defaults, cfg)
	if err != nil {
		log.Fatalf("can't read config: %v\n", err)
	}

	log.Printf("[startup] http.port: '%#v', example.env1 '%#v', example.env2 '%#v'",
		cfg.HTTP.Port,
		viperSvc.GetString("example.env1"),
		viperSvc.GetString("example.env2"),
	)

	grpcServerAddr := viperSvc.GetString("session.grpc_addr")
	log.Println("grpc server addr:", grpcServerAddr)

	listenAddr := viperSvc.GetString("http.port")
	log.Printf("listen addr (config http.port): '%s'", listenAddr)

	// основные настройки к базе
	// mysqlDsn := "root:love@tcp(host.docker.internal:3306)/photolist?charset=utf8&interpolateParams=true"
	mysqlDsn := "%s:%s@tcp(%s)/%s?charset=utf8&interpolateParams=true"
	mysqlDsn = fmt.Sprintf(mysqlDsn, cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Database)
	log.Printf("mysql DSN: %s", mysqlDsn)

	db, err := sql.Open("mysql", mysqlDsn)
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalf("[startup] can't connect to db: %v\n", err)
	}
	defer db.Close()

	// start tracing cfg --------------------------------------------------------------------------
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: appName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: viperSvc.GetString("JAEGER_AGENT_ADDR"),
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
	// end tracing cfg ----------------------------------------------------------------------------

	grpcSessions, err := session.NewSessionsGRPC(grpcServerAddr)
	if err != nil {
		log.Fatalf("[startup] can't connect to grpc svc: %v\n", err)
	}

	usersRepo := user.NewUsersRepository(db)

	userSvc := &user.UserHandler{
		Tmpl:      nil,
		Sessions:  grpcSessions,
		UsersRepo: usersRepo,
	}

	httpHandler := middleware.AccessLog(http.HandlerFunc(userSvc.InternalImagesAuth))
	httpHandler = middleware.RequestIDMiddleware(httpHandler)

	http.Handle("/api/v1/internal/images/auth", httpHandler)

	log.Printf("[startup] listening server at %s", listenAddr)
	http.ListenAndServe(listenAddr, nil)
}
