package main

import (
	"fmt"
	"net/http"
)

// start server
func main() {
	h := &PhotolistHandler{
		// ref. to services implementation
		St:   NewStorage(),
		Tmpl: NewTemplates(),
	}

	http.HandleFunc("/", h.List)
	http.HandleFunc("/upload", h.Upload)

	staticHandler := http.StripPrefix(
		"/images/",
		http.FileServer(http.Dir("./images")),
	)
	http.Handle("/images/", staticHandler)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
