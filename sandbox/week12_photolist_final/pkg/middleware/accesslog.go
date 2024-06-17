package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
)

func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		span, newCtx := opentracing.StartSpanFromContext(ctx, r.URL.Path)
		defer span.Finish()

		requestID := RequestIDFromContext(ctx)
		span.SetTag("myrequestid", requestID)

		start := time.Now()
		r = r.WithContext(newCtx)
		next.ServeHTTP(w, r)

		log.Printf("[access] %s %s %s %s %s", requestID, time.Since(start), r.RemoteAddr, r.Method, r.URL.Path)
	})
}
