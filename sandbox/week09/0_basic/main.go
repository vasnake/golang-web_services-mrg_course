package basic

import (
	"fmt"
	"html/template"
	"log"

	"net/http"
)

var imagesTmpl = `<html>
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
	items  = []*Photo{}
	userID = 0
)

func List(w http.ResponseWriter, r *http.Request) {
	// проблема - каждый раз парсим шаблон
	// проблема - template.Must паникуем при невалидном шаблоне
	tmpl := template.Must(template.New(`list`).Parse(imagesTmpl))

	err := tmpl.Execute(w,
		struct {
			Items []*Photo
		}{
			items, // global var
		})
	if err != nil {
		log.Println("can't execute template", err)
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func Upload(w http.ResponseWriter, r *http.Request) {
	uploadData, _, err := r.FormFile("my_file")
	if err != nil {
		log.Println("can't parse file", err)
		http.Error(w, "Internal error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer uploadData.Close()

	md5Sum, err := SaveFile(uploadData) // magic inside: file name
	if err != nil {
		log.Println("can't save file", err)
		http.Error(w, "Internal error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	realFile := "./images/" + md5Sum + ".jpg" // from magic SaveFile
	err = MakeThumbnails(realFile, md5Sum)
	if err != nil {
		log.Println("can't resize file", err)
		http.Error(w, "Internal error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// проблема - глобальная переменная
	items = append(items, &Photo{
		UserID: userID,
		Path:   md5Sum,
	})

	http.Redirect(w, r, "/", 302)
}

func MainBasic() {
	http.HandleFunc("/", List)
	http.HandleFunc("/upload", Upload)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
