// +build dev
// NB: added only `--tags=dev` used, templates loaded from disk, not from resources

package main

import (
	"github.com/shurcooL/httpfs/union"
	"net/http"
)

//go:generate go run assets_gen.go assets.go

var Assets http.FileSystem = union.New(map[string]http.FileSystem{
	"/templates": http.Dir("./templates/"),
	"/static":    http.Dir("./static/"),
})

/*
https://github.com/shurcooL/vfsgen
https://github.com/shurcooL/httpfs/html

how to use

go generate --tags=dev
-> go run assets_gen.go assets.go
=> assets_vfsdata.go
    static/
    templates/

go build --tags=dev -o ./tmp/dev .
    main.go
    assets.go

go build -o ./tmp/release .
    main.go
    assets_vfsdata.go

*/
