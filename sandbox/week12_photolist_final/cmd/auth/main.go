package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
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
	buildHash string = "_dev"
	buildTime string = "_dev"
)

func AccessLogInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	var requestID string
	md, mdExists := metadata.FromIncomingContext(ctx)
	if mdExists {
		requestID = md["x-request-id"][0]
	} else {
		requestID = "-"
	}

	clientContext, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, traceutils.MetadataReaderWriter{md})
	var serverSpan opentracing.Span
	if err == nil {
		serverSpan = opentracing.StartSpan(info.FullMethod, ext.RPCServerOption(clientContext))
	} else {
		serverSpan = opentracing.StartSpan(info.FullMethod)
	}
	defer serverSpan.Finish()

	reply, err := handler(ctx, req)

	log.Printf("[access] %s %s %s '%v'", requestID, time.Since(start), info.FullMethod, err)
	return reply, err
}

func main() {
	log.Printf("[startup] %s, commit %s, build %s", appName, buildHash, buildTime)
	rand.Seed(time.Now().UnixNano())

	cfg := &config.Config{}
	v1, err := config.Read(appName, config.Defaults, cfg)
	if err != nil {
		log.Fatalf("[startup] cant read config, err: %v\n", err)
	}

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

	// основные настройки к базе
	dsn := "%s:%s@tcp(%s)/%s?charset=utf8&interpolateParams=true"
	dsn = fmt.Sprintf(dsn, cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Database)
	db, err := sql.Open("mysql", dsn)
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalf("[startup] cant connect to db, err: %v\n", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(AccessLogInterceptor),
	)
	svc := &session.AuthService{
		DB: db,
	}
	session.RegisterAuthServer(server, svc)

	listenAddr := v1.GetString("service.port")
	lis, err := net.Listen("tcp", ":"+listenAddr)
	if err != nil {
		log.Fatalln("[startup] cant listen port", err)
	}
	log.Printf("[startup] listening server at %s", listenAddr)
	server.Serve(lis)
}
