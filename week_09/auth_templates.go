package main

import (
	"html/template"
)

var imagesTmpl = `
<html>
	<body>
	<div>
		<form action="/photos/upload" method="post" enctype="multipart/form-data">
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

var loginTmpl = `<html>
<body>
<div>
	<form action="/user/login" method="post" autocomplete="off">
		<input type="text" name="login" placeholder="Login"><br />
		<input type="password" name="password" placeholder="Password"><br />
		<input type="submit" value="Login"> <a href="/user/reg">Registration</a>
	</form>
</div>
</body>
</html>
`

var regTmpl = `<html>
<body>
<div>
	<form action="/user/reg" method="post" autocomplete="off">
		<input type="text" name="login" placeholder="Login"><br />
		<input type="password" name="password" placeholder="Password"><br />
		<input type="submit" value="Registration">
	</form>
</div>
</body>
</html>
`

func NewUserTemplates() *template.Template {
	tmpl := template.New("blank")

	template.Must(tmpl.New("login").Parse(loginTmpl))
	template.Must(tmpl.New("reg").Parse(regTmpl))

	return tmpl
}
