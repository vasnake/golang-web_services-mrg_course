package token

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"photolist/pkg/session"
)

type TokenManager interface {
	Create(*session.Session, int64) (string, error)
	Check(*session.Session, string) (bool, error)
}

var (
	noCSRFUrls = map[string]struct{}{
		"/user/login_oauth": struct{}{},
		"/user/login":       struct{}{},
		"/user/reg":         struct{}{},
		"/api/v1/token":     struct{}{},

		// TODO enable csrf for graphql after playground done
		"/playgroud": struct{}{},
		"/graphql":   struct{}{},
	}

	errorTokenExpired = errors.New("token expired")
)

func CsrfTokenMiddleware(tm TokenManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, skip := noCSRFUrls[r.URL.Path]
		isAPI := strings.HasPrefix(r.URL.Path, "/api")
		skip = skip || (!isAPI && r.Method == http.MethodGet) // check all api and regular forms

		if skip {
			next.ServeHTTP(w, r)
			return
		}

		var CSRFToken string
		if isAPI {
			CSRFToken = r.Header.Get("csrf-token")
		} else {
			CSRFToken = r.FormValue("csrf-token")
		}

		sess, _ := session.SessionFromContext(r.Context())
		tokenValid, err := tm.Check(sess, CSRFToken)
		if tokenValid {
			next.ServeHTTP(w, r)
			return
		}

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
