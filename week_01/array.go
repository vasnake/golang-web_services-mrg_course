package main

import "fmt"

func main() {
	// array size defines array type, you can't use [2]int where [3]int is required

	var a1 [3]int
	fmt.Printf("a1 short %v, a1 full %#v \n", a1, a1) // a1 short [0 0 0], a1 full [3]int{0, 0, 0}

	// const can be used
	const size = 2
	var a2 [2 * size]bool
	fmt.Printf("a2 %#v \n", a2) // a2 [4]bool{false, false, false, false}

	// detect size from values enum
	a3 := [...]int{1, 2, 3}
	fmt.Printf("a3 %#v \n", a3) // a3 [3]int{1, 2, 3}

	// compile time panic
	//a3[size + 1] = 42

}
