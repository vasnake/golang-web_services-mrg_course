package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	// "strings"
	"time"
)

// token service, way to pass a secret
type CryptToken struct {
	Secret []byte
}

// token structure
type TokenData struct {
	SessionID string
	UserID    uint32
	Exp       int64
}

// token service factory
func NewAesCryptHashToken(secret string) (*CryptToken, error) {
	key := []byte(secret)
	_, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("cypher problem %v", err)
	}
	return &CryptToken{Secret: key}, nil
}

// create token from session data
func (tk *CryptToken) Create(s *Session, tokenExpTime int64) (string, error) {
	block, err := aes.NewCipher(tk.Secret)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesgcm.NonceSize()) // header
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// data
	td := &TokenData{SessionID: s.ID, UserID: s.UserID, Exp: tokenExpTime}
	data, _ := json.Marshal(td)

	ciphertext := aesgcm.Seal(nil, nonce, data, nil) // crypt

	res := append([]byte(nil), nonce...)
	res = append(res, ciphertext...)

	token := base64.StdEncoding.EncodeToString(res)
	return token, nil
}

// decrypt, check data equality
func (tk *CryptToken) Check(s *Session, inputToken string) (bool, error) {
	block, err := aes.NewCipher(tk.Secret)
	if err != nil {
		return false, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return false, err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(inputToken)
	if err != nil {
		return false, err
	}

	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return false, fmt.Errorf("ciphertext too short")
	}

	// header, crypted payload
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// decrypt
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return false, fmt.Errorf("decrypt fail: %v", err)
	}

	// decode token
	td := TokenData{}
	err = json.Unmarshal(plaintext, &td)
	if err != nil {
		return false, fmt.Errorf("bad json: %v", err)
	}

	if td.Exp < time.Now().Unix() {
		return false, fmt.Errorf("token expired")
	}

	// check equality
	expected := TokenData{SessionID: s.ID, UserID: s.UserID}
	td.Exp = 0
	return td == expected, nil
}
