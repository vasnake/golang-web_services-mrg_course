package main

import (
	"fmt"
	"net/http"
)

func runServer(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Addr:", addr, "URL:", r.URL.String())
		})

	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Println("starting server at", addr)
	server.ListenAndServe()
}

func main() {
	// two async servers on two ports
	go runServer(":8081")
	runServer(":8080")
}
