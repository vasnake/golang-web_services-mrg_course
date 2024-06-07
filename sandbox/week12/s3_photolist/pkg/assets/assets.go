package assets

import (
	"github.com/shurcooL/httpfs/union"
	"net/http"
)

var Assets http.FileSystem = union.New(map[string]http.FileSystem{
	"/templates": http.Dir("./week12/s3_photolist/templates/"),
	"/static":    http.Dir("./week12/s3_photolist/static/"),
})

/*
https://github.com/shurcooL/vfsgen
https://github.com/shurcooL/httpfs/html

go generate --tags=dev
-> go run assets_gen.go assets.go

go build --tags=dev -o ./tmp/dev .
    main.go
    assets.go
    + static/
    + templates/

go build -o ./tmp/release .
    main.go
    assets_vfsdata.go
*/
