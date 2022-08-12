package main

import (
	"fmt"
	"time"
)

func sayHello() {
	fmt.Println("Hello World")
}

func main() {
	timer := time.AfterFunc(1*time.Second, sayHello) // run `sayHello` after 1 sec

	println("Press enter when ready ...")
	fmt.Scanln()
	timer.Stop() // if you press a button before timer is ready, you won't see a hello
}
