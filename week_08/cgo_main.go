package main

/*
go build

https://golang.org/cmd/cgo/#hdr-C_references_to_Go
см. послений абзац

main.cgo2.o: In function `Multiply':
... multiple definition of `Multiply'
... main.go:10: first defined here
collect2.exe: error: ld returned 1 exit status
*/

/*
void Multiply(int a, int b); // declaration, implementation in other file to avoid errors like mentioned above
*/
import "C" //это псевдо-пакет, он реализуется компилятором
import "fmt"

//export printResultGolang
func printResultGolang(result C.int) {
	// called from C code
	fmt.Printf("result-var internals %T = %+v\n", result, result)
}

/*
	переходы между рантаймами:
	go - main
	cgo - Multiply
	go - printResultGolang
	cgo - Multiply
	go - main
*/

func main() {
	a := 2
	b := 3

	C.Multiply(C.int(a), C.int(b)) // call C code
}
