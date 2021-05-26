package main

import "fmt"

// slice = buffer, (length, capacity)
func main() {
	// creation, make function, default values
	var buf0 []int                  // like array, but w/o size, len:0, cap:0
	buf1 := []int{}                 // empty enumeration using literal
	buf2 := []int{42}               // one elem, len:1, cap:1
	buf3 := make([]int, 0)          // len:0, cap:0
	buf4 := make([]int, 5)          // len:5, cap:5
	buf5 := make([]int, 5, 10)      // len:5, cap:10
	fmt.Printf("buf0 %#v \n", buf0) //buf0 []int(nil)
	fmt.Printf("buf1 %#v \n", buf1) //buf1 []int{}
	fmt.Printf("buf2 %#v \n", buf2) //buf2 []int{42}
	fmt.Printf("buf3 %#v \n", buf3) //buf3 []int{}
	fmt.Printf("buf4 %#v \n", buf4) //buf4 []int{0, 0, 0, 0, 0}
	fmt.Printf("buf5 %#v \n", buf5) //buf5 []int{0, 0, 0, 0, 0}

	someInt := buf2[0]
	buf2[0] = 37
	fmt.Printf("buf2 %#v, oldval %#v \n", buf2, someInt) // buf2 []int{37}, oldval 42

	// append elements
	var buf []int            // len:0, cap:0
	buf = append(buf, 9, 10) // len:2, cap:2
	buf = append(buf, 12)    // len:2, cap:4, capacity x2 every reallocation
	// old storage deleted if unused elsewhere, references broken
	fmt.Printf("buf %#v \n", buf) // buf []int{9, 10, 12}

	// append elements from another slice
	otherBuf := make([]int, 3)     // [0, 0, 0]
	buf = append(buf, otherBuf...) // unpack `from` using `...`
	fmt.Printf("buf %#v \n", buf)  // buf []int{9, 10, 12, 0, 0, 0}

	// get len, cap from slice
	fmt.Printf("buf len: %#v, cap: %#v \n", len(buf), cap(buf)) // buf len: 6, cap: 8
}
