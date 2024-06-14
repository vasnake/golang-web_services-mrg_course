package middleware

import (
	"log"
	"net/http"
	"time"
)

func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		requestID := RequestIDFromContext(r.Context())
		log.Printf("[access] %s %s %s %s %s", requestID, time.Since(start), r.RemoteAddr, r.Method, r.URL.Path)
	})
}
