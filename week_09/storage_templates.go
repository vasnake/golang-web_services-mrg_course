package main

import (
	"html/template"
)

var imagesTmpl = `
<html>
	<body>
	<div>
		<form action="/upload" method="post" enctype="multipart/form-data">
			Image: <input type="file" name="my_file">
			<input type="submit" value="Upload">
		</form>
	</div>
	<br />
	{{range .Items}}
		<div>
			<img src="/images/{{.Path}}_160.jpg" />
			<br />
		</div>
	{{end}}
	</body>
</html>
`

// compile one time, on app start
func NewTemplates() *template.Template {
	tmpl := template.Must(template.New(`example`).Parse(imagesTmpl))
	return tmpl
}
