package main

//go:generate hero -source=./template/

import (
	"bytes"
	"fmt"
	"net/http"

	"coursera/template_adv/item"
	"coursera/template_adv/template"
)

// package item
type Item struct {
	Id          int
	Title       string
	Description string
}

var ExampleItems = []*item.Item{
	&item.Item{1, "rvasily", "Mail.ru Group"},
	&item.Item{2, "username", "freelancer"},
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		buffer := new(bytes.Buffer) // should use pool to go easy on GC

		// use generated templates
		template.Index(ExampleItems, buffer)

		w.Write(buffer.Bytes())
	})

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
