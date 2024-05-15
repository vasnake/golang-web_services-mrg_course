package sql_storage

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func MainSqlStorage() {
	// основные настройки к базе
	// dsn := "root:love@tcp(127.0.0.1:3306)/photolist?charset=utf8&interpolateParams=true"
	dsn := "root@tcp(127.0.0.1:3306)/photolist?charset=utf8&interpolateParams=true"
	db, err := sql.Open("mysql", dsn)

	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalf("can't connect to db, err: %v\n", err)
	}
	defer db.Close()

	h := &PhotolistHandler{
		St:   NewDbStorage(db), // no global vars
		Tmpl: NewTemplates(),
	}

	// same shit, no changes
	http.HandleFunc("/", h.List)
	http.HandleFunc("/upload", h.Upload)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
