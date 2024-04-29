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
// fillFrom write data from 'params' to 'pref'
func (pref *ProfileParams) fillFrom(params url.Values) error {
	var err error = nil

	pref.Login = getOrDefault(params, "login", "")

	return err
}
// validate check data against set of rules
func (cpref *ProfileParams) validate() error {


	if cpref.Login == "" { // required
		return errors.New("login: value required")
	}




	return nil
}
// fillFrom write data from 'params' to 'pref'
func (pref *CreateParams) fillFrom(params url.Values) error {
	var err error = nil

	pref.Login = getOrDefault(params, "login", "")

	pref.Name = getOrDefault(params, "full_name", "")

	pref.Status = getOrDefault(params, "status", "user")

	pref.Age, err = strconv.Atoi(getOrDefault(params, "age", ""))
	if err != nil {
		return errors.New("age must be int")
	}


	return err
}
// validate check data against set of rules
func (cpref *CreateParams) validate() error {


	if cpref.Login == "" { // required
		return errors.New("login: value required")
	}


	if len(cpref.Login) < 10 { // min string
		return errors.New("login len must be >= 10")
	}











	if !contains(cpref.Status, []string{"user", "moderator", "admin"}) { // enum
		return errors.New("status must be one of [user, moderator, admin]")
	}


	if cpref.Age < 0 { // min int
		return errors.New("age must be >= 0")
	}


	if cpref.Age > 128 {
		return errors.New("age must be <= 128")
	}

	return nil
}
// fillFrom write data from 'params' to 'pref'
func (pref *OtherCreateParams) fillFrom(params url.Values) error {
	var err error = nil

	pref.Username = getOrDefault(params, "username", "")

	pref.Name = getOrDefault(params, "account_name", "")

	pref.Class = getOrDefault(params, "class", "warrior")

	pref.Level, err = strconv.Atoi(getOrDefault(params, "level", ""))
	if err != nil {
		return errors.New("level must be int")
	}


	return err
}
// validate check data against set of rules
func (cpref *OtherCreateParams) validate() error {


	if cpref.Username == "" { // required
		return errors.New("username: value required")
	}


	if len(cpref.Username) < 3 { // min string
		return errors.New("username len must be >= 3")
	}











	if !contains(cpref.Class, []string{"warrior", "sorcerer", "rouge"}) { // enum
		return errors.New("class must be one of [warrior, sorcerer, rouge]")
	}


	if cpref.Level < 1 { // min int
		return errors.New("level must be >= 1")
	}


	if cpref.Level > 50 {
		return errors.New("level must be <= 50")
	}

	return nil
}
 
