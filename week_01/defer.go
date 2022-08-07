package main

import (
	"fmt"
	"time"
)

func doStuff(in string) string {
	fmt.Printf("doStuff: `%s`\n", in)
	return fmt.Sprintf("stuff `%s` len: %d bytes", in, len(in))
}

func main() {
	// defer call until exit from outer function, similar to `finally`
	defer fmt.Println("Work is done.")
	defer fmt.Println("closing ...")              // n.b. closing called before done
	defer fmt.Println(doStuff("before work ...")) // n.b. parameters are evaluated right here, doStuff called BEFORE, not after!
	defer doStuff("after work ...")

	// anon functions let you to defer params evaluation
	defer func() {
		fmt.Println(doStuff("first deferred call"))
	}()

	fmt.Println("Working ...")
	time.Sleep(3 * time.Second)
	fmt.Println("now deferred ...")

}
