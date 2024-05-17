package jwt_session

import (
	"html/template"
)

var imagesTmpl = `<html>
	<head>
	<script>
		function rateComment(id, vote) {
			var request = new XMLHttpRequest();
			request.open('POST', '/photos/rate?id='+id+"&vote="+(vote > 0 ? "up" : "down"), true);
			request.setRequestHeader("csrf-token", "{{.CSRFToken}}");
			request.onload = function() {
				var resp = JSON.parse(request.responseText);
				if(resp.err) {
					console.log("rateComment server err:", resp.err);
					return;
				}
				var elem = document.querySelector('#rating-'+resp.id);
				rating = parseInt(elem.innerHTML) + parseInt(vote);
				elem.innerHTML = rating;
			};
			request.send();
		}
	</script>
	</head>
	<body>
	<div style="text-align:right;">
		<a href="/user/change_pass">Change password</a> | 
		<a href="/user/logout">Logout</a>
	</div>
	<!-- div>
		&lt;script&gt;alert(document.cookie)&lt;/script&gt;
		<br />
		&lt;img src=&quot;/photos/rate?id=1&amp;vote=up&quot;&gt;
	</div -->
	<br />
	<div>
		<form action="/photos/upload" method="post" enctype="multipart/form-data">
			<input type="hidden" value="{{.CSRFToken}}" name="csrf-token" />
			Image: <input type="file" name="my_file"><br />
			<textarea rows=4 cols=50 name="comment"></textarea></br>
			<input type="submit" value="Upload">
		</form>
	</div>
	<br />
	{{range .Items}}
		<div>
			<img src="/images/{{.Path}}_160.jpg" />
			<br />
			<div style="border: 1px solid black; padding: 5px; margin: 5px;">
				<button onclick="rateComment({{.ID}}, 1)">&uarr;</button>
				<span id="rating-{{.ID}}">{{.Rating}}</span>
				<button onclick="rateComment({{.ID}}, -1)">&darr;</button>
				&nbsp;

				<!-- text/template по-умолчанию ничего не экранируется --> 
				<!-- html/template по-умолчанию будет экранировать --> 
				{{.Comment}}
			</div>
		</div>
	{{end}}
	</body>
</html>
`

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

var changePassTmpl = `<html>
<body>
<div>
	<form action="/user/change_pass" method="post" autocomplete="off">
		<input type="password" name="old_password" placeholder="Current password"><br />
		<input type="password" name="pass1" placeholder="New Password"><br />
		<input type="password" name="pass2" placeholder="Repeat new Password"><br />
		<input type="submit" value="Change password">
	</form>
</div>
</body>
</html>
`

func NewTemplates() *template.Template {
	tmpl := template.New("blank")

	template.Must(tmpl.New("list").Parse(imagesTmpl))
	template.Must(tmpl.New("login").Parse(loginTmpl))
	template.Must(tmpl.New("reg").Parse(regTmpl))
	template.Must(tmpl.New("change_password").Parse(changePassTmpl))

	return tmpl
}
