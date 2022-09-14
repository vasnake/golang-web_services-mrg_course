package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/shurcooL/httpfs/html/vfstemplate"
)

func NewTemplates(assets http.FileSystem) *template.Template {
	tmpl := template.New("")

	tmpl, err := vfstemplate.ParseGlob(assets, tmpl, "/templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	return tmpl
}
