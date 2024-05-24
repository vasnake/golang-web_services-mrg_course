package fronte

import (
	"html/template"
	"log"
	"net/http"

	"github.com/shurcooL/httpfs/html/vfstemplate"
)

func NewTemplates(assets http.FileSystem) *template.Template {
	tmpl := template.New("")
	tmpl, err := vfstemplate.ParseGlob(assets, tmpl, "/templates/*.html") // filename will be used as ref.
	if err != nil {
		log.Fatal(err)
	}
	return tmpl
}
