package session

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"
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

func (sm *SessionsGRPC) Check(r *http.Request) (*Session, error) {
	sessionCookie, err := r.Cookie(cookieName)
	if err == http.ErrNoCookie {
		log.Println("CheckSession no cookie")
		return nil, ErrNoAuth
	}

	authSess, err := sm.client.Check(context.Background(), &AuthCheckIn{SessKey: sessionCookie.Value})
	if err != nil {
		return nil, err
	}

	return &Session{
		ID:     authSess.GetID(),
		UserID: authSess.GetUserID(),
	}, nil
}

func (sm *SessionsGRPC) Create(w http.ResponseWriter, user UserInterface) error {
	authSess, err := sm.client.Create(context.Background(), &AuthUserIn{
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

func (sm *SessionsGRPC) DestroyCurrent(w http.ResponseWriter, r *http.Request) error {
	cookie := http.Cookie{
		Name:    cookieName,
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	_, err := sm.client.DestroyCurrent(context.Background(), &AuthSession{
		ID: cookie.Value,
	})
	return err
}

func (sm *SessionsGRPC) DestroyAll(w http.ResponseWriter, user UserInterface) error {
	cookie := http.Cookie{
		Name:    cookieName,
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	_, err := sm.client.DestroyAll(context.Background(), &AuthUserIn{
		UserID: user.GetID(),
		Ver:    user.GetVer(),
	})
	return err
}
