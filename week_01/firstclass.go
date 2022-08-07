package main

import "fmt"

// like any other var
// named function vs anon. function
// assign anon.func to var and call
// define var type as some function(in)out
// using closure as currying substitute

// named function
func doNothing() {
	fmt.Println("I'm a regular function")
}

func main() {

	// anonymous function, creation and call
	func(in string) {
		fmt.Println("anon func, in: ", in)
	}("no name function")

	// assign func to var
	printer := func(in string) {
		fmt.Println("printer input:", in)
	}
	printer("ref to function saved in var")

	// define func type
	type stringPrinterType func(msg string)

	// define function which take a callback and call it
	worker := func(callback stringPrinterType) {
		callback("callback called by worker")
	}
	worker(printer)

	// closure, access to var outside a function body
	prefixer := func(prefix string) stringPrinterType {
		return func(in string) {
			fmt.Printf("[%s] %s\n", prefix, in) // n.b. `prefix` linked to outer var
		}
	}
	successLogger := prefixer("SUCCESS")
	successLogger("should be marked as success")

}
