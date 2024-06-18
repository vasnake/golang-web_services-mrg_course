package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"photolist/pkg/config"
	"photolist/pkg/session"
	"photolist/pkg/utils/traceutils"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	// jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

var (
	appName   string = "auth"
	buildHash string = "unknown"
	buildTime string = "unknown"
)

func main() {
	log.Printf("[startup] %s, commit %s, build %s", appName, buildHash, buildTime)

	cfg := &config.Config{}
	viperSvc, err := config.Read(appName, config.Defaults, cfg)
	if err != nil {
		log.Fatalf("[startup] can't read config: %v\n", err)
	}

	listenAddr := viperSvc.GetString("service.port")
	log.Printf("tcp listen addr (config service.port): %s", listenAddr)

	// start tracing cfg -----------------------------------------------------
	jaegerCfg := jaegercfg.Configuration{
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

	tracer, closer, err := jaegerCfg.NewTracer(
		// jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	// end tracing cfg -------------------------------------------------------

	// основные настройки к базе
	dsn := "%s:%s@tcp(%s)/%s?charset=utf8&interpolateParams=true"
	dsn = fmt.Sprintf(dsn, cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Database)
	log.Printf("mysql DSN: %s", dsn)
	db, err := sql.Open("mysql", dsn)
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalf("[startup] can't connect to db, err: %v\n", err)
	}
	defer db.Close()

	authSvc := &session.AuthService{
		DB: db,
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(LoggingInterceptor),
	)

	session.RegisterAuthServer(grpcServer, authSvc)

	lstnr, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalln("[startup] can't listen port", err)
	}

	log.Printf("[startup] listening server at %s", listenAddr)
	grpcServer.Serve(lstnr)
}

func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	reqHandlerFunc grpc.UnaryHandler,
) (interface{}, error) {

	start := time.Now()

	var requestID string
	md, mdExists := metadata.FromIncomingContext(ctx)
	if mdExists {
		requestID = md["x-request-id"][0]
	} else {
		requestID = "unknown"
	}

	var span opentracing.Span
	spanCtx, err := opentracing.GlobalTracer().
		Extract(opentracing.HTTPHeaders, traceutils.MetadataReaderWriter{MD: md})
	if err == nil {
		span = opentracing.StartSpan(info.FullMethod, ext.RPCServerOption(spanCtx))
	} else {
		span = opentracing.StartSpan(info.FullMethod)
	}
	defer span.Finish()

	// do some useful work already
	reply, err := reqHandlerFunc(ctx, req)

	log.Printf("[access] %s %s %s '%v'", requestID, time.Since(start), info.FullMethod, err)
	return reply, err
}
