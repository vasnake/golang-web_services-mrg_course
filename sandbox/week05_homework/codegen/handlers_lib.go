package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	// "runtime/debug"
	"slices"
	// "strconv"
	"strings"
	"time"
)

func writeSrvError(err error, w http.ResponseWriter) {
	switch err.(type) {
	case ApiError:
		writeError((err.(ApiError)).HTTPStatus, err.Error(), w)
	default:
		writeError(http.StatusInternalServerError, err.Error(), w)
	}
}

func writeSuccess(status int, obj interface{}, w http.ResponseWriter) {
	bs, err := json.Marshal(ApiSuccessResponse{"", obj})
	writeResponse(status, err, bs, w)
}

func writeError(status int, message string, w http.ResponseWriter) {
	bs, err := json.Marshal(ApiErrorResponse{message})
	writeResponse(status, err, bs, w)
}

func writeResponse(status int, err error, msg []byte, w http.ResponseWriter) {
	if err != nil {
		panic(err) // TODO: replace panic with StatusInternalServerError
	}
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

func isAuthenticated(r *http.Request) bool {
	return r.Header.Get("X-Auth") == "100500"
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
