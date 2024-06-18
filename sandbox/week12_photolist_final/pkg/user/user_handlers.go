package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/asaskevich/govalidator"

	"photolist/pkg/session"
	"photolist/pkg/utils/httputils"
	"photolist/pkg/utils/randutils"
)

var GitHubEndpoint = oauth2.Endpoint{
	AuthURL:  "https://github.com/login/oauth/authorize",
	TokenURL: "https://github.com/login/oauth/access_token",
}

const (
	REDIRECT_URL = "http://localhost:8080/user/login_oauth" // callback
	AUTH_URL     = "https://github.com/login/oauth/authorize?scope=user:email&client_id=%s"
	API_URL      = "https://api.github.com/user?fields=email,photo_50&access_token=%s"
)

var (
	APP_ID     = "Ov2***gJF"
	APP_SECRET = "ada***860"
)

type Templater interface {
	Render(context.Context, http.ResponseWriter, string, map[string]interface{})
}

type UserHandler struct {
	Tmpl      Templater
	Sessions  session.SessionManager
	UsersRepo *UserRepository
}

var (
	loginRE = regexp.MustCompile(`^[\w-_\.]+$`)
)

func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		uh.Tmpl.Render(r.Context(), w, "login.html", map[string]interface{}{
			"OAuthURL": fmt.Sprintf(AUTH_URL, APP_ID),
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

	uh.Sessions.Create(r.Context(), w, user)
	http.Redirect(w, r, "/photos/", http.StatusFound)
}

func (uh *UserHandler) LoginOauth(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "no ouath code", http.StatusBadRequest)
		return
	}
	show("oauth code: ", code)

	conf := oauth2.Config{
		ClientID:     APP_ID,
		ClientSecret: APP_SECRET,
		RedirectURL:  REDIRECT_URL,
		Endpoint:     GitHubEndpoint,
	}
	ctx := r.Context()
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Println("exchange err", err)
		http.Error(w, "cannot get oauth token", http.StatusInternalServerError)
		return
	}
	show("oauth access token: ", token)

	// ask for data
	httpClient := conf.Client(ctx, token)
	apiResp, err := httpClient.Get(fmt.Sprintf(API_URL, token.AccessToken))
	if err != nil {
		log.Println("cannot request data from provider (api or token not working well)", err)
		http.Error(w, err.Error(), 500)
		return
	}
	defer apiResp.Body.Close()
	// decode api response
	respBodyBytes, err := io.ReadAll(apiResp.Body)
	if err != nil {
		log.Println("cannot read buffer", err)
		http.Error(w, err.Error(), 500)
		return
	}
	show("api response: ", string(respBodyBytes))
	userData := make(map[string]any, 32)
	err = json.Unmarshal(respBodyBytes, &userData)
	if err != nil {
		log.Println("cannot json.Unmarshal", err)
		http.Error(w, err.Error(), 500)
		return
	}
	if len(userData) == 0 {
		log.Println("requested data is empty", err)
		http.Error(w, "you should read the api docs", 500)
		return
	}
	// extract some data
	var email string
	var userID string
	emailAny, emailExists := userData["email"]
	if emailExists {
		email = emailAny.(string)
		show("user email from oauth provider: ", email)
	}
	uidAny, uidExists := userData["id"]
	if uidExists {
		userID = strconv.FormatUint(uint64(uidAny.(float64)), 36) // float64: json package to blame
		show("user id from oauth provider: ", userID)
	}

	// oauth profile adaptor (create vanilla user: login, password)
	login := "ghid_" + userID
	pass := randutils.RandStringRunes(50)
	show("creating app user. login, password: ", login, pass)

	user, err := uh.UsersRepo.Create(login, email, pass)
	if err != nil && err != errUserExists {
		log.Println("db err", err)
		http.Error(w, "Db err", http.StatusInternalServerError)
		return
	}

	uh.Sessions.Create(r.Context(), w, user)
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

	uh.Sessions.Create(r.Context(), w, user)
	http.Redirect(w, r, "/photos/", http.StatusFound)
}

func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	uh.Sessions.DestroyCurrent(r.Context(), w, r)
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

	sess, _ := session.SessionFromContext(r.Context())
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

	uh.Sessions.DestroyAll(r.Context(), w, user)
	uh.Sessions.Create(r.Context(), w, user)
	http.Redirect(w, r, "/photos/", http.StatusFound)
}

