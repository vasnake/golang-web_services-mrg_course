package main

import (
	"./person"
	"fmt"
)

func main() {
	p := person.NewPerson(1, "Foo", "bar")
	fmt.Printf("main, person %+v\n", p)

	// p.secret is private, only Capital letters are exported
	fmt.Println("main, secret", person.GetSecret(p))
}
