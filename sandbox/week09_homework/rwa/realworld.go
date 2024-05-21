package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"
)

// сюда писать код

func GetApp() http.Handler {
	return NewConduitAppHttpHandlers(
		// DI
		NewRAMStorage(),
		NewSimpleSessionManager(),
	)
}

var _ http.Handler = &ConduitAppHttpHandlers{} // type check

type ConduitAppHttpHandlers struct {
	storage  Storage
	sessions SessionManager
}

func NewConduitAppHttpHandlers(stor Storage, sm SessionManager) *ConduitAppHttpHandlers {
	return &ConduitAppHttpHandlers{
		storage:  stor,
		sessions: sm,
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

// updateCurrentUser: update current user profile, use storage to save data between requests
func (srv *ConduitAppHttpHandlers) updateCurrentUser(w http.ResponseWriter, r *http.Request) {
	// read: email, bio, token
	sessionID, err := getSessionIDFromReq(r)
	panicOnError("getSessionIDFromReq failed", err)
	userID, err := srv.sessions.GetUserID(sessionID)
	panicOnError("GetUserID failed", err)
	user, err := getUserFromStorage(srv.storage, userID)
	panicOnError("getUserFromStorage failed", err)

	user.Token = sessionID

	bodyMap, err := unmarshalBody(r)
	panicOnError("unmarshalBody failed", err)

	userData, userExists := bodyMap["user"]
	if userExists {
		userMap, err := giveMeStrings(userData.(map[string]any)) // recover from panic
		panicOnError("giveMeStrings failed", err)

		user.updateFromMap(userMap)

		err = putUser2Storage(srv.storage, user)
		panicOnError("putUserByEmail failed", err)

		err = writeUserResponse(w, user)
		panicOnError("writeUserResponse failed", err)

	} else {
		show("updateCurrentUser, no user in given data")
		http.Error(w, "oops", http.StatusBadRequest)
	}
}

func (srv *ConduitAppHttpHandlers) showCurrentUser(w http.ResponseWriter, r *http.Request) {
	// get userID from token, load user data
	sessionID, err := getSessionIDFromReq(r)
	panicOnError("getSessionIDFromReq failed", err)
	userID, err := srv.sessions.GetUserID(sessionID)
	panicOnError("GetUserID failed", err)
	user, err := getUserFromStorage(srv.storage, userID)
	panicOnError("getUserFromStorage failed", err)

	err = writeUserResponse(w, user)
	panicOnError("writeUserResponse failed", err)
}

func (srv *ConduitAppHttpHandlers) loginUser(w http.ResponseWriter, r *http.Request) {
	// TODO: read from request: email, password
	sessionID, userID := "foobar", "42"

	var user = UserProfile{
		ID:        userID,
		Email:     "golang@example.com",
		CreatedAt: now_RFC3339(),
		UpdatedAt: now_RFC3339(),
		Username:  "golang",
		Token:     sessionID, // for test machinery
	}

	err := srv.addNewSession(userID, sessionID)
	panicOnError("addNewSession failed", err)

	err = writeUserResponse(w, user)
	panicOnError("writeUserResponse failed", err)
}

func (srv *ConduitAppHttpHandlers) registerNewUser(w http.ResponseWriter, r *http.Request) {
	// TODO: read from request: email, password, username
	var user = UserProfile{
		ID:        "42",
		Email:     "golang@example.com",
		CreatedAt: now_RFC3339(),
		UpdatedAt: now_RFC3339(),
		Username:  "golang",
	}

	err := putUser2Storage(srv.storage, user)
	panicOnError("putUserByEmail failed", err)

	var u = &struct{ User UserProfile }{User: user}
	err = writeResponseWithCode(w, u, http.StatusCreated)
	panicOnError("writeResponseWithCode failed", err)
}

func (srv *ConduitAppHttpHandlers) addNewSession(userID, sessionID string) error {
	return srv.sessions.AddSession(userID, sessionID)
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

type SessionManager interface {
	AddSession(userID, sessionID string) error
	GetUserID(sessionID string) (string, error)
}

var _ SessionManager = &SimpleSessionManager{} // type check

type SimpleSessionManager struct {
	data map[string]string // not ready for async
}

func NewSimpleSessionManager() *SimpleSessionManager {
	return &SimpleSessionManager{
		data: make(map[string]string, 16),
	}
}

// AddSession implements SessionManager.
func (sm *SimpleSessionManager) AddSession(userID string, sessionID string) error {
	sm.data[sessionID] = userID
	return nil
}

// GetUserID implements SessionManager.
func (sm *SimpleSessionManager) GetUserID(sessionID string) (string, error) {
	// userID, err := srv.sessions.GetUserID(sessionID)
	uid, exists := sm.data[sessionID]
	if exists {
		return uid, nil
	}
	return "", fmt.Errorf("GetUserID failed, session %#v not exist", sessionID)
}

type Storage interface {
	// Set(key, value any) error
	// Get(key string) (any, error)
	GetAll() ([]any, error)
	SetAll([]any) error
}

var _ Storage = &RAMStorage{} // type check

type RAMStorage struct {
	data []any
}

func NewRAMStorage() *RAMStorage {
	return &RAMStorage{
		data: make([]any, 0, 16),
	}
}

// GetAll implements Storage.
func (stor RAMStorage) GetAll() ([]any, error) {
	return stor.data, nil // not ready for async
}

// SetAll implements Storage.
func (stor *RAMStorage) SetAll(items []any) error {
	stor.data = items // not ready for async
	return nil
}

// getUserFromStorage: use user.id as key to search in stored items
func getUserFromStorage(stor Storage, userID string) (UserProfile, error) {
	items, err := stor.GetAll()
	if err == nil {
		for _, x := range items {
			user := x.(UserProfile) // panic if not only users in storage
			if user.ID == userID {
				return user, nil
			}
		}
	} else {
		return UserProfile{}, err
	}
	return UserProfile{}, fmt.Errorf("getUserByEmail failed. ID %#v not found in %d records", userID, len(items))
}

// putUser2Storage: use user.id as key in stored items
func putUser2Storage(stor Storage, user UserProfile) error {
	items, err := stor.GetAll()
	if err == nil {
		items = slices.DeleteFunc(items, func(x any) bool {
			return (x.(UserProfile)).ID == user.ID // panic if not only users in storage
		})
		items = append(items, user)

		return stor.SetAll(items)
	}
	return err
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

func getSessionIDFromReq(r *http.Request) (string, error) {
	// req.Header.Add("Authorization", "Token "+tplParams[item.TokenName])
	token, found := strings.CutPrefix(
		r.Header.Get("Authorization"),
		"Token ",
	)
	if found {
		return token, nil
	}
	return "", fmt.Errorf("getSessionIDFromReq failed, token not found")
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