func (uh *UserHandler) FollowAPI(w http.ResponseWriter, r *http.Request) {
	sess, _ := session.SessionFromContext(r.Context())
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		httputils.RespJSONError(w, http.StatusBadRequest, nil, "bad id")
		return
	}
	folUser, err := uh.UsersRepo.GetByID(uint32(id))
	if err == errUserNotFound {
		httputils.RespJSONError(w, http.StatusBadRequest, nil, "no user")
		return
	}
	if err != nil {
		httputils.RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("db error: %v", err), "internal")
		return
	}

	rate := 1
	if r.FormValue("unfollow") == "1" {
		rate = -1
	}

	err = uh.UsersRepo.Follow(folUser.ID, sess.UserID, rate)
	if err != nil {
		httputils.RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("db error: %v", err), "internal")
		return
	}
	httputils.RespJSON(w, map[string]interface{}{
		"id": id,
	})
	return
}

type UserResp struct {
	ID    uint32 `json:"id"`
	Login string `json:"login"`
}

func (uh *UserHandler) FollowingAPI(w http.ResponseWriter, r *http.Request) {
	sess, _ := session.SessionFromContext(r.Context())
	users, err := uh.UsersRepo.GetFollowedUsers(sess.UserID)
	if err != nil {
		httputils.RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("db error: %v", err), "internal")
		return
	}
	result := make([]*UserResp, 0, len(users))
	for _, u := range users {
		result = append(result, &UserResp{
			ID:    u.ID,
			Login: u.Login,
		})
	}
	httputils.RespJSON(w, map[string]interface{}{
		"users":    result,
		"followed": true,
	})
}

func (uh *UserHandler) RecomendsAPI(w http.ResponseWriter, r *http.Request) {
	sess, _ := session.SessionFromContext(r.Context())
	users, err := uh.UsersRepo.GetRecomendedUsers(sess.UserID)
	if err != nil {
		httputils.RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("db error: %v", err), "internal")
		return
	}
	result := make([]*UserResp, 0, len(users))
	for _, u := range users {
		result = append(result, &UserResp{
			ID:    u.ID,
			Login: u.Login,
		})
	}
	httputils.RespJSON(w, map[string]interface{}{
		"users":    result,
		"followed": false,
	})
}

func (uh *UserHandler) InternalImagesAuth(w http.ResponseWriter, r *http.Request) {
	params := strings.Split(r.Header.Get("X-Original-URI"), "/")
	if len(params) != 4 {
		log.Println("bad params:", params)
		http.Error(w, "No auth", http.StatusForbidden)
		return
	}
	// log.Println("InternalImagesAuth params", params)

	sess, err := uh.Sessions.Check(r.Context(), r)
	if err != nil {
		log.Println("bad params:", err)
		http.Error(w, "Bad params", http.StatusForbidden)
		return
	}
	// log.Println("InternalImagesAuth sess", sess)

	targetUserID, err := strconv.Atoi(params[2])
	if err != nil {
		log.Println("bad uid:", err, params[2])
		http.Error(w, "No auth", http.StatusForbidden)
		return
	}

	if sess.UserID == uint32(targetUserID) {
		// 200 OK
		return
	}

	followed, err := uh.UsersRepo.IsFollowed(uint32(targetUserID), sess.UserID)
	if err != nil {
		log.Println("IsFollowed err:", err)
		http.Error(w, "Internal", http.StatusForbidden)
		return
	}

	if !followed {
		// no logs required - regular situation
		// log.Println("IsFollowed false", uint32(targetUserID), sess.UserID)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// 200 OK
}

// ts returns current timestamp in RFC3339 with milliseconds
func ts() string {
	/*
		https://pkg.go.dev/time#pkg-constants
		https://stackoverflow.com/questions/35479041/how-to-convert-iso-8601-time-in-golang
	*/
	const (
		RFC3339      = "2006-01-02T15:04:05Z07:00"
		RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
	)
	return time.Now().UTC().Format(RFC3339Milli)
}

// show writes message to standard output. Message combined from prefix msg and slice of arbitrary arguments
func show(msg string, xs ...any) {
	var line = ts() + ": " + msg

	for _, x := range xs {
		// https://pkg.go.dev/fmt
		// line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
}
