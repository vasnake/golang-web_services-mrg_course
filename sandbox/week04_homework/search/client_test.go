package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSmoke(t *testing.T) {

	var serverLogic = func(w http.ResponseWriter, r *http.Request) {
		// show("Request: ", r)
		switch r.FormValue("id") {
		case "42":
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `{"status": 200, "balance": 100500}`)
		case "__internal_error":
			fallthrough
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	serverMock := httptest.NewServer(http.HandlerFunc(serverLogic))
	defer serverMock.Close()
	sc := &SearchClient{
		URL: serverMock.URL,
	}

	idx := 1
	req := SearchRequest{}
	resp, err := sc.FindUsers(req)
	if err != nil || resp == nil {
		t.Errorf("[%d] unexpected error: %#v", idx, err)
	}
}
