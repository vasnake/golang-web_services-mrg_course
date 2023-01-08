// export app_dir=~/go_sandbox/hello_world mkdir -p $app_dir && pushd $app_dir
// touch main.go
// go mod init hello_world
// pushd ..
// go work init
// go work use ./hello_world/
// popd
// edit main.go
// gofmt -w main.go
// go run main.go

package main

import "fmt"

func main() {
	name := "Gofer"
	fmt.Println("Hello World! Yours truly,", name, 0_123 == 0o123, 072.40, 0x1p-2)
	// Hello World! Yours truly, Gofer true 72.4 0.25
}
