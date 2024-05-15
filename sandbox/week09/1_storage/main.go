package storage

import (
	"fmt"
	"net/http"
)

func MainStorage() {
	// w/o global vars
	h := &PhotolistHandler{
		St:   NewStorage(),
		Tmpl: NewTemplates(),
	}

	http.HandleFunc("/", h.List)
	http.HandleFunc("/upload", h.Upload)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
