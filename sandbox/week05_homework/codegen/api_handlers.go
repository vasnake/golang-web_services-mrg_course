package main

import (
	"errors"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
) 
// ServeHTTP implements http.Handler
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
// ServeHTTP implements http.Handler
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
// handlerProfile implements http.Handler for 'Profile' method
func (srv *MyApi ) handlerProfile(w http.ResponseWriter, r *http.Request) {
 

	r.ParseForm()
	paramsRef := new(ProfileParams)
	err := paramsRef.fillFrom(r.Form)
	if err != nil {
		writeError(http.StatusBadRequest, err.Error(), w)
		return
	}

	err = paramsRef.validate()
	if err != nil {
		writeError(http.StatusBadRequest, err.Error(), w)
		return
	}

	resultRef, err := srv.Profile(r.Context(), *paramsRef)
	if err != nil {
		writeSrvError(err, w)
		return
	}

	writeSuccess(http.StatusOK, resultRef, w)
}
// handlerCreate implements http.Handler for 'Create' method
func (srv *MyApi ) handlerCreate(w http.ResponseWriter, r *http.Request) {
 
	if r.Method != "POST" {
		writeError(http.StatusNotAcceptable, "bad method", w)
		return
	}


	if !isAuthenticated(r) {
		writeError(http.StatusForbidden, "unauthorized", w)
		return
	}

	r.ParseForm()
	paramsRef := new(CreateParams)
	err := paramsRef.fillFrom(r.Form)
	if err != nil {
		writeError(http.StatusBadRequest, err.Error(), w)
		return
	}

	err = paramsRef.validate()
	if err != nil {
		writeError(http.StatusBadRequest, err.Error(), w)
		return
	}

	resultRef, err := srv.Create(r.Context(), *paramsRef)
	if err != nil {
		writeSrvError(err, w)
		return
	}

	writeSuccess(http.StatusOK, resultRef, w)
}
// handlerCreate implements http.Handler for 'Create' method
func (srv *OtherApi ) handlerCreate(w http.ResponseWriter, r *http.Request) {
 
	if r.Method != "POST" {
		writeError(http.StatusNotAcceptable, "bad method", w)
		return
	}


	if !isAuthenticated(r) {
		writeError(http.StatusForbidden, "unauthorized", w)
		return
	}

	r.ParseForm()
	paramsRef := new(OtherCreateParams)
	err := paramsRef.fillFrom(r.Form)
	if err != nil {
		writeError(http.StatusBadRequest, err.Error(), w)
		return
	}

	err = paramsRef.validate()
	if err != nil {
		writeError(http.StatusBadRequest, err.Error(), w)
		return
	}

	resultRef, err := srv.Create(r.Context(), *paramsRef)
	if err != nil {
		writeSrvError(err, w)
		return
	}

	writeSuccess(http.StatusOK, resultRef, w)
}






 
