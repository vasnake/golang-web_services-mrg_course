package main

import (
	"fmt"
	"html/template"
	"log"

	"net/http"
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

type Photo struct {
	ID     int
	UserID int
	Path   string
}

var (
	// global list of loaded files
	items  = []*Photo{}
	userID = 0
)

// render list of loaded files
func List(w http.ResponseWriter, r *http.Request) {
	// проблема - каждый раз парсим шаблон
	// проблема - template.Must паникуем при невалидном шаблоне
	tmpl := template.Must(template.New(`list`).Parse(imagesTmpl))

	err := tmpl.Execute(w,
		struct {
			Items []*Photo
		}{
			items,
		})

	if err != nil {
		log.Println("cant execute template", err)
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

// put file to storage, create preview, append to list
func Upload(w http.ResponseWriter, r *http.Request) {
	uploadData, _, err := r.FormFile("my_file")
	if err != nil {
		log.Println("cant parse file", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	defer uploadData.Close()

	md5Sum, err := SaveFile(uploadData)
	if err != nil {
		log.Println("cant save file", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	realFile := "./images/" + md5Sum + ".jpg"
	err = MakeThumbnails(realFile, md5Sum)
	if err != nil {
		log.Println("cant resize file", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	// проблема - глобальная переменная
	items = append(items, &Photo{
		UserID: userID,
		Path:   md5Sum,
	})

	http.Redirect(w, r, "/", 302)
}

// start server
func main() {
	http.HandleFunc("/", List)
	http.HandleFunc("/upload", Upload)

	// should be nginx or other performant file server
	staticHandler := http.StripPrefix(
		"/images/",
		http.FileServer(http.Dir("./images")),
	)
	http.Handle("/images/", staticHandler)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
