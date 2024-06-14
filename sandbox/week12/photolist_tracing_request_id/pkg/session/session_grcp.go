package session

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"photolist/pkg/middleware"
)

var (
	_ SessionManager = (*SessionsGRPC)(nil)
)

type SessionsGRPC struct {
	client AuthClient
}

func NewSessionsGRPC(addr string) (*SessionsGRPC, error) {
	grcpConn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("cant connect to grpc")
	}
	return &SessionsGRPC{
		client: NewAuthClient(grcpConn),
	}, nil
}

func ctxWithRID(ctx context.Context) context.Context {
	requestID := middleware.RequestIDFromContext(ctx)
	return metadata.AppendToOutgoingContext(ctx, "X-Request-ID", requestID)
}

func (sm *SessionsGRPC) Check(ctx context.Context, r *http.Request) (*Session, error) {
	sessionCookie, err := r.Cookie(cookieName)
	if err == http.ErrNoCookie {
		log.Println("CheckSession no cookie")
		return nil, ErrNoAuth
	}

	authSess, err := sm.client.Check(ctxWithRID(ctx), &AuthCheckIn{SessKey: sessionCookie.Value})
	if err != nil {
		return nil, err
	}

	return &Session{
		ID:     authSess.GetID(),
		UserID: authSess.GetUserID(),
	}, nil
}

func (sm *SessionsGRPC) Create(ctx context.Context, w http.ResponseWriter, user UserInterface) error {
	authSess, err := sm.client.Create(ctxWithRID(ctx), &AuthUserIn{
		UserID: user.GetID(),
		Ver:    user.GetVer(),
	})
	if err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:    cookieName,
		Value:   authSess.GetID(),
		Expires: time.Now().Add(90 * 24 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	return nil
}

func (sm *SessionsGRPC) DestroyCurrent(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	cookie := http.Cookie{
		Name:    cookieName,
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	_, err := sm.client.DestroyCurrent(ctxWithRID(ctx), &AuthSession{
		ID: cookie.Value,
	})
	return err
}

func (sm *SessionsGRPC) DestroyAll(ctx context.Context, w http.ResponseWriter, user UserInterface) error {
	cookie := http.Cookie{
		Name:    cookieName,
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	_, err := sm.client.DestroyAll(ctxWithRID(ctx), &AuthUserIn{
		UserID: user.GetID(),
		Ver:    user.GetVer(),
	})
	return err
}
