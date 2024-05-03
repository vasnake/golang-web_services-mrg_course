package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"
	"slices"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// NewDbExplorer create http handler for db_explorer app
func NewDbExplorer(dbRef *sql.DB) (http.Handler, error) {
	srv := &MysqlExplorerHttpHandlers{
		DB: dbRef,
	}
	r := mux.NewRouter() // gorilla/mux
	r.HandleFunc("/", srv.DefaultRead).Methods("GET")
	return r, nil
}

type MysqlExplorerHttpHandlers struct {
	DB *sql.DB
}

type GenericMap map[string]interface{}

func (srv *MysqlExplorerHttpHandlers) DefaultRead(w http.ResponseWriter, r *http.Request) {
	// http.Error(w, "not yet", http.StatusInternalServerError)
	show("DefaultRead, request: ", r)

	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			show("DefaultRead recover from error: ", err)
			writeError(http.StatusInternalServerError, "Internal server error", w)
		}
	}()

	resultRef := &GenericMap{
		"response": GenericMap{
			"tables": []string{"items", "users"},
		},
	}

	writeSuccess(http.StatusOK, resultRef, w)
}

func writeSuccess(status int, obj interface{}, w http.ResponseWriter) {
	// bytes, err := json.Marshal(ApiSuccessResponse{"", obj})
	bytes, err := json.Marshal(obj)
	writeResponse(status, err, bytes, w)
}

func writeError(status int, message string, w http.ResponseWriter) {
	bytes, err := json.Marshal(ApiErrorResponse{message})
	writeResponse(status, err, bytes, w)
}

func writeResponse(status int, err error, msg []byte, w http.ResponseWriter) {
	panicOnError("writeResponse, got an error: ", err) // TODO: replace panic with StatusInternalServerError
	w.WriteHeader(status)
	w.Write(msg)
}

type ApiSuccessResponse struct {
	Error    string      `json:"error"`
	Response interface{} `json:"response"`
}

type ApiErrorResponse struct {
	Error string `json:"error"`
}

func getOrDefault(values url.Values, key string, defaultValue string) string {
	items, ok := values[key]
	if !ok {
		return defaultValue
	}

	if len(items) == 0 {
		return defaultValue
	}

	return items[0] // TODO: or find first not empty
}

func contains(str string, lst []string) bool {
	return slices.Contains(lst, str)
}

func split(str, sep string) []string {
	return strings.Split(str, sep)
}

// panicOnError throw the panic with given error and msg prefix, if err != nil
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
		line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		// line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
}
