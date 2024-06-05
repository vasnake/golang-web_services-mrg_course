package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/99designs/gqlgen/graphql"
)

// auth user stuff

// CheckAuthorizedMiddleware: graphql middleware (cfg.directives)
func CheckAuthorizedMiddleware(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	session, err := SessionFromContext(ctx)
	if err != nil {
		return nil, err
	}
	show("CheckAuthorizedMiddleware, session: ", session)
	return next(ctx)
}

// auth middleware: add to context session data
func (udb *UserSessionAuth) InjectSession2ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerAuthValue := r.Header.Get("Authorization")
		if headerAuthValue != "" {
			sessToken := strings.TrimPrefix(headerAuthValue, "Token ")
			if len(sessToken) != len(headerAuthValue) && sessToken != "" {
				user, err := udb.GetUserBySession(sessToken)
				if err == nil {
					var sess AppSession = user
					ctx := SessionToContext(r.Context(), sess)

					show("InjectSession2ContextMiddleware, session: ", sess)
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			}
		}

		// show("InjectSession2ContextMiddleware, no session")
		next.ServeHTTP(w, r)
	})
}

type UserSessionAuth struct {
	usersTable []*AppUserRow
}

func (udb *UserSessionAuth) New() *UserSessionAuth {
	return &UserSessionAuth{
		usersTable: make([]*AppUserRow, 0, 16),
	}
}

func (udb *UserSessionAuth) UserExists(aur *AppUserRow) bool {
	for _, ur := range udb.usersTable {
		if ur.username == aur.username {
			return true
		}
	}
	return false
}

func (udb *UserSessionAuth) createNewSession(user *AppUserRow) (updatedUser *AppUserRow, token string) {
	token = nextID_36()
	updatedUser = user
	updatedUser.sessions = append(updatedUser.sessions, token)
	return
}

func (udb *UserSessionAuth) insertUser(user *AppUserRow) {
	udb.usersTable = append(udb.usersTable, user)
}

func (udb *UserSessionAuth) RegisterNewUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "RegisterNewUserHandler, only POST method accepted", http.StatusBadRequest)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	panicOnError("RegisterNewUserHandler failed, can't read req. body", err)
	r.Body.Close()

	newUser, err := loadUserFromJsonBytes(bodyBytes)
	panicOnError("RegisterNewUserHandler failed, can't decode user struct", err)

	if udb.UserExists(newUser) {
		http.Error(w, "RegisterNewUserHandler, user exists already", http.StatusBadRequest)
		return
	}

	newUser, token := udb.createNewSession(newUser)
	newUser.ID = token // TODO: rewrite this shit (create new user, create new session(userid), don't use user as sessions storage)
	udb.insertUser(newUser)

	respContent := make(map[string]any, 1)
	respContent["token"] = token
	writeJsonResponse(w, respContent, "")
}

func (udb *UserSessionAuth) GetUserBySession(sessToken string) (*AppUserRow, error) {
	for _, ur := range udb.usersTable {
		if slices.Contains(ur.sessions, sessToken) {
			return ur, nil
		}
	}
	return nil, fmt.Errorf("GetUserBySession failed, can't find user with session %s", sessToken)
}

func loadUserFromJsonBytes(jsonBytes []byte) (*AppUserRow, error) {
	data := make(map[string]any, 16)
	err := json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return nil, fmt.Errorf("loadUserFromJsonBytes failed, can't unmarshal json: %w", err)
	}

	userAny, userExists := data["user"]
	if !userExists {
		return nil, fmt.Errorf("loadUserFromJsonBytes failed, can't find user data in given json %s", string(jsonBytes))
	}

	userMap, isMap := userAny.(map[string]any)
	if !isMap {
		return nil, fmt.Errorf("loadUserFromJsonBytes, malformed user data %s", string(jsonBytes))
	}

	email, err := loadStringFromMap(userMap, "email")
	panicOnError("no email", err)
	username, err := loadStringFromMap(userMap, "username")
	panicOnError("no username", err)
	password, err := loadStringFromMap(userMap, "password")
	panicOnError("no password", err)

	user := AppUserRow{
		email:    email,
		username: username,
		password: password,
		sessions: make([]string, 0, 16),
	}

	return &user, nil
}

func SessionFromContext(ctx context.Context) (AppSession, error) {
	sess, isSession := ctx.Value(CONTEXT_SESSION_KEY).(AppSession)
	if !isSession || sess == nil {
		return nil, ERROR_NO_USER
	} else {
		return sess, nil
	}
}

func SessionToContext(ctx context.Context, sess AppSession) context.Context {
	return context.WithValue(ctx, CONTEXT_SESSION_KEY, sess)
}

type AppSession interface{}

const (
	CONTEXT_SESSION_KEY = "SESSION_CONTEXT_KEY"
)

var (
	ERROR_NO_USER = errors.New("User not authorized")
)

type AppUserRow struct {
	ID       string
	email    string
	username string
	password string
	sessions []string
}
