package main

import (
	"fmt"
	"net/http"
)

// сюда код писать не надо

func main() {
	addr := ":8080"
	h := GetApp()
	fmt.Println("start server at", addr)
	http.ListenAndServe(addr, h)
}
