package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"

	"github.com/asaskevich/govalidator"
)

// oauth2

const (
	VK_APP_ID  = "7065390"
	VK_APP_KEY = "cQZe3Vvo4mHotmetUdXK"
	// куда идти с токеном за информацией
	VK_API_URL = "https://api.vk.com/method/users.get?fields=photo_50&access_token=%s&v=5.52"
	// куда идти для получения токена
	VK_AUTH_URL = "https://oauth.vk.com/authorize?client_id=7065390&redirect_uri=http://localhost:8080/user/login_oauth&response_type=code&scope=email"
)

type VKOauthResp struct {
	Response []struct {
		FirstName string `json:"first_name"`
		Photo     string `json:"photo_50"`
	}
}

// user methods

type UserHandler struct {
	Tmpl      *MyTemplate
	Sessions  SessionManager
	UsersRepo *UserRepository
}

var (
	loginRE = regexp.MustCompile(`^[\w-_\.]+$`)
)

func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		uh.Tmpl.Render(r.Context(), w, "login.html", map[string]interface{}{
			"VKAuthURL": VK_AUTH_URL,
		})
		return
	}

	login := r.FormValue("login")
	pass := r.FormValue("password")

	user, err := uh.UsersRepo.CheckPasswordByLogin(login, pass)
	switch err {
	case nil:
		// all is ok
	case errUserNotFound:
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
	if !okEmail || !okID {
		log.Printf("cant convert data: UID: %T %v, Email: %T %v", token.Extra("user_id"), token.Extra("user_id"), token.Extra("email"), token.Extra("email"))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	login := "vk" + strconv.Itoa(int(userIDraw))
	pass := RandStringRunes(50)
	user, err := uh.UsersRepo.Create(login, email, pass)
	if err != nil && err != errUserExists {
		log.Println("db err", err)
		http.Error(w, "Db err", http.StatusInternalServerError)
		return
	}

	uh.Sessions.Create(w, user)
	http.Redirect(w, r, "/photos/", http.StatusFound)
}

func (uh *UserHandler) Reg(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		uh.Tmpl.Render(r.Context(), w, "reg.html", nil)
		return
	}

	login := r.FormValue("login")
	pass := r.FormValue("password")
	email := r.FormValue("email")

	if !govalidator.IsEmail(email) {
		http.Error(w, "Bad email", http.StatusBadRequest)
		return
	}

	if !loginRE.MatchString(login) {
		http.Error(w, "Bad login", http.StatusBadRequest)
		return
	}

	user, err := uh.UsersRepo.Create(login, email, pass)
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
		uh.Tmpl.Render(r.Context(), w, "change_pass.html", nil)
		return
	}

	if r.FormValue("pass1") == "" || r.FormValue("pass1") != r.FormValue("pass2") {
		http.Error(w, "New password mistmatch", http.StatusBadRequest)
		return
	}

	sess, _ := SessionFromContext(r.Context())
	user, err := uh.UsersRepo.CheckPasswordByUserID(sess.UserID, r.FormValue("old_password"))
	if err != nil {
		http.Error(w, "Bad pass", http.StatusBadRequest)
		return
	}

	err = uh.UsersRepo.UpdatePassword(user.ID, r.FormValue("pass1"))
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

func (uh *UserHandler) FollowAPI(w http.ResponseWriter, r *http.Request) {
	sess, _ := SessionFromContext(r.Context())
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		RespJSONError(w, http.StatusBadRequest, nil, "bad id")
		return
	}
	folUser, err := uh.UsersRepo.GetByID(uint32(id))
	if err == errUserNotFound {
		RespJSONError(w, http.StatusBadRequest, nil, "no user")
		return
	}
	if err != nil {
		RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("db error: %v", err), "internal")
		return
	}

	rate := 1
	if r.FormValue("unfollow") == "1" {
		rate = -1
	}

	err = uh.UsersRepo.Follow(folUser.ID, sess.UserID, rate)
	if err != nil {
		RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("db error: %v", err), "internal")
		return
	}
	RespJSON(w, map[string]interface{}{
		"id": id,
	})
	return
}

type UserResp struct {
	ID    uint32 `json:"id"`
	Login string `json:"login"`
}

func (uh *UserHandler) FollowingAPI(w http.ResponseWriter, r *http.Request) {
	sess, _ := SessionFromContext(r.Context())
	users, err := uh.UsersRepo.GetFollowedUsers(sess.UserID)
	if err != nil {
		RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("db error: %v", err), "internal")
		return
	}
	result := make([]*UserResp, 0, len(users))
	for _, u := range users {
		result = append(result, &UserResp{
			ID:    u.ID,
			Login: u.Login,
		})
	}
	RespJSON(w, map[string]interface{}{
		"users":    result,
		"followed": true,
	})
}

func (uh *UserHandler) RecomendsAPI(w http.ResponseWriter, r *http.Request) {
	sess, _ := SessionFromContext(r.Context())
	users, err := uh.UsersRepo.GetRecomendedUsers(sess.UserID)
	if err != nil {
		RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("db error: %v", err), "internal")
		return
	}
	result := make([]*UserResp, 0, len(users))
	for _, u := range users {
		result = append(result, &UserResp{
			ID:    u.ID,
			Login: u.Login,
		})
	}
	RespJSON(w, map[string]interface{}{
		"users":    result,
		"followed": false,
	})
}
