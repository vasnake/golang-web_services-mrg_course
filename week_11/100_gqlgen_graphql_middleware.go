package main

import (
	"context"
	"net/http"
	"time"
)

// go run github.com/vektah/dataloaden UserLoader uint32 *coursera/3p/photolist/100_gqlgen/main.User

func UserLoaderMiddleware(resolver *Resolver, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := UserLoaderConfig{
			MaxBatch: 100,
			Wait:     1 * time.Millisecond,
			Fetch: func(ids []uint32) ([]*User, []error) {
				sess, _ := SessionFromContext(r.Context())
				return resolver.UsersRepo.LookupByIDs(sess.UserID, ids)
			},
		}
		userLoader := NewUserLoader(cfg)
		ctx := context.WithValue(r.Context(), "userLoaderKey", userLoader)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func UserLoaderFromContext(ctx context.Context) *UserLoader {
	return ctx.Value("userLoaderKey").(*UserLoader)
}
