package auth

import (
	"bytes"
	"database/sql"
	"golang.org/x/crypto/argon2"
	"html/template"
	"log"
	"net/http"
)

// http.Handler for reg, login, logout
type UserHandler struct {
	DB   *sql.DB
	Tmpl *template.Template
}

func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	// post: do login, not-post: show login form
	if r.Method != http.MethodPost {
		uh.Tmpl.ExecuteTemplate(w, "login", nil)
		return
	}

	login := r.FormValue("login")
	givenPass := r.FormValue("password")
	var (
		dbPass []byte
		userID uint32
	)

	row := uh.DB.QueryRow("SELECT id, password FROM users WHERE login = ?", login)
	err := row.Scan(&userID, &dbPass)
	if err == sql.ErrNoRows {
		http.Error(w, "no such user in db", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "db query failed", http.StatusInternalServerError)
		return
	}

	salt := string(dbPass[0:8]) // magic: 8 bytes
	if !bytes.Equal(hashPass(givenPass, salt), dbPass) {
		http.Error(w, "wrong password, guess again", http.StatusBadRequest)
		return
	}

	err = CreateSession(w, r, uh.DB, userID)
	if err != nil {
		panic("not implemented error handler")
	}

	http.Redirect(w, r, "/photos/", http.StatusFound) // should be redirected to `/`?
}

func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	DestroySession(w, r, uh.DB)
	http.Redirect(w, r, "/user/login", http.StatusFound) // should be redirected to `/`?
}

func (uh *UserHandler) Reg(w http.ResponseWriter, r *http.Request) {
	// post: do register, not-post: show form
	if r.Method != http.MethodPost {
		uh.Tmpl.ExecuteTemplate(w, "reg", nil)
		return
	}

	login := r.FormValue("login")
	salt := RandStringRunes(8) // bytes vs runes, magic; 8: magic
	pass := hashPass(r.FormValue("password"), salt)

	// ошибки игнорируются. никогда так не делайте :)
	// это будет исправлено в следующей итерации примера
	// сейчас так чтобы не отвлекаться от темы лекции
	result, err := uh.DB.Exec("INSERT INTO users(login, password) VALUES(?, ?)", login, pass)
	if err != nil {
		log.Println("insert error", err)
		http.Error(w, "db insert failed", http.StatusInternalServerError)
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "Looks like user exists", http.StatusBadRequest)
		return
	}
	userID, _ := result.LastInsertId()

	CreateSession(w, r, uh.DB, uint32(userID))
	http.Redirect(w, r, "/photos/", http.StatusFound) // or `/`?
}

// https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/Password_Storage_Cheat_Sheet.md
// // [protected form] = [salt] + protect([protection func], [salt] + [credential]);
func hashPass(plainPassword, salt string) []byte {
	saltBytes := []byte(salt)
	hashedPass := argon2.IDKey([]byte(plainPassword), saltBytes, 1, 64*1024, 4, 32)
	saltAndHashedPass := make([]byte, len(salt)+len(hashedPass))
	copy(saltAndHashedPass, saltBytes)
	return append(saltAndHashedPass, hashedPass...)
}
