package main

import (
	"bytes"
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/argon2"
)

type User struct {
	ID    uint32
	Login string
	Ver   int32
}

type UserHandler struct {
	DB       *sql.DB
	Tmpl     *template.Template
	Sessions SessionManager // interface
}

func (uh *UserHandler) hashPass(plainPassword, salt string) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), []byte(salt), 1, 64*1024, 4, 32)
	res := []byte(salt)
	return append(res, hashedPass...)
}

var (
	errNoRec   = errors.New("No user record found")
	errBadPass = errors.New("Wrong password")
)

// get user from DB, compare hashed passwords
func (uh *UserHandler) passwordIsValid(pass string, row *sql.Row) (*User, error) {
	var (
		dbPass []byte
		user   = &User{}
	)

	err := row.Scan(&user.ID, &user.Login, &user.Ver, &dbPass)
	if err == sql.ErrNoRows {
		return nil, errNoRec
	} else if err != nil {
		return nil, err
	}

	// compare hashed passwords
	salt := string(dbPass[0:8])
	if !bytes.Equal(uh.hashPass(pass, salt), dbPass) {
		return nil, errBadPass
	}

	return user, nil
}

func (uh *UserHandler) checkPasswordByUserID(uid uint32, pass string) (*User, error) {
	row := uh.DB.QueryRow("SELECT id, login, ver, password FROM users WHERE id = ?", uid)
	return uh.passwordIsValid(pass, row)
}

func (uh *UserHandler) checkPasswordByLogin(login, pass string) (*User, error) {
	row := uh.DB.QueryRow("SELECT id, login, ver, password FROM users WHERE login = ?", login)
	return uh.passwordIsValid(pass, row)
}

// check password, create session, redirect to homepage
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		uh.Tmpl.ExecuteTemplate(w, "login", nil)
		return
	}

	login := r.FormValue("login")
	pass := r.FormValue("password")

	user, err := uh.checkPasswordByLogin(login, pass)

	switch err {
	case nil:
		// all is ok
	case errNoRec:
		http.Error(w, "No user", http.StatusBadRequest)
	case errBadPass:
		http.Error(w, "Bad pass", http.StatusBadRequest)
	default:
		http.Error(w, "Db err", http.StatusInternalServerError)
	}
	if err != nil {
		return
	}

	uh.Sessions.Create(w, user)
	http.Redirect(w, r, "/photos/", http.StatusFound)
}

// create new user record in DB, create session, redirect to homepage
func (uh *UserHandler) Reg(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		uh.Tmpl.ExecuteTemplate(w, "reg", nil)
		return
	}
	login := r.FormValue("login")
	salt := RandStringRunes(8)
	pass := uh.hashPass(r.FormValue("password"), salt)

	result, err := uh.DB.Exec("INSERT INTO users(login, password) VALUES(?, ?)", login, pass)
	if err != nil {
		log.Println("insert error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "Looks like user exists", http.StatusBadRequest)
		return
	}
	userID, _ := result.LastInsertId()

	user := &User{
		ID:    uint32(userID),
		Ver:   0,
		Login: login,
	}
	uh.Sessions.Create(w, user)
	http.Redirect(w, r, "/photos/", http.StatusFound)
}

// drop current session
func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	uh.Sessions.DestroyCurrent(w, r)
	http.Redirect(w, r, "/user/login", http.StatusFound)
}

// check old password, update user record in DB, destroy all user sessions, create new session, redirect to homepage
func (uh *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		uh.Tmpl.ExecuteTemplate(w, "change_password", nil)
		return
	}

	if r.FormValue("pass1") == "" || r.FormValue("pass1") != r.FormValue("pass2") {
		http.Error(w, "New password mistmatch", http.StatusBadRequest)
		return
	}

	sess, _ := SessionFromContext(r.Context())
	user, err := uh.checkPasswordByUserID(sess.UserID, r.FormValue("old_password"))
	if err != nil {
		http.Error(w, "Bad pass", http.StatusBadRequest)
		return
	}

	salt := RandStringRunes(8)
	pass := uh.hashPass(r.FormValue("pass1"), salt)

	_, err = uh.DB.Exec("UPDATE users SET password = ?, ver = ver + 1 WHERE id = ?",
		pass, user.ID)
	if err != nil {
		log.Println("update password error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	user.Ver++ // во избежание рейсов лучше подгрузить из базы

	uh.Sessions.DestroyAll(w, user)
	uh.Sessions.Create(w, user)
	http.Redirect(w, r, "/photos/", http.StatusFound)
}
