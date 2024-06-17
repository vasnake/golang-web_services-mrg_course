package token

import (
	"fmt"
	"time"

	"photolist/pkg/session"

	jwt "github.com/dgrijalva/jwt-go"
)

type JwtToken struct {
	Secret []byte
}

type JwtCsrfClaims struct {
	SessionID string `json:"sid"`
	UserID    uint32 `json:"uid"`
	jwt.StandardClaims
}

func NewJwtToken(secret string) (*JwtToken, error) {
	return &JwtToken{Secret: []byte(secret)}, nil
}

func (tk *JwtToken) Create(s *session.Session, tokenExpTime int64) (string, error) {
	data := JwtCsrfClaims{
		SessionID: s.ID,
		UserID:    s.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpTime,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	return token.SignedString(tk.Secret)
}

func (tk *JwtToken) parseSecretGetter(token *jwt.Token) (interface{}, error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok || method.Alg() != "HS256" {
		return nil, fmt.Errorf("bad sign method")
	}
	return tk.Secret, nil
}

func (tk *JwtToken) Check(s *session.Session, inputToken string) (bool, error) {
	payload := &JwtCsrfClaims{}
	_, err := jwt.ParseWithClaims(inputToken, payload, tk.parseSecretGetter)
	if err != nil {
		return false, fmt.Errorf("cant parse jwt token: %v", err)
	}
	// проверка exp, iat
	if payload.Valid() != nil {
		return false, errorTokenExpired
	}
	return payload.SessionID == s.ID && payload.UserID == s.UserID, nil
}
