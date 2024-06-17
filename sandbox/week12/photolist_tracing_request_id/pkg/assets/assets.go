package assets

import (
	"github.com/shurcooL/httpfs/union"
	"net/http"
)

var Assets http.FileSystem = union.New(map[string]http.FileSystem{
	"/templates": http.Dir("./week12/photolist_tracing_request_id/templates/"),
	"/static":    http.Dir("./week12/photolist_tracing_request_id/static/"),
})
