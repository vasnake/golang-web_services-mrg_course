/*
embedded assets generator

GOWORK=off go run cmd/assets_gen/assets_gen.go # creates pkg/assets/assets_vfsdata.go
*/

package main

import (
	"fmt"

	"photolist/pkg/assets"

	"github.com/shurcooL/vfsgen"
)

// all paths relative to project root

func main() {
	err := vfsgen.Generate(assets.Assets, vfsgen.Options{
		PackageName:  "assets",
		BuildTags:    "prod",
		VariableName: "Assets",
		Filename:     "pkg/assets/assets_vfsdata.go",
	})
	if err != nil {
		panic(fmt.Errorf("vfsgen.Generate failed: %w", err))
	}
}
