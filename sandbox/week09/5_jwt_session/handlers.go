package jwt_session

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

// updated shit

type Storage interface {
	Add(*Photo) error
	GetPhotos(uint32) ([]*Photo, error)
	Rate(uint32, int) error
}

type TokenManager interface {
	Create(*Session, int64) (string, error)
	Check(*Session, string) (bool, error)
}

// -----------------------------

// service http handlers (list, upload, rate-pic)
type PhotolistHandler struct {
	St     Storage // storage interface
	Tmpl   *template.Template
	Tokens TokenManager // token interface (create, check), csrf
}

// http handler: get session (from auth-middleware), get pics from storage, create new csrf token;
// render template (token, pics)
func (svcH *PhotolistHandler) List(w http.ResponseWriter, r *http.Request) {
	sess, _ := SessionFromContext(r.Context())
	items, err := svcH.St.GetPhotos(sess.UserID)
	if err != nil {
		log.Println("cant get items", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	token, err := svcH.Tokens.Create(sess, time.Now().Add(24*time.Hour).Unix())
	if err != nil {
		log.Println("csrf token creation error:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	err = svcH.Tmpl.ExecuteTemplate(w, "list",
		struct {
			Items     []*Photo
			CSRFToken string
		}{
			Items:     items,
			CSRFToken: token,
		})
	if err != nil {
		log.Println("cant execute template", err)
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}

// http handler: get session from auth-middleware, get csrf token from form, check csrf token, add pic. file to storage
func (svcH *PhotolistHandler) Upload(w http.ResponseWriter, r *http.Request) {
	sess, _ := SessionFromContext(r.Context())
	CSRFToken := r.FormValue("csrf-token")
	ok, err := svcH.Tokens.Check(sess, CSRFToken)
	if !ok || err != nil {
		log.Println("csrf token check fail:", ok, err)
		http.Error(w, "bad token", http.StatusUnauthorized)
		return
	}

	uploadData, _, err := r.FormFile("my_file")
	if err != nil {
		log.Println("cant parse file", err)
		http.Error(w, "request error", http.StatusInternalServerError)
		return
	}
	defer uploadData.Close()

	comment := r.FormValue("comment")

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

	err = svcH.St.Add(&Photo{
		UserID:  sess.UserID,
		Path:    md5Sum,
		Comment: comment,
	})
	if err != nil {
		log.Println("cant store item", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/photos", 302)
}

// http handler: get session from middleware, get csrf token from http header, update rating in storage
func (svcH *PhotolistHandler) Rate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	sess, _ := SessionFromContext(r.Context())
	CSRFToken := r.Header.Get("csrf-token")
	ok, err := svcH.Tokens.Check(sess, CSRFToken)
	if !ok || err != nil {
		log.Println("csrf token check fail:", ok, err)
		http.Error(w, `{"err": "bad token"}`, http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, `{"err": "bad id"}`, http.StatusBadRequest)
		return
	}
	vote := r.FormValue("vote")
	rate := 0
	switch vote {
	case "up":
		rate = 1
	case "down":
		rate = -1
	default:
		http.Error(w, `{"err": "bad vote"}`, http.StatusBadRequest)
		return
	}

	err = svcH.St.Rate(uint32(id), rate)
	if err != nil {
		log.Println("rate err: ", err)
		http.Error(w, `{"err": "db err"}`, http.StatusBadRequest)
		return
	}

	result := map[string]interface{}{
		"id": id,
	}
	resp, _ := json.Marshal(result)
	w.Write(resp)
}
