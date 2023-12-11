package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func GetUserHttpHandler(w http.ResponseWriter, r *http.Request) {
	// Some business logic, should be tested. How? Like any other function.
	key := r.FormValue("id")
	if key == "42" {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"status": 200, "resp": {"user": 42}}`)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status": 500, "err": "db_error"}`)
	}
}

func TestGetUserHttpHandler(t *testing.T) {
	type TestCase struct {
		ID                 string
		ExpectedResponse   string
		ExpectedStatusCode int
	}

	cases := []TestCase{
		{
			ID:                 "42",
			ExpectedResponse:   `{"status": 200, "resp": {"user": 42}}`,
			ExpectedStatusCode: http.StatusOK,
		},
		{
			ID:                 "43",
			ExpectedResponse:   `{"status": 500, "err": "db_error"}`,
			ExpectedStatusCode: http.StatusInternalServerError,
		},
	}

	for idx, testCase := range cases {
		req := httptest.NewRequest("GET", "http://example.com/api/user?id="+testCase.ID, nil)
		respWriter := httptest.NewRecorder()

		// show("Test request: ", req)
		GetUserHttpHandler(respWriter, req)

		if respWriter.Code != testCase.ExpectedStatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				idx, respWriter.Code, testCase.ExpectedStatusCode)
		}

		resp := respWriter.Result()
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyStr := string(bodyBytes)
		if bodyStr != testCase.ExpectedResponse {
			t.Errorf("[%d] wrong Response: got %+v, expected %+v",
				idx, bodyStr, testCase.ExpectedResponse)
		}
	}
}
