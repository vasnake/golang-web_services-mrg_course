package main

// codegen example, file with struct for wich code should be generated

import "fmt"

// lets generate code for this struct
// cgen: binpack
type User struct {
	ID       int
	RealName string `cgen:"-"`
	Login    string
	Flags    int
}

// other stuff to make parsing interesting

type Avatar struct {
	ID  int
	Url string
}

var test = 42

func foo() {
	fmt.Printf("Unpacked user %#v", test)
}
