package csrf_token

import (
	"bytes"
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/argon2"
)

// same shit as before

// http handler: reg, login, logout
type UserHandler struct {
	DB   *sql.DB
	Tmpl *template.Template
}

func (u *UserHandler) Reg(w http.ResponseWriter, r *http.Request) {
	// show form or do register a new user
	if r.Method != http.MethodPost {
		u.Tmpl.ExecuteTemplate(w, "reg", nil)
		return
	}

	login := r.FormValue("login")
	pass := hashPass(r.FormValue("password"), RandStringRunes(8))

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

	userID, _ := result.LastInsertId()

	CreateSession(w, r, u.DB, uint32(userID))
	http.Redirect(w, r, "/photos/", http.StatusFound)
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		u.Tmpl.ExecuteTemplate(w, "login", nil)
		return
	}

	login := r.FormValue("login")
	pass := r.FormValue("password")

	row := u.DB.QueryRow("SELECT id, password FROM users WHERE login = ?", login)
	var (
		dbPass []byte
		userID uint32
	)
	err := row.Scan(&userID, &dbPass)
	if err == sql.ErrNoRows {
		http.Error(w, "No user", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Db err", http.StatusInternalServerError)
		return
	}

	salt := string(dbPass[0:8])
	if !bytes.Equal(hashPass(pass, salt), dbPass) {
		http.Error(w, "Bad pass", http.StatusBadRequest)
		return
	}

	CreateSession(w, r, u.DB, userID)
	http.Redirect(w, r, "/photos/", http.StatusFound)
}

func (u *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	DestroySession(w, r, u.DB)
	http.Redirect(w, r, "/user/login", http.StatusFound)
}

func hashPass(plainPassword, salt string) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), []byte(salt), 1, 64*1024, 4, 32)
	res := make([]byte, len(salt))
	copy(res, salt[:len(salt)])
	return append(res, hashedPass...)
}
