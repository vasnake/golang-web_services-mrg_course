// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/httpfs/union"
	"github.com/shurcooL/vfsgen"
)

// all paths relative to project root - 102_build ( first introduced )

func main() {
	var Assets http.FileSystem = union.New(map[string]http.FileSystem{
		"/templates": http.Dir("./templates/"),
		"/static":    http.Dir("./static/"),
	})

	err := vfsgen.Generate(Assets, vfsgen.Options{
		PackageName:  "assets",
		BuildTags:    "!dev",
		VariableName: "Assets",
		Filename:     "pkg/assets/assets_vfsdata.go",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
