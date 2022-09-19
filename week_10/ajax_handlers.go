package main

import (
	"fmt"
	// "html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	// "time"
	"database/sql"
)

type Storage interface {
	Add(*Photo) error
	GetPhotos(uint32, uint32) ([]*Photo, error)
	Rate(uint32, uint32, int) error
}

type PhotolistHandler struct {
	St     Storage
	Tmpl   *MyTemplate
	Tokens TokenManager
	UserDB *sql.DB
}

func (h *PhotolistHandler) List(w http.ResponseWriter, r *http.Request) {
	sess, _ := SessionFromContext(r.Context())
	CurrentUser, err := GetUserByID(h.UserDB, sess.UserID)
	if err != nil {
		log.Println("GetUserByID error", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// user wall by default
	TargetUser := CurrentUser

	// wall for other user?
	login := strings.Replace(r.URL.Path, "/photos/", "", 1)
	if login != "" {
		TargetUser, err = GetUserByLogin(h.UserDB, login)
		if err == errUserNotFound {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}
		if err != nil {
			log.Println("GetUserByLogin error", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
	}

	// setup data for template
	vars := map[string]interface{}{
		"CurrentUser": CurrentUser,
		"TargetUser":  TargetUser,
	}
	h.Tmpl.Render(r.Context(), w, "list.html", vars)
}

// save file to storage, add record to DB
func (h *PhotolistHandler) UploadAPI(w http.ResponseWriter, r *http.Request) {
	sess, _ := SessionFromContext(r.Context())

	r.ParseMultipartForm(5 * 1024 * 1025)

	uploadData, _, err := r.FormFile("my_file")
	if err != nil {
		RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("cant parse file: %v", err), "internal")
		return
	}
	defer uploadData.Close()

	comment := r.FormValue("comment")

	md5Sum, err := SaveFile(uploadData)
	if err != nil {
		RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("cant save file: %v", err), "internal")
		return
	}

	realFile := "./images/" + md5Sum + ".jpg"
	err = MakeThumbnails(realFile, md5Sum)
	if err != nil {
		RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("cant resize photo: %v", err), "internal")
		return
	}

	err = h.St.Add(&Photo{
		UserID:  sess.UserID,
		Path:    md5Sum,
		Comment: comment,
	})
	if err != nil {
		RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("cant store item: %v", err), "internal")
		return
	}

	RespJSON(w, map[string]interface{}{
		"status": "ok",
	})
}

// list files for user
func (h *PhotolistHandler) ListAPI(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("uid"))
	if err != nil {
		RespJSONError(w, http.StatusBadRequest, nil, "bad id")
		return
	}

	sess, _ := SessionFromContext(r.Context())

	items, err := h.St.GetPhotos(uint32(id), sess.UserID)
	if err != nil {
		RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("cate get photos: %v", err), "internal")
		return
	}

	RespJSON(w, map[string]interface{}{
		"photolist": items,
	})
}

// update file rating
func (h *PhotolistHandler) RateAPI(w http.ResponseWriter, r *http.Request) {
	sess, _ := SessionFromContext(r.Context())

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		RespJSONError(w, http.StatusBadRequest, nil, "bad id")
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
		RespJSONError(w, http.StatusBadRequest, nil, "bad vote")
		return
	}

	err = h.St.Rate(uint32(id), sess.UserID, rate)
	if err != nil {
		RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("rate db err: %v", err), "internal")
		return
	}

	RespJSON(w, map[string]interface{}{
		"id": id,
	})
}
