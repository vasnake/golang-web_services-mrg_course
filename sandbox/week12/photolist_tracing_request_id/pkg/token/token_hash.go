package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"photolist/pkg/session"
)

type HashToken struct {
	Secret []byte
}

func NewHMACHashToken(secret string) (*HashToken, error) {
	return &HashToken{Secret: []byte(secret)}, nil
}

func (tk *HashToken) Create(s *session.Session, tokenExpTime int64) (string, error) {
	h := hmac.New(sha256.New, []byte(tk.Secret))
	data := fmt.Sprintf("%s:%d:%d", s.ID, s.UserID, tokenExpTime)
	h.Write([]byte(data))
	token := hex.EncodeToString(h.Sum(nil)) + ":" + strconv.FormatInt(tokenExpTime, 10)
	return token, nil
}

func (tk *HashToken) Check(s *session.Session, inputToken string) (bool, error) {
	tokenData := strings.Split(inputToken, ":")
	if len(tokenData) != 2 {
		return false, fmt.Errorf("bad token data")
	}

	tokenExp, err := strconv.ParseInt(tokenData[1], 10, 64)
	if err != nil {
		return false, fmt.Errorf("bad token time")
	}

	if tokenExp < time.Now().Unix() {
		return false, errorTokenExpired
	}

	h := hmac.New(sha256.New, []byte(tk.Secret))
	data := fmt.Sprintf("%s:%d:%d", s.ID, s.UserID, tokenExp)
	h.Write([]byte(data))
	expectedMAC := h.Sum(nil)
	messageMAC, err := hex.DecodeString(tokenData[0])
	if err != nil {
		return false, fmt.Errorf("cand hex decode token")
	}

	return hmac.Equal(messageMAC, expectedMAC), nil
}
