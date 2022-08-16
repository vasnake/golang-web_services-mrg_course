package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Main page")
}

func main() {

	// N.B. route w.o. ending slash
	http.HandleFunc("/page",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Single page:", r.URL.String())
		})

	// N.B. route with ending slash, it's a prefix for a class of supported routes
	http.HandleFunc("/pages/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Multiple pages:", r.URL.String())
		})

	http.HandleFunc("/", handler)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
