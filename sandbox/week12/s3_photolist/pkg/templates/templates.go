package templates

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/shurcooL/httpfs/html/vfstemplate"

	"photolist/pkg/session"
)

type TokenManager interface {
	Create(*session.Session, int64) (string, error)
	Check(*session.Session, string) (bool, error)
}

type MyTemplate struct {
	Tmpl   *template.Template
	Tokens TokenManager
}

func NewTemplates(assets http.FileSystem, tm TokenManager) *MyTemplate {
	tmpl := template.New("")
	tmpl, err := vfstemplate.ParseGlob(assets, tmpl, "/templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
	return &MyTemplate{
		Tmpl:   tmpl,
		Tokens: tm,
	}
}

func (tpl *MyTemplate) Render(ctx context.Context, w http.ResponseWriter, tmplName string, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{}, 3)
	}

	sess, err := session.SessionFromContext(ctx)
	if err == nil {
		data["Authorized"] = true
		data["Session"] = sess

		token, err := tpl.Tokens.Create(sess, time.Now().Add(24*time.Hour).Unix())
		if err != nil {
			log.Println("csrf token creation error:", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		data["CSRFToken"] = token
	} else {
		data["Authorized"] = false
	}

	err = tpl.Tmpl.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		log.Println("cant execute template", err)
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}
