//go:build !prod
// +build !prod

//go:generate GOWORK=off go run cmd/assets_gen/assets_gen.go
//# pipeline: go generate; go build --tags=prod
// creates: pkg/assets/assets_vfsdata.go

/*
Как стало:
эмбеддинг файлов в бинарь будет сделан при сборке с тегом "prod".
Иначе ассеты должны быть доступны из традиционной файл.системы.

Так было:
go:build dev
+build dev
go:generate go run assets_gen.go

https://github.com/shurcooL/vfsgen
https://github.com/shurcooL/httpfs/html

go generate --tags=dev
-> go run assets_gen.go assets.go # creates assets_vfsdata.go

зависимости при билде go build --tags=dev ...
    main.go
    assets.go
    static/
    templates/

зависимости при билде go build -o --tags=not_dev ...
    main.go
    assets_vfsdata.go
*/

package assets

import (
	"github.com/shurcooL/httpfs/union"
	"net/http"
)

var Assets http.FileSystem = union.New(map[string]http.FileSystem{
	"/templates": http.Dir("./templates/"),
	"/static":    http.Dir("./static/"),
})
