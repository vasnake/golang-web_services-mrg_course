package middleware

import (
	"context"
	"net/http"

	"photolist/pkg/utils/randutils"
)

const REQUEST_ID_KEY = "requestID"

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		// create new RID
		if requestID == "" {
			// https://github.com/opentracing/specification/blob/master/rfc/trace_identifiers.md
			requestID = randutils.RandBytesHex(16)
			r.Header.Set("X-Request-ID", requestID)
			r.Header.Set("trace-id", requestID)
			w.Header().Set("X-Request-ID", requestID)
			w.Header().Set("trace-id", requestID)
		}

		ctx := context.WithValue(r.Context(), REQUEST_ID_KEY, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestIDFromContext(ctx context.Context) string {
	requestID, ok := ctx.Value(REQUEST_ID_KEY).(string)
	if !ok {
		return "unknown"
	}
	return requestID
}
