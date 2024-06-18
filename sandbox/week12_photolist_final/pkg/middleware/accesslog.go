package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
)

func AccessLog(httpHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ctx := r.Context()
		span, spanCtx := opentracing.StartSpanFromContext(ctx, r.URL.Path)
		defer span.Finish()

		requestID := RequestIDFromContext(ctx)
		span.SetTag("myrequestid", requestID)

		req := r.WithContext(spanCtx)

		httpHandler.ServeHTTP(w, req)

		log.Printf("[access] %s %s %s %s %s", requestID, time.Since(start), req.RemoteAddr, req.Method, req.URL.Path)
	})
}
