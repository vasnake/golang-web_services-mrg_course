package main

import (
	"errors"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
) 
func (srv *MyApi ) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			writeError(http.StatusInternalServerError, "Internal server error", w)
		}
	}()

	switch r.URL.Path {

	case "/user/profile":
		srv.handlerProfile(w, r)

	case "/user/create":
		srv.handlerCreate(w, r)

	default:
		writeError(http.StatusNotFound, "unknown method", w)
	}
}
func (srv *OtherApi ) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			writeError(http.StatusInternalServerError, "Internal server error", w)
		}
	}()

	switch r.URL.Path {

	case "/user/create":
		srv.handlerCreate(w, r)

	default:
		writeError(http.StatusNotFound, "unknown method", w)
	}
}
 
