package main

import (
	"net/http"
)

// redirect to login or user page, according to session value in context (set by middleware)
func Index(w http.ResponseWriter, r *http.Request) {
	_, err := SessionFromContext(r.Context())
	if err != nil {
		http.Redirect(w, r, "/user/login", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/photos/", http.StatusFound)
}
