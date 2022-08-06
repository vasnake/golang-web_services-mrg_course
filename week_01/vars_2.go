package main

import "fmt"

func main() {
	var platformInt int = 1          // platform dependent
	var autoInt = 2                  // compiler decision, should be platform int
	var bigInt int64 = (1 << 63) - 1 // you can select from int8 .. int64
	var uInt uint = 1 << 63          // unsigned, same as int
	fmt.Println("ints: ", platformInt, autoInt, bigInt, uInt)

	var pi float32 = 3.14 // explicit 32 or 64
	var e float64 = 2.718
	var _pi = 3.14          // platform dependent float
	var flotDefault float32 // zero default
	fmt.Println("floats: ", pi, e, _pi, flotDefault)

	var defaultBool bool // default 0 == false
	var b bool = true
	condition := 0 == 0
	fmt.Println("bools: ", defaultBool, b, condition)

	// complex
	var c complex64 = -1.1 + 7.12i
	complex := 1.2 + 3.4i // complex128 by default
	fmt.Println("complex: ", c, complex)
}
