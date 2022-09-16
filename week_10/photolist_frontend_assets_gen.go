// +build ignore
// NB: can't be built unless mentioned explicitly, e.g. `go run assets_gen.go assets.go`

package main

import (
	"github.com/shurcooL/vfsgen"
	"log"
)

// generate resources file `assets_vfsdata.go`
func main() {
	err := vfsgen.Generate(Assets, vfsgen.Options{
		PackageName:  "main",
		BuildTags:    "!dev",
		VariableName: "Assets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
