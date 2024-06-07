package graphql

import (
	"context"
	"net/http"
	"time"

	"photolist/pkg/session"
	"photolist/pkg/user"
)

// go run github.com/vektah/dataloaden UserLoader uint32 *photolist.User

func UserLoaderMiddleware(resolver *Resolver, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := user.UserLoaderConfig{
			MaxBatch: 100,
			Wait:     1 * time.Millisecond,
			Fetch: func(ids []uint32) ([]*user.User, []error) {
				sess, _ := session.SessionFromContext(r.Context())
				return resolver.UsersRepo.LookupByIDs(sess.UserID, ids)
			},
		}
		userLoader := user.NewUserLoader(cfg)
		ctx := context.WithValue(r.Context(), "userLoaderKey", userLoader)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func UserLoaderFromContext(ctx context.Context) *user.UserLoader {
	return ctx.Value("userLoaderKey").(*user.UserLoader)
}
