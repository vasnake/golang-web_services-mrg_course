package main

import (
	"fmt"
	"go/types"
	"reflect"
)

func main() {
	var num0 int // platform dependent int, 32 or 64 or ? with default value
	fmt.Println("default int: ", num0)

	var num1 int = 1 // set init value
	fmt.Println("with init value: ", num1)

	var num2 = 2 // type inference
	fmt.Println("detected type: ", num2, types.BasicInfo(num2), reflect.TypeOf(num2).String())
	fmt.Printf("type: %T\n", num2)

	num := 3 // skip `var` and `type`
	fmt.Println("new var: ", num, reflect.TypeOf(num).String())

	num += 1 // increment
	num++    // postfix increment
	fmt.Println("incremented: ", num)

	var a, b = 11, 12 // define multiple vars
	a, c := 13, 14    // only if `c` is new
	a, b, c = 1, 2, 3 // multi assign
	fmt.Println("multi: ", a, b, c)
}
