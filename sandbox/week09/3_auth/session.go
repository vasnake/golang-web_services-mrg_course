package auth

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"
)

type Session struct {
	UserID uint32
	ID     string
}

// линтер ругается если используем базовые типы в Value контекста
// типа так безопаснее разграничивать
type ctxKey int

const sessionKey ctxKey = 1

var (
	ErrNoAuth = errors.New("No session found")
)

var unit = struct{}{}
var (
	noAuthUrls = map[string]struct{}{
		"/user/login": unit,
		"/user/reg":   unit,
		"/":           unit,
	}
)

// CreateSession: write session to db, create cookie
func CreateSession(w http.ResponseWriter, r *http.Request, db *sql.DB, userID uint32) error {
	sessID := RandStringRunes(32) // not so unique as you think
	db.Exec("INSERT INTO sessions(id, user_id) VALUES(?, ?)", sessID, userID)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessID,
		Expires: time.Now().Add(90 * 24 * time.Hour), // magic
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	return nil
}

// DestroySession: drop row from db and drop cookie
func DestroySession(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	sess, err := SessionFromContext(r.Context())
	if err == nil {
		db.Exec("DELETE FROM sessions WHERE id = ?", sess.ID)
	}
	cookie := http.Cookie{
		Name:    "session_id", // magic
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	return nil
}

// AuthMiddleware: add session check to req. handler
func AuthMiddleware(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := noAuthUrls[r.URL.Path]; ok {
			next.ServeHTTP(w, r) // public path
			return
		}

		sess, err := CheckSession(db, r) // read session data from db
		if err != nil {
			http.Error(w, "you need to auth first", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), sessionKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx)) // session object available downstream
	})
}

// CheckSession: read data from db using id from cookie
func CheckSession(db *sql.DB, r *http.Request) (*Session, error) {
	sessionCookie, err := r.Cookie("session_id") // magic
	if err == http.ErrNoCookie {
		log.Println("CheckSession: no cookie")
		return nil, ErrNoAuth
	}

	// read from db
	sess := &Session{ID: sessionCookie.Value}
	row := db.QueryRow(`SELECT user_id FROM sessions WHERE id = ?`, sess.ID)
	err = row.Scan(&sess.UserID)
	if err == sql.ErrNoRows {
		log.Println("CheckSession: no rows")
		return nil, ErrNoAuth
	} else if err != nil {
		log.Println("CheckSession: db err:", err)
		return nil, err
	}

	return sess, nil
}

// SessionFromContext: extract session object from context
func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(sessionKey).(*Session)
	if !ok {
		return nil, ErrNoAuth
	}
	return sess, nil
}
