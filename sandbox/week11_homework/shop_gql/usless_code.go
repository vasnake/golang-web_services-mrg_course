package main

import (
	"fmt"
	"io"
	"net/http"
)

type ShopHttpApp struct {
	data any
}

var _ http.Handler = NewShopHttpApp() // type check

func NewShopHttpApp() *ShopHttpApp {
	return &ShopHttpApp{
		data: nil,
	}
}

// ServeHTTP implements http.Handler.
func (s *ShopHttpApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	show("ShopHttpApp.ServeHTTP, request: ", r.Method, r.URL.Path, r.Form)
	bodyBytes, err := io.ReadAll(r.Body)
	panicOnError("io.ReadAll failed", err)
	show(fmt.Sprintf("request body: %s", string(bodyBytes)))
	http.Error(w, "oops", http.StatusNotImplemented)
}
