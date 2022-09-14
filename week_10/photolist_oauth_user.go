package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/argon2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"
)

type User struct {
	ID    uint32
	Login string
	Ver   int32
}

// vk oauth
const (
	VK_APP_ID  = "7065390"
	VK_APP_KEY = "cQZe3Vvo4mHotmetUdXK"
	// куда идти для получения токена
	VK_AUTH_URL = "https://oauth.vk.com/authorize?client_id=7065390&redirect_uri=http://localhost:8080/user/login_oauth&response_type=code&scope=email"
	// куда идти с токеном за информацией
	VK_API_URL = "https://api.vk.com/method/users.get?fields=photo_50&access_token=%s&v=5.52"
)

// user data from vk
type VKOauthResp struct {
	Response []struct {
		FirstName string `json:"first_name"`
		Photo     string `json:"photo_50"`
	}
}

var (
	errNoRec      = errors.New("No user record found")
	errBadPass    = errors.New("Bad password")
	errUserExists = errors.New("User Exists")
)

// handlers for login, reg, ...
type UserHandler struct {
	DB       *sql.DB
	Tmpl     *template.Template
	Sessions SessionManager
}

// do pass hashing
func (uh *UserHandler) hashPass(plainPassword, salt string) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), []byte(salt), 1, 64*1024, 4, 32)
	res := []byte(salt)
	return append(res, hashedPass...)
}

// check pass vs db
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

// create session if login+password ok
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		uh.Tmpl.ExecuteTemplate(w, "login", map[string]string{
			"VKAuthURL": VK_AUTH_URL,
		})
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

// use vk to create session
func (uh *UserHandler) LoginOauth(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "no ouath code", http.StatusBadRequest)
		return
	}

	conf := oauth2.Config{
		ClientID:     VK_APP_ID,
		ClientSecret: VK_APP_KEY,
		RedirectURL:  "http://localhost:8080/user/login_oauth",
		Endpoint:     vk.Endpoint,
	}

	token, err := conf.Exchange(r.Context(), code)
	if err != nil {
		log.Println("exchange err", err)
		http.Error(w, "cannot get oauth token", http.StatusInternalServerError)
		return
	}

	emailRaw := token.Extra("email")
	email := ""
	okEmail := true
	if emailRaw != nil {
		email, okEmail = emailRaw.(string)
	}
	userIDraw, okID := token.Extra("user_id").(float64)

	// have not user data
	if !okEmail || !okID {
		log.Printf("cant convert data: UID: %T %v, Email: %T %v", token.Extra("user_id"), token.Extra("user_id"), token.Extra("email"), token.Extra("email"))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// email or vkid as login
	login := email
	if login == "" {
		login = "vk" + strconv.Itoa(int(userIDraw))
	}

	// random password to disable login with password
	pass := RandStringRunes(50)

	// write to db (or not, if exists)
	user, err := uh.createUser(login, pass)
	if err != nil && err != errUserExists {
		log.Println("db err", err)
		http.Error(w, "Db err", http.StatusInternalServerError)
		return
	}

	uh.Sessions.Create(w, user)
	http.Redirect(w, r, "/photos/", http.StatusFound)
}

// create or return existing user
func (uh *UserHandler) createUser(login, passIn string) (*User, error) {
	salt := RandStringRunes(8)
	pass := uh.hashPass(passIn, salt)

	user := &User{
		ID:    0,
		Ver:   0,
		Login: login,
	}

	// exist?
	err := uh.DB.QueryRow("SELECT id, ver FROM users WHERE login = ?", login).
		Scan(&user.ID, &user.Ver)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("db error: %v", err)
	}
	if err != sql.ErrNoRows {
		return user, errUserExists
	}

	// create
	result, err := uh.DB.Exec("INSERT INTO users(login, password) VALUES(?, ?)", login, pass)
	if err != nil {
		return nil, fmt.Errorf("insert error: %v", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, fmt.Errorf("no rows affected")
	}
	uid, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("LastInsertId err: %v", err)
	}
	user.ID = uint32(uid)

	return user, nil
}

// register new user, create session
func (uh *UserHandler) Reg(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		uh.Tmpl.ExecuteTemplate(w, "reg", nil)
		return
	}

	login := r.FormValue("login")
	pass := r.FormValue("password")

	user, err := uh.createUser(login, pass)
	switch err {
	case nil:
		// all is ok
	case errUserExists:
		http.Error(w, "Looks like user exists", http.StatusBadRequest)
	default:
		log.Println("db err", err)
		http.Error(w, "Db err", http.StatusInternalServerError)
	}
	if err != nil {
		return
	}

	uh.Sessions.Create(w, user)
	http.Redirect(w, r, "/photos/", http.StatusFound)
}

func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	uh.Sessions.DestroyCurrent(w, r)
	http.Redirect(w, r, "/user/login", http.StatusFound)
}

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
