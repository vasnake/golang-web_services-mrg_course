package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// сюда писать код

func GetApp() http.Handler {
	return NewConduitAppHttpHandlers()
}

type ConduitAppHttpHandlers struct {
	storage Storage
}

func NewConduitAppHttpHandlers() *ConduitAppHttpHandlers {
	return &ConduitAppHttpHandlers{
		storage: nil,
	}
}

func (srv *ConduitAppHttpHandlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	show("ConduitAppHttpHandlers.ServeHTTP, req: ", r.URL, r)
	// http.Error(w, "foo", http.StatusCreated)
	switch r.Method {

	case "POST":
		srv.servePost(w, r)
	case "GET":
		srv.serveGet(w, r)
	case "PUT":
		srv.servePut(w, r)

	default:
		show("unknown method: ", r.Method)
	}
}

func (srv *ConduitAppHttpHandlers) servePut(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {

	case "/api/user":
		srv.updateCurrentUser(w, r)

	default:
		show("unknown put endpoint: ", r.URL.Path)
	}
}

func (srv *ConduitAppHttpHandlers) servePost(w http.ResponseWriter, r *http.Request) {
	// show("ConduitAppHttpHandlers.servePost, req: ", r.URL.Path, r)
	switch r.URL.Path {

	case "/api/users":
		srv.registerNewUser(w, r)
	case "/api/users/login":
		srv.loginUser(w, r)

	default:
		show("unknown post endpoint: ", r.URL.Path)
	}
}

func (srv *ConduitAppHttpHandlers) serveGet(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {

	case "/api/user":
		srv.showCurrentUser(w, r)

	default:
		show("unknown get endpoint: ", r.URL.Path)
	}
}

func (srv *ConduitAppHttpHandlers) showCurrentUser(w http.ResponseWriter, r *http.Request) {
	var x = &struct {
		User UserProfile
	}{
		User: UserProfile{
			Email:     "golang@example.com",
			CreatedAt: now_RFC3339(),
			UpdatedAt: now_RFC3339(),
			Username:  "golang",
		},
	}
	err := writeResponse(w, x)
	panicOnError("writeResponse failed", err)
}

func (srv *ConduitAppHttpHandlers) updateCurrentUser(w http.ResponseWriter, r *http.Request) {
	// TODO: update current user profile, use storage to save data between requests
	/*
		Email:     "golang@example.com",
		Username:  "golang",
		Bio:"",
		=>
		Email:"u_golang@example.com",
		Username:"golang",
		Bio:"Info about golang",
	*/

	var x = &struct {
		User UserProfile
	}{
		User: UserProfile{
			// Email:     "u_golang@example.com",
			CreatedAt: now_RFC3339(),
			UpdatedAt: now_RFC3339(),
			Username:  "golang",
			// Bio:       "Info about golang",
		},
	}

	err := writeResponse(w, x)
	panicOnError("writeResponse failed", err)
}

func (srv *ConduitAppHttpHandlers) registerNewUser(w http.ResponseWriter, r *http.Request) {
	var x = &struct {
		User UserProfile
	}{
		User: UserProfile{
			Email:     "golang@example.com",
			CreatedAt: now_RFC3339(),
			UpdatedAt: now_RFC3339(),
			Username:  "golang",
		},
	}

	err := writeResponseWithCode(w, x, http.StatusCreated)
	panicOnError("writeResponse failed", err)
}

func (srv *ConduitAppHttpHandlers) loginUser(w http.ResponseWriter, r *http.Request) {
	show("loginUser, (email, password): ", r.FormValue("email"), r.FormValue("password"))
	var x = &struct {
		User UserProfile
	}{
		User: UserProfile{
			Email:     "golang@example.com",
			CreatedAt: now_RFC3339(),
			UpdatedAt: now_RFC3339(),
			Username:  "golang",
		},
	}

	err := writeResponse(w, x)
	panicOnError("writeResponse failed", err)
}

func writeResponse(w http.ResponseWriter, x any) error {
	return writeResponseWithCode(w, x, http.StatusOK)
}

func writeResponseWithCode(w http.ResponseWriter, x any, code int) error {
	resp, err := json.Marshal(x)
	if err == nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(code)
		_, err = w.Write(resp)
		return err
	}
	return err
}

type Storage interface {
	Set(key, value string) error
	Get(key string) (value string, err error)
}

type UserProfile struct {
	ID        string `json:"id" testdiff:"ignore"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"` // RFC3339     = "2006-01-02T15:04:05Z07:00"
	UpdatedAt string `json:"updatedAt"` // RFC3339     = "2006-01-02T15:04:05Z07:00"
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Token     string `json:"token" testdiff:"ignore"`
	Following bool
}

func panicOnError(msg string, err error) {
	if err != nil {
		panic(msg + ": " + err.Error())
	}
}

// ts returns current timestamp in RFC3339 with milliseconds
func ts() string {
	/*
		https://pkg.go.dev/time#pkg-constants
		https://stackoverflow.com/questions/35479041/how-to-convert-iso-8601-time-in-golang
	*/
	const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
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

func now_RFC3339() string {
	const RFC3339 = "2006-01-02T15:04:05Z07:00"
	return time.Now().UTC().Format(RFC3339)
}
