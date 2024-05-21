package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"
)

// сюда писать код

func GetApp() http.Handler {
	return NewConduitAppHttpHandlers(NewRAMStorage())
}

type ConduitAppHttpHandlers struct {
	storage Storage
}

func NewConduitAppHttpHandlers(stor Storage) *ConduitAppHttpHandlers {
	return &ConduitAppHttpHandlers{
		storage: stor,
	}
}

func (srv *ConduitAppHttpHandlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	show("ConduitAppHttpHandlers.ServeHTTP, req: ", r.URL, r)
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

// updateCurrentUser: update current user profile, use storage to save data between requests
func (srv *ConduitAppHttpHandlers) updateCurrentUser(w http.ResponseWriter, r *http.Request) {
	var currentUser = UserProfile{
		Email:     "golang@example.com",
		CreatedAt: now_RFC3339(),
		UpdatedAt: now_RFC3339(),
		Username:  "golang",
	}
	var err error

	bodyMap, err := unmarshalBody(r)
	panicOnError("unmarshalBody failed", err)

	userData, userExists := bodyMap["user"]
	if userExists {
		userMap, err := giveMeStrings(userData.(map[string]any)) // recover from panic
		panicOnError("giveMeStrings failed", err)

		user, err := getUserByEmail(srv.storage, currentUser.Email)
		panicOnError("getUserByEmail failed", err)

		user.updateFromMap(userMap)

		err = putUserByEmail(srv.storage, user)
		panicOnError("putUserByEmail failed", err)

		err = writeUserResponse(w, user)
		panicOnError("writeUserResponse failed", err)

	} else {
		show("updateCurrentUser, no user in given data")
		http.Error(w, "oops", http.StatusBadRequest)
	}
}

func (srv *ConduitAppHttpHandlers) loginUser(w http.ResponseWriter, r *http.Request) {
	// show("loginUser, (email, password): ", r.FormValue("email"), r.FormValue("password"))
	var user = UserProfile{
		Email:     "golang@example.com",
		CreatedAt: now_RFC3339(),
		UpdatedAt: now_RFC3339(),
		Username:  "golang",
	}

	// err := putUserByEmail(srv.storage, user)
	// panicOnError("putUserByEmail failed", err)

	err := writeUserResponse(w, user)
	panicOnError("writeUserResponse failed", err)
}

func (srv *ConduitAppHttpHandlers) registerNewUser(w http.ResponseWriter, r *http.Request) {
	var user = UserProfile{
		Email:     "golang@example.com",
		CreatedAt: now_RFC3339(),
		UpdatedAt: now_RFC3339(),
		Username:  "golang",
	}

	err := putUserByEmail(srv.storage, user)
	panicOnError("putUserByEmail failed", err)

	var u = &struct{ User UserProfile }{User: user}
	err = writeResponseWithCode(w, u, http.StatusCreated)
	panicOnError("writeResponseWithCode failed", err)
}

func unmarshalBody(r *http.Request) (map[string]any, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err == nil {
		var bodyDecoded map[string]any
		err = json.Unmarshal(bodyBytes, &bodyDecoded)
		return bodyDecoded, err
	}
	return nil, err
}

func writeUserResponse(w http.ResponseWriter, user UserProfile) error {
	var x = &struct{ User UserProfile }{User: user}
	return writeResponse(w, x)
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
	// Set(key, value any) error
	// Get(key string) (any, error)
	GetAll() ([]any, error)
	SetAll([]any) error
}

type RAMStorage struct {
	data []any
}

var _ Storage = &RAMStorage{} // type check

func NewRAMStorage() *RAMStorage {
	return &RAMStorage{
		data: make([]any, 0, 16),
	}
}

// GetAll implements Storage.
func (stor RAMStorage) GetAll() ([]any, error) {
	return stor.data, nil
}

// SetAll implements Storage.
func (stor *RAMStorage) SetAll(items []any) error {
	stor.data = items
	return nil
}

// getUserByEmail: use email as key to search in stored items
func getUserByEmail(stor Storage, email string) (UserProfile, error) {
	items, err := stor.GetAll()
	if err == nil {
		for _, x := range items {
			user := x.(UserProfile) // panic if not only users in storage
			if user.Email == email {
				return user, nil
			}
		}
	} else {
		return UserProfile{}, err
	}
	return UserProfile{}, fmt.Errorf("getUserByEmail failed. email %#v not found in %d records", email, len(items))
}

// putUserByEmail: use email as key in stored items
func putUserByEmail(stor Storage, user UserProfile) error {
	items, err := stor.GetAll()
	if err == nil {
		items = slices.DeleteFunc(items, func(x any) bool {
			return (x.(UserProfile)).Email == user.Email // panic if not only users in storage
		})
		items = append(items, user)

		return stor.SetAll(items)
	}
	return err // fmt.Errorf("putUserByEmail failed. email: %#v; user: %#v", email, user)
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

func (user *UserProfile) updateFromMap(data map[string]string) *UserProfile {
	var x string
	var exists bool

	// TODO: add other user fields

	x, exists = data["email"]
	if exists {
		user.Email = x
	}

	x, exists = data["bio"]
	if exists {
		user.Bio = x
	}

	return user
}

func giveMeStrings(xs map[string]any) (map[string]string, error) {
	var ys = make(map[string]string, 16)
	for k, v := range xs {
		ys[k] = v.(string) // panic
	}
	return ys, nil
}

// --- system-wide tools ---

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
