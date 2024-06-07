package session

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"photolist/pkg/utils/randutils"
)

type SessionsDB struct {
	DB *sql.DB
}

func NewSessionsDB(db *sql.DB) *SessionsDB {
	return &SessionsDB{
		DB: db,
	}
}

func (sm *SessionsDB) Check(r *http.Request) (*Session, error) {
	sessionCookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		log.Println("CheckSession no cookie")
		return nil, ErrNoAuth
	}

	sess := &Session{}
	row := sm.DB.QueryRow(`SELECT user_id FROM sessions WHERE id = ?`, sessionCookie.Value)
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

func (sm *SessionsDB) Create(w http.ResponseWriter, user UserInterface) error {
	sessID := randutils.RandStringRunes(32)
	_, err := sm.DB.Exec("INSERT INTO sessions(id, user_id) VALUES(?, ?)", sessID, user.GetID())
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

func (sm *SessionsDB) DestroyAll(w http.ResponseWriter, user UserInterface) error {
	result, err := sm.DB.Exec("DELETE FROM sessions WHERE user_id = ?",
		user.GetID())
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	log.Println("destroyed sessions", affected, "for user", user.GetID())

	return nil
}
