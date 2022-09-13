package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

type SessionsDB struct {
	DB *sql.DB
}

func NewSessionsDB(db *sql.DB) *SessionsDB {
	return &SessionsDB{
		DB: db,
	}
}

// session id saved in cookie, load session from DB
func (sm *SessionsDB) Check(r *http.Request) (*Session, error) {
	sessionCookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		log.Println("CheckSession no cookie")
		return nil, ErrNoAuth
	}

	sess := &Session{} // create empty placeholder
	row := sm.DB.QueryRow(`SELECT user_id FROM sessions WHERE id = ?`, sessionCookie.Value)

	// load user id
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

// write new session to DB, set cookie
func (sm *SessionsDB) Create(w http.ResponseWriter, user *User) error {
	sessID := RandStringRunes(32)
	_, err := sm.DB.Exec("INSERT INTO sessions(id, user_id) VALUES(?, ?)", sessID, user.ID)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessID,
		Expires: time.Now().Add(90 * 24 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	return nil
}

// delete current session record, expire cookie
func (sm *SessionsDB) DestroyCurrent(w http.ResponseWriter, r *http.Request) error {
	sess, err := SessionFromContext(r.Context())
	if err == nil {
		_, err = sm.DB.Exec("DELETE FROM sessions WHERE id = ?", sess.ID)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:    "session_id",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}

	http.SetCookie(w, &cookie)
	return nil
}

// delete all user sessions records
func (sm *SessionsDB) DestroyAll(w http.ResponseWriter, user *User) error {
	result, err := sm.DB.Exec("DELETE FROM sessions WHERE user_id = ?",
		user.ID)
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	log.Println("destroyed sessions", affected, "for user", user.ID)

	return nil
}
