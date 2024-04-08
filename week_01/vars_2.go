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

/*
	// int - платформозависимый тип, 32/64
	var i int = 10

	// автоматически выбранный int
	var autoInt = -10

	// int8, int16, int32, int64
	var bigInt int64 = 1<<32 - 1

	// платформозависимый тип, 32/64
	var unsignedInt uint = 100500

	// uint8, unit16, uint32, unit64
	var unsignedBigInt uint64 = 1<<64 - 1

	fmt.Println(i, autoInt, bigInt, unsignedInt, unsignedBigInt)

	// float32, float64
	var pi float32 = 3.141
	var e = 2.718
	goldenRatio := 1.618

	fmt.Println(pi, e, goldenRatio)

	// bool
	var b bool // false по-умолчанию
	var isOk bool = true
	var success = true
	cond := true

	fmt.Println(b, isOk, success, cond)

	// complex64, complex128
	var c complex128 = -1.1 + 7.12i
	c2 := -1.1 + 7.12i

	fmt.Println(c, c2)
}

*/
