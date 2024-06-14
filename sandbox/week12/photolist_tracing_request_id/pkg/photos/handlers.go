package photos

import (
	"fmt"
	// "html/template"
	"context"
	// "io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"

	"photolist/pkg/blobstorage"
	"photolist/pkg/session"
	"photolist/pkg/user"
	"photolist/pkg/utils/httputils"
)

type PhotosRepoInterface interface {
	Add(*Photo) (uint32, error)
	GetPhotos(uint32, uint32) ([]*Photo, error)
	Rate(uint32, uint32, int) error
}

// -----------------------------

type Templater interface {
	Render(context.Context, http.ResponseWriter, string, map[string]interface{})
}

type PhotolistHandler struct {
	PhotosRepo  PhotosRepoInterface
	Tmpl        Templater
	UsersRepo   *user.UserRepository
	BlobStorage *blobstorage.S3Storage
}

func (h *PhotolistHandler) ListREST(w http.ResponseWriter, r *http.Request) {
	h.List(w, r, "list.html")
}

func (h *PhotolistHandler) ListGQL(w http.ResponseWriter, r *http.Request) {
	h.List(w, r, "list_gql.html")
}

func (h *PhotolistHandler) List(w http.ResponseWriter, r *http.Request, tmpl string) {
	sess, _ := session.SessionFromContext(r.Context())
	CurrentUser, err := h.UsersRepo.GetByID(sess.UserID)
	if err != nil {
		log.Println("GetUserByID error", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	TargetUser := CurrentUser

	login := strings.Replace(r.URL.Path, "/photos/", "", 1)
	if login != "" {
		TargetUser, err = h.UsersRepo.GetByLogin(login)
		if user.IsErrUserNotFound(err) {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}
		if err != nil {
			log.Println("GetUserByLogin error", err)
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
	}

	vars := map[string]interface{}{
		"CurrentUser": CurrentUser,
		"TargetUser":  TargetUser,
	}
	h.Tmpl.Render(r.Context(), w, tmpl, vars)
}

func (h *PhotolistHandler) UploadAPI(w http.ResponseWriter, r *http.Request) {
	sess, _ := session.SessionFromContext(r.Context())

	r.ParseMultipartForm(5 * 1024 * 1025)
	uploadData, _, err := r.FormFile("my_file")
	if err != nil {
		httputils.RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("cant parse file: %v", err), "internal")
		return
	}
	defer uploadData.Close()

	photoUUID, _ := uuid.NewV4()

	err = h.BlobStorage.Put(uploadData,
		photoUUID.String()+".jpg", "image/jpeg",
		sess.UserID)
	if err != nil {
		httputils.RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("cant save file: %v", err), "internal")
		return
	}

	err = MakeThumbnails(h.BlobStorage,
		uploadData, photoUUID.String(),
		sess.UserID)
	if err != nil {
		httputils.RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("cant save thumbnails: %v", err), "internal")
		return
	}

	comment := r.FormValue("comment")
	_, err = h.PhotosRepo.Add(&Photo{
		UserID:  sess.UserID,
		URL:     photoUUID.String(),
		Comment: comment,
	})
	if err != nil {
		httputils.RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("cant store item: %v", err), "internal")
		return
	}

	httputils.RespJSON(w, map[string]interface{}{
		"status": "ok",
	})
}

func (h *PhotolistHandler) ListAPI(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("uid"))
	if err != nil {
		httputils.RespJSONError(w, http.StatusBadRequest, nil, "bad id")
		return
	}

	sess, _ := session.SessionFromContext(r.Context())
	items, err := h.PhotosRepo.GetPhotos(uint32(id), sess.UserID)
	if err != nil {
		httputils.RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("cate get photos: %v", err), "internal")
		return
	}

	httputils.RespJSON(w, map[string]interface{}{
		"photolist": items,
	})
}

func (h *PhotolistHandler) RateAPI(w http.ResponseWriter, r *http.Request) {
	sess, _ := session.SessionFromContext(r.Context())

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		httputils.RespJSONError(w, http.StatusBadRequest, nil, "bad id")
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
		httputils.RespJSONError(w, http.StatusBadRequest, nil, "bad vote")
		return
	}

	err = h.PhotosRepo.Rate(uint32(id), sess.UserID, rate)
	if err != nil {
		httputils.RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("rate db err: %v", err), "internal")
		return
	}

	httputils.RespJSON(w, map[string]interface{}{
		"id": id,
	})
}
