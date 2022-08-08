package main

import "fmt"

// recover from panic using defer

// panic: program stop with stacktrace.
// panic .. recover is not a try .. catch, don't use it this way.
// you should catch panic if your program is a demon of some sort, ant it should continue processing another requests.

func panicButton() {
	// catch panic event on exit and suppress it
	defer func() {
		// process panic event if there was some shit
		if err := recover(); err != nil {
			fmt.Println("panic detected:", err)
		}
	}()

	fmt.Println("working ...")
	panic("some shit happened!")
}

func main() {
	panicButton()
	return
}
