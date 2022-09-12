package main

import (
	"bytes"
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/argon2"
)

// user register, login processing

type UserHandler struct {
	DB   *sql.DB
	Tmpl *template.Template
}

// https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Password_Storage_Cheat_Sheet.md
func (u *UserHandler) hashPass(plainPassword, salt string) []byte {
	// [protected form] = [salt] + protect([protection func], [salt] + [credential]);
	hashedPass := argon2.IDKey([]byte(plainPassword), []byte(salt), 1, 64*1024, 4, 32)
	res := []byte(salt)
	return append(res, hashedPass...)
}

// login page processing
func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	// wrong method
	if r.Method != http.MethodPost {
		u.Tmpl.ExecuteTemplate(w, "login", nil)
		return
	}

	// go to DB

	login := r.FormValue("login")
	pass := r.FormValue("password")

	row := u.DB.QueryRow("SELECT id, password FROM users WHERE login = ?", login)

	var (
		dbPass []byte
		userID uint32
	)
	err := row.Scan(&userID, &dbPass)
	if err == sql.ErrNoRows {
		http.Error(w, "Bad user", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Db err", http.StatusInternalServerError)
		return
	}

	salt := string(dbPass[0:8])
	// check pass equality
	if !bytes.Equal(u.hashPass(pass, salt), dbPass) {
		http.Error(w, "Bad pass", http.StatusBadRequest)
		return
	}

	// user passed, create session
	CreateSession(w, r, u.DB, userID)

	// go to user home page
	http.Redirect(w, r, "/photos/", http.StatusFound)
}

// drop user session data from storage
func (u *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	DestroySession(w, r, u.DB)
	http.Redirect(w, r, "/user/login", http.StatusFound)
}

// register a new user
func (u *UserHandler) Reg(w http.ResponseWriter, r *http.Request) {
	// wrong method
	if r.Method != http.MethodPost {
		u.Tmpl.ExecuteTemplate(w, "reg", nil)
		return
	}

	login := r.FormValue("login")
	salt := RandStringRunes(8)
	pass := u.hashPass(r.FormValue("password"), salt)

	// ошибки игнорируются. никогда так не делайте :)
	// это будет исправлено в следующей итерации примера
	// сейчас так чтобы не отвлекаться от темы лекции
	result, err := u.DB.Exec("INSERT INTO users(login, password) VALUES(?, ?)", login, pass)
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

	// or, just go to login page
	userID, _ := result.LastInsertId()
	CreateSession(w, r, u.DB, uint32(userID))
	http.Redirect(w, r, "/photos/", http.StatusFound)
}
