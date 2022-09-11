package main

import (
	"fmt"
)

func test_error(is_ok bool) error { // underscore!
	if !is_ok { // underscore!
		fmt.Errorf("failed") // not used (missed return)
	}
	return nil
}

func main() {
	flag := true
	result := test_error(flag)
	fmt.Printf("result is\n", result) // not used
	fmt.Printf("%v is %v\n", flag)    // missed param
}
