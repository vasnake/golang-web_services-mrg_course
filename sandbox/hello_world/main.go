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
// GOOS=linux GOARCH=amd64 go build -o hello_world-amd64-linux main.go
// hello_world>set GOOS=linux&&set GOARCH=amd64&&go build -o hello_world-amd64-linux main.go
// mv -v ~/data/github/golang-web_services-mrg_course/sandbox/hello_world/hello_world-amd64-linux ~/data/gitlab/dm.dmgrinder/dmgrinder/

/*
# /workdir/project_files/hello_world-amd64-linux > /dev/null
ERR: Answer is, int(42)
# /workdir/project_files/hello_world-amd64-linux 2> /dev/null
OUT: Hello World! Yours truly, string(Gofer)
*/

package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	name := "Gofer"
	answer := 42
	StdOut("Hello World! Yours truly", name)
	StdErr("Answer is", answer)
	// fmt.Fprintf(os.Stderr, "number of foo: %d", nFoo)
}

// StdOut print line to stdout
func StdOut(msg string, xs ...any) {
	var line = MakeLine("OUT: "+msg, xs)
	PrintLine(line, os.Stdout)
}

// StdErr print line to stderr
func StdErr(msg string, xs ...any) {
	var line = MakeLine("ERR: "+msg, xs)
	PrintLine(line, os.Stderr)
}

// MakeLine adds to msg all xs values
func MakeLine(msg string, xs []any) string {
	var line string = msg
	for _, x := range xs {
		line += fmt.Sprintf(", %T(%v)", x, x) // %#v
	}
	return line
}

// PrintLine print line msg to writer w
func PrintLine(msg string, w io.Writer) {
	fmt.Fprintln(w, msg)
}
