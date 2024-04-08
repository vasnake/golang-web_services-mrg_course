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
	year          = 2021
	yearTyped int = 2017
)

func main() {
	fmt.Println("const: ", pi, e, hi, zero, three, two, KB, MB)

	// const type inferred in-place
	var bigY int64 = 2021
	var smallY uint16 = 2021
	fmt.Println("untyped const: ", smallY+year, bigY+year) // 4042, 4042
	// fmt.Println("typed const: ", smallY+yearTyped, bigY+yearTyped) // prohibited
	fmt.Println("typed const: ", smallY+uint16(yearTyped), bigY+int64(yearTyped)) // explisit cast

}

/*
const pi = 3.141
const (
	hello = "Привет"
	e     = 2.718
)
const (
	zero = iota
	_    // пустая переменная, пропуск iota
	two
	three // = 3
)
const (
	_         = iota             // пропускаем первое значение
	KB uint64 = 1 << (10 * iota) // 1 << (10 * 1) = 1024
	MB                           // 1 << (10 * 2) = 1048576
)
const (
	// нетипизированная константа
	year = 2017
	// типизированная константа
	yearTyped int = 2017
)

func main() {
	var month int32 = 13
	fmt.Println(month + year)

	// month + yearTyped (mismatched types int32 and int)
	// fmt.Println( month + yearTyped )
}

*/
