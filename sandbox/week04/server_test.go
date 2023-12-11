package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type Cart struct {
	PaymentApiURL string
}

type CheckoutResult struct {
	Status  int
	Balance int
	Err     string
}

func (cart *Cart) CartCheckout(id string) (*CheckoutResult, error) {
	// some business logic for testing

	url := cart.PaymentApiURL + "?id=" + id

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dataBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	resultRef := &CheckoutResult{}
	err = json.Unmarshal(dataBytes, resultRef)
	if err != nil {
		return nil, err
	}
	return resultRef, nil
}

func TestCartCheckout(t *testing.T) {
	type TestCase struct {
		ID              string
		ExpectedResult  *CheckoutResult
		ExpectedIsError bool
	}

	cases := []TestCase{
		{
			ID: "42",
			ExpectedResult: &CheckoutResult{
				Status:  200,
				Balance: 100500,
				Err:     "",
			},
			ExpectedIsError: false,
		},
		{
			ID: "100500",
			ExpectedResult: &CheckoutResult{
				Status:  400,
				Balance: 0,
				Err:     "bad_balance",
			},
			ExpectedIsError: false,
		},
		{
			ID:              "__broken_json",
			ExpectedResult:  nil,
			ExpectedIsError: true,
		},
		{
			ID:              "__internal_error",
			ExpectedResult:  nil,
			ExpectedIsError: true,
		},
	}

	var PaymentAPIMockHandler = func(w http.ResponseWriter, r *http.Request) {
		// show("Request: ", r)
		switch r.FormValue("id") {
		case "42":
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `{"status": 200, "balance": 100500}`)
		case "100500":
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `{"status": 400, "err": "bad_balance"}`)
		case "__broken_json":
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `{"status": 400`)
		case "__internal_error":
			fallthrough
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	serverMock := httptest.NewServer(http.HandlerFunc(PaymentAPIMockHandler))
	defer serverMock.Close()

	for idx, testCase := range cases {
		cart := &Cart{
			PaymentApiURL: serverMock.URL,
		}
		result, err := cart.CartCheckout(testCase.ID)

		if err != nil && !testCase.ExpectedIsError {
			t.Errorf("[%d] unexpected error: %#v", idx, err)
		}

		if err == nil && testCase.ExpectedIsError {
			t.Errorf("[%d] expected error, got nil", idx)
		}

		if !reflect.DeepEqual(testCase.ExpectedResult, result) {
			t.Errorf("[%d] wrong result, expected %#v, got %#v", idx, testCase.ExpectedResult, result)
		}
	}
}
