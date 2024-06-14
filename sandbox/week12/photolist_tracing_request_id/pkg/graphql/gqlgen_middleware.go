package graphql

import (
	"context"
	"log"
	"net/http"
	"time"

	"photolist/pkg/middleware"
	"photolist/pkg/session"
	"photolist/pkg/user"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler/extension"
)

// go run github.com/vektah/dataloaden UserLoader uint32 *coursera/3p/photolist/100_gqlgen/main.User

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

// -----

func ResolverMiddleware(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	reqCtx := graphql.GetFieldContext(ctx)
	start := time.Now()
	res, err = next(ctx)
	requestID := middleware.RequestIDFromContext(ctx)
	log.Printf("[resolver] %s %s %v '%v'", requestID, time.Since(start), reqCtx.Path(), err)
	return
}

func RequestMiddleware(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	reqCtx := graphql.GetOperationContext(ctx)
	cs := extension.GetComplexityStats(ctx)
	start := time.Now()
	result := next(ctx)
	requestID := middleware.RequestIDFromContext(ctx)
	log.Printf("[RequestMiddleware] %s %s %s %d", requestID, time.Since(start), reqCtx.OperationName, cs.Complexity)
	return result
}
