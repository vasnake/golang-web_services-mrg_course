package session

import (
	"context"
	"errors"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Session struct {
	UserID uint32
	ID     string
}

/*
	https://medium.com/@dotronglong/interface-naming-convention-in-golang-f53d9f471593
	user
	User
	Userer
	IUser
	UserI
	UserInterface
*/

type UserInterface interface {
	GetID() uint32
	GetVer() int32
}

type SessionManager interface {
	Check(context.Context, *http.Request) (*Session, error)
	Create(context.Context, http.ResponseWriter, UserInterface) error
	DestroyCurrent(context.Context, http.ResponseWriter, *http.Request) error
	DestroyAll(context.Context, http.ResponseWriter, UserInterface) error
}

// линтер ругается если используем базовые типы в Value контекста
// типа так безопаснее разграничивать
type ctxKey int

const sessionKey ctxKey = 1

const cookieName = "session_id"

var (
	ErrNoAuth = errors.New("No session found")
)

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(sessionKey).(*Session)
	if !ok {
		return nil, ErrNoAuth
	}
	return sess, nil
}

var (
	noAuthUrls = map[string]struct{}{
		"/user/login_oauth": struct{}{},
		"/user/login":       struct{}{},
		"/user/reg":         struct{}{},
		"/":                 struct{}{},
	}
)

func AuthMiddleware(sm SessionManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if _, ok := noAuthUrls[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}
		span, _ := opentracing.StartSpanFromContext(ctx, "auth")
		ext.SpanKind.Set(span, "server")
		ext.Component.Set(span, "auth")
		sess, err := sm.Check(ctx, r)
		span.Finish()
		if err != nil {
			http.Error(w, "No auth", http.StatusUnauthorized)
			return
		}
		ctx = context.WithValue(ctx, sessionKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
