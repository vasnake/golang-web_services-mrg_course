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

func SearchServerSlowMock(w http.ResponseWriter, r *http.Request) {
	time.Sleep(client.Timeout + (client.Timeout / 9)) // sleep longer than client timeout
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(w, "[]")
}

// SearchServerMock is a handler for http requests. Reads from Request, writes result to ResponseWriter
func SearchServerMock(w http.ResponseWriter, r *http.Request) {
	// show("Request: ", r)
	marshal := func(v any) []byte {
		bytes, err := json.Marshal(v)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
			return nil
		} else {
			return bytes
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8") // must be set before writing

	// different error cases

	accessTok := r.Header.Get("AccessToken")
	if accessTok == "fatal" { // imitation of fatal error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if accessTok != "good" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch r.FormValue("order_by") {
	case "42":
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "malformed json: ][")
		return
	case "37":
		bytes := marshal(SearchErrorResponse{Error: "Wrong `order_by` value: 37"})
		if bytes != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(bytes)
		}
		return
	case "73":
		bytes := marshal(SearchErrorResponse{Error: "Malformer result: list of users"})
		if bytes != nil {
			w.Write(bytes)
		}
		return
	default:
		_ = "no errors here"
	}

	switch r.FormValue("order_field") {
	case "foo":
		bytes := marshal(SearchErrorResponse{Error: "ErrorBadOrderField"})
		if bytes != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(bytes)
		}
		return
	default:
		_ = "no errors here"
	}

	// search and find users:
	bytes := marshal(found3Users)
	if bytes != nil {
		w.Write(bytes)
	}
}

func TestSmoke(t *testing.T) {
	serverMock := httptest.NewServer(http.HandlerFunc(SearchServerMock))
	defer serverMock.Close()

	sc := &SearchClient{
		URL:         serverMock.URL,
		AccessToken: "good",
	}
	req := SearchRequest{}

	resp, err := sc.FindUsers(req)

	idx := 0
	if resp == nil {
		t.Errorf("[%d] response is nil, unexpected error: %#v", idx, err)
	}
	if err != nil {
		t.Errorf("[%d] unexpected error: %#v", idx, err)
	}
}

type FindUserTEstCase struct {
	request     SearchRequest
	response    *SearchResponse
	err         error
	accessToken string
	url         string
}

func TestErrorCases(t *testing.T) {
	serverMockClosed := httptest.NewServer(http.HandlerFunc(SearchServerMock))
	serverMock := httptest.NewServer(http.HandlerFunc(SearchServerMock))
	serverMockSlow := httptest.NewServer(http.HandlerFunc(SearchServerSlowMock))
	defer func() { serverMock.Close(); serverMockSlow.Close() }()
	serverMockClosed.Close()

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
			url:      serverMockClosed.URL,
			request:  SearchRequest{},
			response: nil,
			err:      fmt.Errorf("unknown error Get"),
		},
		{
			url:      serverMockSlow.URL,
			request:  SearchRequest{},
			response: nil,
			err:      fmt.Errorf("timeout for limit"),
		},
		{
			accessToken: "bad",
			request:     SearchRequest{},
			response:    nil,
			err:         fmt.Errorf("Bad AccessToken"),
		},
		{
			accessToken: "fatal",
			request:     SearchRequest{},
			response:    nil,
			err:         fmt.Errorf("SearchServer fatal error"),
		},
		{
			request:  SearchRequest{OrderBy: 42}, // bad request, with unpackable error message
			response: nil,
			err:      fmt.Errorf("cant unpack error json: "),
		},
		{
			request:  SearchRequest{OrderField: "foo"}, // bad request, with meaningful error message
			response: nil,
			err:      fmt.Errorf("OrderFeld "),
		},
		{
			request:  SearchRequest{OrderBy: 37}, // bad request, with meaningful unknown error message
			response: nil,
			err:      fmt.Errorf("unknown bad request error: "),
		},
		{
			request:  SearchRequest{OrderBy: 73}, // good request, unpackable result
			response: nil,
			err:      fmt.Errorf("cant unpack result json: "),
		},
	}

	eeq := IsErrorStartsWith
	getToken := func(row FindUserTEstCase) string {
		if row.accessToken == "" {
			return "good"
		}
		return row.accessToken
	}
	getURL := func(row FindUserTEstCase) string {
		if row.url == "" {
			return serverMock.URL
		}
		return row.url
	}

	sc := SearchClient{}
	for idx, row := range testTable {
		sc.AccessToken = getToken(row)
		sc.URL = getURL(row)

		resp, err := sc.FindUsers(row.request)

		if !eeq(err, row.err) {
			t.Errorf("row [%d], expected error: `%#v`; got error: `%#v`", idx, row.err, err)
		}
		if resp != row.response {
			t.Errorf("row [%d], expected resp: `%#v`; got resp: `%#v`", idx, row.response, resp)
		}
	}
}

var found3Users = []User{
	{Id: 1},
	{Id: 2},
	{Id: 3},
}

func TestHappyCases(t *testing.T) {
	testTable := []FindUserTEstCase{
		{ // if len(data) == req.Limit
			request:  SearchRequest{Limit: 2}, // limit: 3
			response: &SearchResponse{NextPage: true, Users: found3Users[:2]},
		},
		{ // if len(data) != req.Limit
			request:  SearchRequest{Limit: 3}, // limit: 4
			response: &SearchResponse{Users: found3Users[:3]},
		},
		{ // if len(data) != req.Limit
			request:  SearchRequest{Limit: 33}, // limit: 26
			response: &SearchResponse{Users: found3Users[:3]},
		},
	}

	serverMock := httptest.NewServer(http.HandlerFunc(SearchServerMock))
	defer serverMock.Close()
	sc := SearchClient{URL: serverMock.URL, AccessToken: "good"}
	eeq := IsErrorsEqual
	req := IsResponseEqual

	for idx, row := range testTable {

		resp, err := sc.FindUsers(row.request)

		if !eeq(err, row.err) {
			t.Errorf("row [%d], expected error: `%#v`; got error: `%#v`", idx, row.err, err)
		}
		if !req(resp, row.response) {
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

func IsResponseEqual(r1, r2 *SearchResponse) bool {
	// first field
	if r1.NextPage != r2.NextPage {
		return false
	}

	// last field

	if r1.Users == nil && r2.Users == nil {
		return true
	}
	if r1.Users == nil || r2.Users == nil {
		return false
	}
	if len(r1.Users) != len(r2.Users) {
		return false
	}

	for i, u1 := range r1.Users {
		if !IsUserEqual(u1, r2.Users[i]) {
			return false
		}
	}
	return true
}

func IsUserEqual(u1, u2 User) bool {
	return u1 == u2
	// return u1.Id == u2.Id &&
	// 	u1.Age == u2.Age &&
	// 	u1.About == u2.About &&
	// 	u1.Gender == u2.Gender &&
	// 	u1.Name == u2.Name
}
