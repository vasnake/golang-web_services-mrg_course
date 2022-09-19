package main

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

// CSRF token support

type TokenManager interface {
	Create(*Session, int64) (string, error)
	Check(*Session, string) (bool, error)
}

var (
	noCSRFUrls = map[string]struct{}{
		"/user/login_oauth": struct{}{},
		"/user/login":       struct{}{},
		"/user/reg":         struct{}{},
		"/api/v1/token":     struct{}{},
	}

	errorTokenExpired = errors.New("token expired")
)

func CsrfTokenMiddleware(tm TokenManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, skip := noCSRFUrls[r.URL.Path] // if key in map
		isAPI := strings.HasPrefix(r.URL.Path, "/api")

		// skip declared, skip GET-and-not-API
		skip = skip || (!isAPI && r.Method == http.MethodGet) // check all api and regular forms

		if skip {
			next.ServeHTTP(w, r)
			return
		}

		// get token
		var CSRFToken string
		if isAPI {
			CSRFToken = r.Header.Get("csrf-token")
		} else {
			CSRFToken = r.FormValue("csrf-token")
		}

		sess, _ := SessionFromContext(r.Context())

		// check token, if OK => serve request
		tokenValid, err := tm.Check(sess, CSRFToken)
		if tokenValid {
			next.ServeHTTP(w, r)
			return
		}

		// process errors

		if err == errorTokenExpired {
			log.Println("csrf token expired,", sess)
			if isAPI {
				w.Header().Add("Content-Type", "application/json")
				http.Error(w, `{"error": "token expired"}`, http.StatusForbidden)
			} else {
				http.Error(w, "Token expired", http.StatusForbidden)
			}
			return
		}

		log.Println("bad token", tokenValid, err, sess)
		if isAPI {
			w.Header().Add("Content-Type", "application/json")
			http.Error(w, `{"error": "bad token"}`, http.StatusForbidden)
		} else {
			http.Error(w, "Bad token", http.StatusForbidden)
		}
	})
}
