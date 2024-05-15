package storage

import (
	"html/template"
	"log"
	"net/http"
)

type PhotolistHandler struct {
	St   *StMem
	Tmpl *template.Template
}

// http handler
func (h *PhotolistHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.St.GetPhotos(userID) // using interface, good
	if err != nil {
		log.Println("can't get items", err)
		http.Error(w, "storage error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.Tmpl.Execute(w,
		struct {
			Items []*Photo
		}{
			items,
		}) // not so good, but ok for now
	if err != nil {
		log.Println("can't execute template", err)
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// http handler
func (h *PhotolistHandler) Upload(w http.ResponseWriter, r *http.Request) {
	uploadData, _, err := r.FormFile("my_file")
	if err != nil {
		log.Println("can't parse file", err)
		http.Error(w, "request error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer uploadData.Close()

	md5Sum, err := SaveFile(uploadData) // still magic
	if err != nil {
		log.Println("can't save file", err)
		http.Error(w, "Internal error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	realFile := "./images/" + md5Sum + ".jpg"
	err = MakeThumbnails(realFile, md5Sum)
	if err != nil {
		log.Println("can't resize file", err)
		http.Error(w, "Internal error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.St.Add(&Photo{UserID: userID, Path: md5Sum})
	if err != nil {
		log.Println("can't store item", err)
		http.Error(w, "storage error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/photos", 302) // wrong, we have not /photos route
}

// global var, problem
var (
	userID = 0
)
