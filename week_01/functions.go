package main

import "fmt"

// simple definition
// multi-parameter definition
// named return var
// return tuple
// return named vars as tuple // beware default values for uninitialized vars
// unpacked list of parameters, triple-dot notation

func sqrt(in int) int {
	return in * in
}

func sum3(a, b int, c int) int { // n.b. a and b are the same type
	return a + b + c
}

func namedReturn() (out int) {
	out = 42
	return
	// return 37 // also will work
}

func withError(in int) (int, error) {
	if in > 1 {
		return 1, fmt.Errorf("shit happens")
	}
	return 0, nil
}

func withErrorNamed(cond bool) (res int, err error) {
	if cond {
		err = fmt.Errorf("shit happens")
		return // n.b. res will be 0 by default
		// return 0, fmt.Errorf("???")
	}
	return 42, nil
}

func sum(in ...int) (res int) {
	fmt.Printf("in: %#v \n", in) // slice in: []int{1, 2, 3, 4}
	for _, v := range in {
		res += v
	}
	return
}

func main() {
	println("10: ", sum(1, 2, 3, 4))
}
