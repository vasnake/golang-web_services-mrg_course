package session

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"photolist/pkg/utils/randutils"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	_ SessionManager = (*SessionsJWTVer)(nil)
)

type SessionsJWTVer struct {
	Secret []byte
	DB     *sql.DB
}

type SessionJWTVerClaims struct {
	UserID uint32 `json:"uid"`
	Ver    int32  `json:"ver,omitempty"`
	jwt.StandardClaims
}

func NewSessionsJWTVer(secret string, db *sql.DB) *SessionsJWTVer {
	return &SessionsJWTVer{
		Secret: []byte(secret),
		DB:     db,
	}
}

func (sm *SessionsJWTVer) parseSecretGetter(token *jwt.Token) (interface{}, error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok || method.Alg() != "HS256" {
		return nil, fmt.Errorf("bad sign method")
	}
	return sm.Secret, nil
}

func (sm *SessionsJWTVer) Check(ctx context.Context, r *http.Request) (*Session, error) {
	sessionCookie, err := r.Cookie(cookieName)
	if err == http.ErrNoCookie {
		log.Println("CheckSession no cookie")
		return nil, ErrNoAuth
	}

	payload := &SessionJWTVerClaims{}
	_, err = jwt.ParseWithClaims(sessionCookie.Value, payload, sm.parseSecretGetter)
	if err != nil {
		return nil, fmt.Errorf("cant parse jwt token: %v", err)
	}
	// проверка exp, iat
	if payload.Valid() != nil {
		return nil, fmt.Errorf("invalid jwt token: %v", err)
	}

	var ver int32
	row := sm.DB.QueryRow(`SELECT ver FROM users WHERE id = ?`, payload.UserID)
	err = row.Scan(&ver)
	if err == sql.ErrNoRows {
		log.Println("CheckSession no rows")
		return nil, ErrNoAuth
	} else if err != nil {
		log.Println("CheckSession err:", err)
		return nil, err
	}

	if payload.Ver != ver {
		log.Println("CheckSession invalid version, sess", payload.Ver, "user", ver)
		return nil, ErrNoAuth
	}

	return &Session{
		ID:     payload.Id,
		UserID: payload.UserID,
	}, nil
}

func (sm *SessionsJWTVer) Create(ctx context.Context, w http.ResponseWriter, user UserInterface) error {
	data := SessionJWTVerClaims{
		UserID: user.GetID(),
		Ver:    user.GetVer(), // изменилось по сравнению со stateless-сессией
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(90 * 24 * time.Hour).Unix(), // 90 days
			IssuedAt:  time.Now().Unix(),
			Id:        randutils.RandStringRunes(32),
		},
	}
	sessVal, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, data).SignedString(sm.Secret)

	cookie := &http.Cookie{
		Name:    cookieName,
		Value:   sessVal,
		Expires: time.Now().Add(90 * 24 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	return nil
}

func (sm *SessionsJWTVer) DestroyCurrent(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	cookie := http.Cookie{
		Name:    cookieName,
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)

	// но!
	// если куку украли - ее не отозвать
	// ¯ \ _ (ツ) _ / ¯

	return nil
}

func (sm *SessionsJWTVer) DestroyAll(ctx context.Context, w http.ResponseWriter, user UserInterface) error {
	// но!
	// мы никак не можем дотянуться до других сессий
	// ¯ \ _ (ツ) _ / ¯
	return nil
}
