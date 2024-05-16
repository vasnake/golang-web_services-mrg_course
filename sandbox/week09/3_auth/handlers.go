package auth

import (
	"html/template"
	"log"
	"net/http"
)

// Storage interface for photos
type Storage interface {
	Add(*Photo) error
	GetPhotos(uint32) ([]*Photo, error)
}

// -----------------------------

type PhotolistHandler struct {
	St   Storage
	Tmpl *template.Template
}

// List: show pics (for sessio.user) selected from db
func (h *PhotolistHandler) List(w http.ResponseWriter, r *http.Request) {
	sess, _ := SessionFromContext(r.Context()) // middleware problem
	items, err := h.St.GetPhotos(sess.UserID)
	if err != nil {
		log.Println("cant get items", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	err = h.Tmpl.Execute(w,
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

// Upload: save user (from session) pic to file and db
func (h *PhotolistHandler) Upload(w http.ResponseWriter, r *http.Request) {
	uploadData, _, err := r.FormFile("my_file")
	if err != nil {
		log.Println("cant parse file", err)
		http.Error(w, "request error", http.StatusInternalServerError)
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

	sess, _ := SessionFromContext(r.Context()) // middleware problem
	err = h.St.Add(&Photo{UserID: sess.UserID, Path: md5Sum})
	if err != nil {
		log.Println("cant store item", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/photos", 302)
}
