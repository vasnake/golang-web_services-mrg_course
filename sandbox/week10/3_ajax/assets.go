package ajax3

import (
	"github.com/shurcooL/httpfs/union"
	"net/http"
)

var Assets http.FileSystem = union.New(map[string]http.FileSystem{
	"/templates": http.Dir("./week10/3_ajax/templates/"),
	"/static":    http.Dir("./week10/3_ajax/static/"),
})
