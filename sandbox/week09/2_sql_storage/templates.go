package sql_storage

import (
	"html/template"
)

// same shit

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

func NewTemplates() *template.Template {
	tmpl := template.Must(template.New(`example`).Parse(imagesTmpl))
	return tmpl
}
