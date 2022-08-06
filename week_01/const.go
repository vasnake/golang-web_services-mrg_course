package main

import "fmt"

const pi = 3.1415 // type inferred, but not here, in application place

// multiple constants
const (
	e  = 2.718
	hi = "Hi"
)

// enum by iota
const (
	zero = iota
	_    // skip `1`
	two
	three
)

// complex enum by iota
const (
	_         = iota             // skip zero
	KB uint64 = 1 << (10 * iota) // 1024
	MB                           // 1048576
)

// untyped const
const (
	year = 2021
)

func main() {
	fmt.Println("const: ", pi, e, hi, zero, three, two, KB, MB)

	// const type inferred in-place
	var bigY int64 = 2021
	var smallY uint16 = 2021
	fmt.Println("untyped const: ", smallY+year, bigY+year) // 4042, 4042

}
