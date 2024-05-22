package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"
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

// ServeHTTP: implements http.Handler; routing based on http method name
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
	case "/api/articles":
		srv.createNewArticle(w, r)

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
	sessionID, err := getSessionIDFromReq(r)
	panicOnError("getSessionIDFromReq failed", err)
	userID, err := srv.sessions.GetUserID(sessionID)
	panicOnError("GetUserID failed", err)
	user, err := getUserFromStorage(srv.storage, userID)
	panicOnError("getUserFromStorage failed", err)

	bodyMap, err := unmarshalBody(r)
	panicOnError("unmarshalBody failed", err)

	userData, userExists := bodyMap["user"]
	if !userExists {
		show("updateCurrentUser, no user in given data")
		http.Error(w, "oops", http.StatusBadRequest)
	}

	userMap, err := giveMeStrings(userData.(map[string]any)) // recover from panic
	panicOnError("giveMeStrings failed", err)

	user.updateFromMap(userMap)

	err = putUser2Storage(srv.storage, user)
	panicOnError("putUserByEmail failed", err)

	user.Token = sessionID
	err = writeUserResponse(w, user)
	panicOnError("writeUserResponse failed", err)
}

func (srv *ConduitAppHttpHandlers) showCurrentUser(w http.ResponseWriter, r *http.Request) {
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
	bodyMap, err := unmarshalBody(r)
	panicOnError("unmarshalBody failed", err)
	userData, userExists := bodyMap["user"]
	if !userExists {
		show("loginUser, no user in given data")
		http.Error(w, "oops", http.StatusBadRequest)
	}

	userMap, err := giveMeStrings(userData.(map[string]any)) // recover from panic
	panicOnError("giveMeStrings failed", err)

	var dubiousUser = UserProfile{}
	dubiousUser.updateFromMap(userMap) // email, password
	// end of reading parameters.

	// load user from db and check

	goodUser, err := getUserFromStorageByEmail(srv.storage, dubiousUser.Email)
	panicOnError("getUserFromStorageByEmail failed", err)
	// TODO: check password

	sessionID := srv.newSessionID()
	err = srv.addNewSession(goodUser.ID, sessionID)
	panicOnError("addNewSession failed", err)
	goodUser.Token = sessionID

	err = writeUserResponse(w, goodUser)
	panicOnError("writeUserResponse failed", err)
}

func (srv *ConduitAppHttpHandlers) registerNewUser(w http.ResponseWriter, r *http.Request) {
	bodyMap, err := unmarshalBody(r)
	panicOnError("unmarshalBody failed", err)
	userData, userExists := bodyMap["user"]
	if !userExists {
		show("registerNewUser, no user in given data")
		http.Error(w, "oops", http.StatusBadRequest)
	}
	userMap, err := giveMeStrings(userData.(map[string]any)) // recover from panic
	panicOnError("giveMeStrings failed", err)

	sessionID, userID := srv.newSessionID(), srv.newUserID()
	var user = UserProfile{
		ID:        userID,
		CreatedAt: strRef(now_RFC3339()),
	}
	user.updateFromMap(userMap)
	// end of loading given params.

	// save
	err = putUser2Storage(srv.storage, user)
	panicOnError("putUser2Storage failed", err)

	// auth
	err = srv.addNewSession(userID, sessionID)
	panicOnError("addNewSession failed", err)
	user.Token = sessionID

	// response
	var u = &struct{ User UserProfile }{User: user}
	err = writeResponseWithCode(w, u, http.StatusCreated)
	panicOnError("writeResponseWithCode failed", err)
}

func (srv *ConduitAppHttpHandlers) createNewArticle(w http.ResponseWriter, r *http.Request) {
	// load params
	bodyMap, err := unmarshalBody(r)
	panicOnError("unmarshalBody failed", err)
	articleData, articleExists := bodyMap["article"]
	if !articleExists {
		show("createNewArticle, no article in given data")
		http.Error(w, "oops", http.StatusBadRequest)
	}

	// decode article
	var article = PostedArticle{
		CreatedAt: now_RFC3339(),
	}
	err = article.updateFromMap(articleData.(map[string]any))
	panicOnError("article.loadFromJson failed", err)

	// load user
	sessionID, err := getSessionIDFromReq(r)
	panicOnError("getSessionIDFromReq failed", err)
	userID, err := srv.sessions.GetUserID(sessionID)
	panicOnError("GetUserID failed", err)
	user, err := getUserFromStorage(srv.storage, userID)
	panicOnError("getUserFromStorage failed", err)

	// no idea why is that. I'm to lazy to search for contract specs
	user.Email = ""
	user.CreatedAt = nil
	user.UpdatedAt = nil

	article.Author = user // show("user: ", user.Email)
	article.Slug = srv.newArticleID()

	// save
	err = putArticle2Storage(srv.storage, article)
	panicOnError("putArticle2Storage failed", err)

	// response
	var a = &struct{ Article PostedArticle }{Article: article}
	err = writeResponseWithCode(w, a, http.StatusCreated)
	panicOnError("writeResponseWithCode failed", err)
}

func (srv *ConduitAppHttpHandlers) addNewSession(userID, sessionID string) error {
	return srv.sessions.AddSession(userID, sessionID)
}

func (srv *ConduitAppHttpHandlers) newSessionID() string {
	return nextID()
}

func (srv *ConduitAppHttpHandlers) newUserID() string {
	return nextID()
}

func (srv *ConduitAppHttpHandlers) newArticleID() string {
	return nextID()
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
	GetAllUsers() ([]any, error)
	SetAllUsers([]any) error
	GetAllArticles() ([]any, error)
	SetAllArticles([]any) error
}

