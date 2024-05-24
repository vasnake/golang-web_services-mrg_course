//go:build ignore
// +build ignore

package fronte

import (
	"github.com/shurcooL/vfsgen"
	"log"
)

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
