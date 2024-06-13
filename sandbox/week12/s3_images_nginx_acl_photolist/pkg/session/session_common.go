package session

import (
	"context"
	"errors"
	"net/http"
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
	Check(*http.Request) (*Session, error)
	Create(http.ResponseWriter, UserInterface) error
	DestroyCurrent(http.ResponseWriter, *http.Request) error
	DestroyAll(http.ResponseWriter, UserInterface) error
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
		if _, ok := noAuthUrls[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}
		sess, err := sm.Check(r)
		if err != nil {
			http.Error(w, "No auth", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), sessionKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
