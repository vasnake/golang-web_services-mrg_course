package main

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
	ID     string // session id
}

// линтер ругается если используем базовые типы в Value контекста
// типа так безопаснее разграничивать
type ctxKey int

const sessionKey ctxKey = 1

var (
	ErrNoAuth = errors.New("No session found")
)

// get session or error from context
func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(sessionKey).(*Session)
	if !ok {
		return nil, ErrNoAuth
	}
	return sess, nil
}

// use cookie value to get session from DB, if no cookie or no DB record => error
func CheckSession(db *sql.DB, r *http.Request) (*Session, error) {
	sessionCookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		log.Println("CheckSession no cookie")
		return nil, ErrNoAuth
	}

	sess := &Session{}

	// session id as cookie value
	sess.ID = sessionCookie.Value

	row := db.QueryRow(`SELECT user_id FROM sessions WHERE id = ?`, sessionCookie.Value)

	// user data from DB
	err = row.Scan(&sess.UserID)

	if err == sql.ErrNoRows {
		log.Println("CheckSession no rows")
		return nil, ErrNoAuth
	} else if err != nil {
		log.Println("CheckSession err:", err)
		return nil, err
	}

	return sess, nil
}

// create new session for user
func CreateSession(w http.ResponseWriter, r *http.Request, db *sql.DB, userID uint32) error {
	// TODO: add errors processing

	// generate session id
	sessID := RandStringRunes(32)
	// write to DB
	db.Exec("INSERT INTO sessions(id, user_id) VALUES(?, ?)", sessID, userID)

	// set cooklie
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessID,
		Expires: time.Now().Add(90 * 24 * time.Hour),
		Path:    "/",
	}

	http.SetCookie(w, cookie)

	return nil
}

// drop session from DB and from cooklie
func DestroySession(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	sess, err := SessionFromContext(r.Context())
	if err == nil {
		db.Exec("DELETE FROM sessions WHERE id = ?", sess.ID)
	}

	cookie := http.Cookie{
		Name:    "session_id",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)

	return nil
}

// auth middleware
var (
	// empty struct requires no memory to store, map become a set of keys
	noAuthUrls = map[string]struct{}{
		"/user/login": struct{}{},
		"/user/reg":   struct{}{},
		"/":           struct{}{},
	}
)

func AuthMiddleware(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// no need no auth here
		if _, ok := noAuthUrls[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}

		// require auth to access page
		sess, err := CheckSession(db, r)
		if err != nil {
			http.Error(w, "No auth", http.StatusUnauthorized)
			return
		}

		// add session to context and go downstream
		ctx := context.WithValue(r.Context(), sessionKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
