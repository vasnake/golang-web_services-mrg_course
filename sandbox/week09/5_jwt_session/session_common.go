package jwt_session

import (
	"context"
	"errors"
	"net/http"
)

// session data
type Session struct {
	UserID uint32
	ID     string
}

// session interface
type SessionManager interface {
	Check(*http.Request) (*Session, error)
	Create(http.ResponseWriter, *User) error
	DestroyCurrent(http.ResponseWriter, *http.Request) error
	DestroyAll(http.ResponseWriter, *User) error
}

// линтер ругается если используем базовые типы в Value контекста
// типа так безопаснее разграничивать
type ctxKey int

const sessionKey ctxKey = 1

var (
	ErrNoAuth = errors.New("No session found")
)

var (
	noAuthUrls = map[string]struct{}{
		"/user/login": struct{}{},
		"/user/reg":   struct{}{},
		"/":           struct{}{},
	}
)

// create session from session manager 'check' method, add session to context
func AuthMiddleware(sm SessionManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := noAuthUrls[r.URL.Path]; ok {
			next.ServeHTTP(w, r) // serve w/o auth
			return
		}

		sess, err := sm.Check(r)
		if err != nil {
			http.Error(w, "No auth", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), sessionKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx)) // serve with auth
	})
}

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(sessionKey).(*Session)
	if !ok {
		return nil, ErrNoAuth
	}
	return sess, nil
}
