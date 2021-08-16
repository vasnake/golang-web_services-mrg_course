package main

import (
	"./person"
	"fmt"
)

func main() {
	// package code contains in a few files

	p := person.NewPerson(1, "Foo", "bar")
	fmt.Printf("main, person %+v\n", p)
	//main, person &{ID:1 Name:Foo secret:bar}

	// p.secret is private, only Capital letters are exported
	fmt.Println("main, secret", person.GetSecret(p))
	// main, secret bar
}

// GO111MODULE=off go run visibility/main.go
