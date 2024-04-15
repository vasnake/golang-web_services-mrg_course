package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func SearchServerSlow(w http.ResponseWriter, r *http.Request) {
	// time.Sleep(1100 * time.Millisecond)
	time.Sleep(client.Timeout + 42)
	switch r.FormValue("id") {
	case "__internal_error":
		w.WriteHeader(http.StatusServiceUnavailable)
	default:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		io.WriteString(w, "[]")
	}
}

// SearchServer is a handler for http requests. Reads from Request, writes result to ResponseWriter
func SearchServer(w http.ResponseWriter, r *http.Request) {
	// show("Request: ", r)
	data := []User{}

	at := r.Header.Get("AccessToken")
	if at != "good" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch r.FormValue("id") {
	case "42":
		w.WriteHeader(http.StatusPaymentRequired)
		// w.WriteHeader(http.StatusOK)
		// io.WriteString(w, `{"status": 200, "balance": 100500}`)

	case "__internal_error":
		w.WriteHeader(http.StatusServiceUnavailable)
		// fallthrough

	default:
		bytes, err := json.Marshal(data)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8") // must be set before writing
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
		} else {
			w.Header().Set("Content-Type", "application/json; charset=utf-8") // must be set before writing
			w.Write(bytes)
		}
	} // end switch
}

func TestSmoke(t *testing.T) {
	serverMock := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer serverMock.Close()

	sc := &SearchClient{
		URL:         serverMock.URL,
		AccessToken: "good",
	}
	req := SearchRequest{}

	resp, err := sc.FindUsers(req)

	idx := 0
	if err == nil && resp == nil {
		t.Errorf("[%d] response is nil, unexpected error: %#v", idx, err)
	} else if err != nil {
		t.Errorf("[%d] unexpected error: %#v", idx, err)
	}
}

type FindUserTEstCase struct {
	request  SearchRequest
	response *SearchResponse
	err      error
}

func TestErrorCases(t *testing.T) {
	testTable := []FindUserTEstCase{
		{
			request:  SearchRequest{Limit: -1},
			response: nil,
			err:      fmt.Errorf("limit must be > 0"),
		},
		{
			request:  SearchRequest{Offset: -1},
			response: nil,
			err:      fmt.Errorf("offset must be > 0"),
		},
		{
			request:  SearchRequest{},
			response: nil,
			err:      fmt.Errorf("unknown error Get"),
		},
		{
			request:  SearchRequest{},
			response: nil,
			err:      fmt.Errorf("timeout for limit"),
		},
		{
			request:  SearchRequest{},
			response: nil,
			err:      fmt.Errorf("Bad AccessToken"),
		},
		// http.StatusInternalServerError:  return nil, fmt.Errorf("SearchServer fatal error")
		// http.StatusBadRequest:
		// 		fmt.Errorf("unknown bad request error: %s", errResp.Error)
		// 		fmt.Errorf("cant unpack error json: %s", err)
		// 		fmt.Errorf("OrderFeld %s invalid", req.OrderField)
		// return nil, fmt.Errorf("cant unpack result json: %s", err)
	}

	serverMockClosed := httptest.NewServer(http.HandlerFunc(SearchServer))
	serverMockClosed.Close()
	serverMockSlow := httptest.NewServer(http.HandlerFunc(SearchServerSlow))
	serverMock := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer func() { serverMock.Close(); serverMockSlow.Close() }()
	sc := SearchClient{}

	for idx, row := range testTable {
		// update env config
		sc.AccessToken = "good"
		sc.URL = serverMock.URL
		eq := func(e1, e2 error) bool { return IsErrorsEqual(e1, e2) }
		if row.err.Error() == "Bad AccessToken" {
			sc.AccessToken = "bad"
		}
		if row.err.Error() == "timeout for limit" {
			sc.URL = serverMockSlow.URL
			eq = func(e1, e2 error) bool { return IsErrorStartsWith(e1, e2) }
		}
		if row.err.Error() == "unknown error Get" {
			sc.URL = serverMockClosed.URL
			eq = func(e1, e2 error) bool { return IsErrorStartsWith(e1, e2) }
		}

		resp, err := sc.FindUsers(row.request)

		if !eq(err, row.err) {
			t.Errorf("row [%d], expected error: `%#v`; got error: `%#v`", idx, row.err, err)
		}
		if resp != row.response {
			t.Errorf("row [%d], expected resp: `%#v`; got resp: `%#v`", idx, row.response, resp)
		}
	}
}

func TestHappyCases(t *testing.T) {
	testTable := []FindUserTEstCase{}

	serverMock := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer serverMock.Close()
	sc := SearchClient{URL: serverMock.URL}
	var eq = func(e1, e2 error) bool { return IsErrorsEqual(e1, e2) }

	for idx, row := range testTable {
		resp, err := sc.FindUsers(row.request)
		if eq(err, row.err) {
			t.Errorf("row [%d], expected error: `%#v`; got error: `%#v`", idx, row.err, err)
		}
		if resp != row.response {
			t.Errorf("row [%d], expected resp: `%#v`; got resp: `%#v`", idx, row.response, resp)
		}
	}
}

func show(msg string, xs ...any) {
	var line = ts() + ": " + msg

	for _, x := range xs {
		// https://pkg.go.dev/fmt
		line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		// line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
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

func IsErrorsEqual(e1, e2 error) bool {
	if e1 == nil {
		return e2 == nil
	}
	if e2 == nil {
		return e1 == nil
	}
	return e1.Error() == e2.Error()
}

func IsErrorStartsWith(e1, e2 error) bool {
	if e1 == nil && e2 == nil {
		return true
	}
	if e1 == nil || e2 == nil {
		return false
	}
	return strings.Contains(e1.Error(), e2.Error())
}
