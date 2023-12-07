package main

import (
	"fmt"
	"net/http"
)

// servehttp demo

// Customizable handler, define it's behaviour with parameter. Implements net/http/Handler
//
// type Handler interface { ServeHTTP(ResponseWriter, *Request) }
type Handler_servehttp struct {
	Name string
}

func (h *Handler_servehttp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Name:", h.Name, "URL:", r.URL.String())
}
