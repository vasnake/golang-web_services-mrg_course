package main

import "fmt"

func main() {
	show("Sveiki!")
	integers()
}

func integers() {
	show("Int, Decimal: ", 0, 123, 123_456)
}

func show(msg string, xs ...any) {
	var line string = msg
	for _, x := range xs {
		line += fmt.Sprintf("%T(%v); ", x, x)
	}
	fmt.Println(line)
}
