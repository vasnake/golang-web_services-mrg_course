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
	time.Sleep(client.Timeout + (111 * 42))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(w, "[]")
}

// SearchServer is a handler for http requests. Reads from Request, writes result to ResponseWriter
func SearchServer(w http.ResponseWriter, r *http.Request) {
	// show("Request: ", r)
	/*
	   2024-04-16T06:55:24.553Z: Request: *http.Request(&{
	   	GET /?
	   		limit=1&
	   		offset=0&
	   		order_by=42&
	   		order_field=&
	   		query=
	   	HTTP/1.1 1 1 map[
	   		Accept-Encoding:[gzip]
	   		Accesstoken:[good]
	   		User-Agent:[Go-http-client/1.1]]
	   	{} <nil> 0 [] false 127.0.0.1:33707 map[] map[] <nil> map[] 127.0.0.1:59348 /?
	   		limit=1&
	   		offset=0&
	   		order_by=42&
	   		order_field=&
	   		query= <nil> <nil> <nil> 0xc000212190 <nil> [] map[]});
	*/
	data := []User{}

	accessTok := r.Header.Get("AccessToken")
	if accessTok == "fatal" { // TODO: remove this mock of fatal error
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
		w.Header().Set("Content-Type", "application/json; charset=utf-8") // must be set before writing
		io.WriteString(w, "unpackable json: ][")
		return
	case "37":
		w.Header().Set("Content-Type", "application/json; charset=utf-8") // must be set before writing
		bytes, err := json.Marshal(SearchErrorResponse{Error: "Wrong `order_by` value: 37"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(bytes)
		}
		return
	case "73":
		w.Header().Set("Content-Type", "application/json; charset=utf-8") // must be set before writing
		bytes, err := json.Marshal(SearchErrorResponse{Error: "Malformer result: list of users"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
		} else {
			w.Write(bytes)
		}
		return

	case "-1":
		fallthrough
	case "0":
		fallthrough
	case "1":
		fallthrough
	default:
		_ = "foo"
	}

	switch r.FormValue("order_field") {
	case "foo":
		w.Header().Set("Content-Type", "application/json; charset=utf-8") // must be set before writing
		bytes, err := json.Marshal(SearchErrorResponse{Error: "ErrorBadOrderField"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(bytes)
		}
		return

	default:
		_ = "foo"
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
	request     SearchRequest
	response    *SearchResponse
	err         error
	accessToken string
	url         string
}

func TestErrorCases(t *testing.T) {
	serverMockClosed := httptest.NewServer(http.HandlerFunc(SearchServer))
	serverMock := httptest.NewServer(http.HandlerFunc(SearchServer))
	serverMockSlow := httptest.NewServer(http.HandlerFunc(SearchServerSlow))
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

	eq := IsErrorStartsWith
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
