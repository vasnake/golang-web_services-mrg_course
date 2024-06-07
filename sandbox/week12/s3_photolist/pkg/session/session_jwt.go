package session

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"photolist/pkg/utils/randutils"

	jwt "github.com/golang-jwt/jwt"
)

type SessionsJWT struct {
	Secret []byte
}

type SessionJWTClaims struct {
	UserID uint32 `json:"uid"`
	jwt.StandardClaims
}

func NewSessionsJWT(secret string) *SessionsJWT {
	return &SessionsJWT{
		Secret: []byte(secret),
	}
}

func (sm *SessionsJWT) parseSecretGetter(token *jwt.Token) (interface{}, error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok || method.Alg() != "HS256" {
		return nil, fmt.Errorf("bad sign method")
	}
	return sm.Secret, nil
}

func (sm *SessionsJWT) Check(r *http.Request) (*Session, error) {
	sessionCookie, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		log.Println("CheckSession no cookie")
		return nil, ErrNoAuth
	}

	payload := &SessionJWTClaims{}
	_, err = jwt.ParseWithClaims(sessionCookie.Value, payload, sm.parseSecretGetter)
	if err != nil {
		return nil, fmt.Errorf("cant parse jwt token: %v", err)
	}
	// проверка exp, iat
	if payload.Valid() != nil {
		return nil, fmt.Errorf("invalid jwt token: %v", err)
	}

	return &Session{
		ID:     payload.Id,
		UserID: payload.UserID,
	}, nil
}

func (sm *SessionsJWT) Create(w http.ResponseWriter, user UserInterface) error {
	data := SessionJWTClaims{
		UserID: user.GetID(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(90 * 24 * time.Hour).Unix(), // 90 days
			IssuedAt:  time.Now().Unix(),
			Id:        randutils.RandStringRunes(32),
		},
	}
	sessVal, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, data).SignedString(sm.Secret)

	cookie := &http.Cookie{
		Name:    "session",
		Value:   sessVal,
		Expires: time.Now().Add(90 * 24 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	return nil
}

func (sm *SessionsJWT) DestroyCurrent(w http.ResponseWriter, r *http.Request) error {
	cookie := http.Cookie{
		Name:    "session",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)

	// но!
	// если куку украли - ее не отозвать
	// ¯ \ _ (ツ) _ / ¯

	return nil
}

func (sm *SessionsJWT) DestroyAll(w http.ResponseWriter, user UserInterface) error {
	// но!
	// мы никак не можем дотянуться до других сессий
	// ¯ \ _ (ツ) _ / ¯
	return nil
}
