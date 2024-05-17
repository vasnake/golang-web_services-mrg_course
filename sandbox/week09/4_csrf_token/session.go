package csrf_token

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"
)

// same shit as before

var (
	noAuthUrls = map[string]struct{}{
		"/user/login": struct{}{},
		"/user/reg":   struct{}{},
		"/":           struct{}{},
	}
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

func CreateSession(w http.ResponseWriter, r *http.Request, db *sql.DB, userID uint32) error {
	sessID := RandStringRunes(32)
	db.Exec("INSERT INTO sessions(id, user_id) VALUES(?, ?)", sessID, userID)

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessID,
		Expires: time.Now().Add(90 * 24 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	return nil
}

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

// AuthMiddleware: on every request go to db, read session data, set context
func AuthMiddleware(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := noAuthUrls[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}
		sess, err := CheckSession(db, r)
		if err != nil {
			http.Error(w, "No auth", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), sessionKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CheckSession: get id from cookie, read session data from db
func CheckSession(db *sql.DB, r *http.Request) (*Session, error) {
	sessionCookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		log.Println("CheckSession no cookie")
		return nil, ErrNoAuth
	}

	sess := &Session{}
	row := db.QueryRow(`SELECT user_id FROM sessions WHERE id = ?`, sessionCookie.Value)
	err = row.Scan(&sess.UserID)
	if err == sql.ErrNoRows {
		log.Println("CheckSession no rows")
		return nil, ErrNoAuth
	} else if err != nil {
		log.Println("CheckSession err:", err)
		return nil, err
	}

	sess.ID = sessionCookie.Value
	return sess, nil
}

// simple getter
func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(sessionKey).(*Session)
	if !ok {
		return nil, ErrNoAuth
	}
	return sess, nil
}
