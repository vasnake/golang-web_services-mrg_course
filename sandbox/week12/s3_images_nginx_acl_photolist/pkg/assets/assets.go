package assets

import (
	"github.com/shurcooL/httpfs/union"
	"net/http"
)

var Assets http.FileSystem = union.New(map[string]http.FileSystem{
	"/templates": http.Dir("./week12/s3_images_nginx_acl_photolist/templates/"),
	"/static":    http.Dir("./week12/s3_images_nginx_acl_photolist/static/"),
})