var _ Storage = &RAMStorage{} // type check

type RAMStorage struct {
	usersData    []any
	articlesData []any
}

func NewRAMStorage() *RAMStorage {
	return &RAMStorage{
		usersData:    make([]any, 0, 16),
		articlesData: make([]any, 0, 16),
	}
}

// GetAllUsers implements Storage.
func (stor RAMStorage) GetAllUsers() ([]any, error) {
	return stor.usersData, nil // not ready for async
}

// SetAllUsers implements Storage.
func (stor *RAMStorage) SetAllUsers(items []any) error {
	stor.usersData = items // not ready for async
	return nil
}

// GetAllArticles implements Storage.
func (stor *RAMStorage) GetAllArticles() ([]any, error) {
	return stor.articlesData, nil // not async code
}

// SetAllArticles implements Storage.
func (stor *RAMStorage) SetAllArticles(items []any) error {
	stor.articlesData = items // not async code
	return nil
}

// getUserFromStorage: use user.id as key to search in stored items
func getUserFromStorage(stor Storage, userID string) (UserProfile, error) {
	items, err := stor.GetAllUsers()
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
	return UserProfile{}, fmt.Errorf("getUserFromStorage failed. ID %#v not found in %d records", userID, len(items))
}

// getUserFromStorageByEmail: use email as key to search in stored items
func getUserFromStorageByEmail(stor Storage, email string) (UserProfile, error) {
	items, err := stor.GetAllUsers()
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
	return UserProfile{}, fmt.Errorf("getUserFromStorageByEmail failed. email %#v not found in %d records", email, len(items))
}

// putUser2Storage: use user.id as key in stored items
func putUser2Storage(stor Storage, user UserProfile) error {
	items, err := stor.GetAllUsers()
	if err == nil {
		items = slices.DeleteFunc(items, func(x any) bool {
			return (x.(UserProfile)).ID == user.ID
		})
		items = append(items, user)

		return stor.SetAllUsers(items)
	}
	return err
}

// putArticle2Storage: use article.slug as key
func putArticle2Storage(stor Storage, article PostedArticle) error {
	items, err := stor.GetAllArticles()
	if err == nil {
		items = slices.DeleteFunc(items, func(x any) bool {
			return (x.(PostedArticle)).Slug == article.Slug
		})
		items = append(items, article)

		return stor.SetAllArticles(items)
	}
	return err
}

// getArticleFromStorage: use article.slug as key
func getArticleFromStorage(stor Storage, slug string) (PostedArticle, error) {
	items, err := stor.GetAllArticles()
	if err == nil {
		for _, x := range items {
			article := x.(PostedArticle)
			if article.Slug == slug {
				return article, nil
			}
		}
	} else {
		return PostedArticle{}, err
	}
	return PostedArticle{}, fmt.Errorf("getArticleFromStorage failed. slug %#v not found in %d records", slug, len(items))
}

// time as `RFC3339     = "2006-01-02T15:04:05Z07:00"`
type UserProfile struct {
	ID        string  `json:"id" testdiff:"ignore"`
	Email     string  `json:"email"`
	CreatedAt *string `json:"createdAt"`
	UpdatedAt *string `json:"updatedAt"`
	Username  string  `json:"username"`
	Bio       string  `json:"bio"`
	Image     string  `json:"image"`
	Token     string  `json:"token" testdiff:"ignore"`
	Following bool
	password  string
}

func (user *UserProfile) updateFromMap(data map[string]string) *UserProfile {
	// TODO: add other user fields
	for k, v := range data {
		switch k {
		case "username":
			user.Username = v
		case "password":
			user.password = v
		case "email":
			user.Email = v
		case "bio":
			user.Bio = v
		default:
			show("unknown user field: ", k, v)
		}
	}
	user.UpdatedAt = strRef(now_RFC3339()) // TODO: if updated

	return user
}

type PostedArticle struct {
	// simple decoding
	Slug        string `json:"slug" testdiff:"ignore"` // id
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	// complicated decoding
	Author         UserProfile `json:"author"`
	TagList        []string    `json:"tagList"`
	FavoritesCount int         `json:"favoritesCount"`
	Favorited      bool        `json:"favorited"`
}

func (a *PostedArticle) updateFromMap(data map[string]any) error {
	var _str = func(x any) string {
		return x.(string)
	}
	var _list = func(x any) []string {
		var list = make([]string, 0, 4)
		for _, x := range x.([]any) {
			list = append(list, x.(string))
		}
		return list
	}

	for k, v := range data {
		switch k {
		// case "slug":
		case "title":
			a.Title = _str(v)
		case "description":
			a.Description = _str(v)
		case "body":
			a.Body = _str(v)
		case "tagList":
			a.TagList = _list(v)
		default:
			show("unknown article field: ", k, v)
		}
	}
	a.UpdatedAt = now_RFC3339() // TODO: if updated

	return nil // TODO: if no fields found: error
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

var globalCounter = new(atomic.Uint64)

func nextID() string {
	return strconv.FormatInt(int64(globalCounter.Add(1)), 36)
}

func giveMeStrings(xs map[string]any) (map[string]string, error) {
	var ys = make(map[string]string, 16)
	for k, v := range xs {
		ys[k] = v.(string) // panic
	}
	return ys, nil
}

func panicOnError(msg string, err error) {
	if err != nil {
		panic(msg + ": " + err.Error())
	}
}

func strRef(in string) *string {
	return &in
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
