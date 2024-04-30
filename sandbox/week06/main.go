package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql" // database/sql implementation
	"github.com/gorilla/mux"
)

const (
	port    = 8080
	portStr = ":8080"
	host    = "127.0.0.1"
)

func main() {
	mysqlSimple()
}

func lessonTemplate() {
	show("lessonTemplate: program started ...")
	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err := http.ListenAndServe(portStr, nil)
	show("end of program. ", err)
}

func mysqlSimple() {
	/*
		-- items.sql // sandbox/week06/mysql_items.sql
		SET NAMES utf8;
		SET time_zone = '+00:00';
		SET foreign_key_checks = 0;
		SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

		DROP TABLE IF EXISTS `items`;
		CREATE TABLE `items` (
		  `id` int(11) NOT NULL AUTO_INCREMENT,
		  `title` varchar(255) NOT NULL,
		  `description` text NOT NULL,
		  `updated` varchar(255) DEFAULT NULL,
		  PRIMARY KEY (`id`)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;

		INSERT INTO `items` (`id`, `title`, `description`, `updated`) VALUES
		(1,	'database/sql',	'Рассказать про базы данных',	'foo'),
		(2,	'memcache',	'Рассказать про мемкеш с примером использования',	NULL);
	*/

	show("mysqlItems: program started ...")

	// основные настройки к базе
	dsn := "root@tcp(localhost:3306)/coursera?"
	// указываем кодировку
	dsn += "&charset=utf8"
	// отказываемся от prapared statements // параметры подставляются сразу
	dsn += "&interpolateParams=true"
	db, err := sql.Open("mysql", dsn)
	db.SetMaxOpenConns(10)
	err = db.Ping() // вот тут будет первое подключение к базе
	__err_panic(err)
	show("connected to DB ", db)

	tmpl := template.Must(template.ParseGlob("./week06/crud_templates/*")) // sandbox\week06\crud_templates\
	show("loaded templates ", tmpl)

	srv := &MysqlSimpleHttpHandlers{
		DB:   db,
		Tmpl: tmpl,
	}

	// в целях упрощения примера пропущена авторизация и csrf
	r := mux.NewRouter() // gorilla/mux
	r.HandleFunc("/", srv.ListAll).Methods("GET")
	r.HandleFunc("/items", srv.ListAll).Methods("GET")
	r.HandleFunc("/items/new", srv.ShowCreateForm).Methods("GET")
	r.HandleFunc("/items/new", srv.CreateItem).Methods("POST")
	r.HandleFunc("/items/{id}", srv.ShowUpdateForm).Methods("GET")
	r.HandleFunc("/items/{id}", srv.UpdateItem).Methods("POST")
	r.HandleFunc("/items/{id}", srv.DeleteItem).Methods("DELETE")

	show("Starting server at: ", host+portStr)
	show(fmt.Sprintf("Open url http://localhost%s/", portStr))
	err = http.ListenAndServe(portStr, r)
	show("end of program. ", err)
}

// sql table repr.
type MysqlSimplePostItem struct {
	Id          int
	Title       string
	Description string
	Updated     sql.NullString
}

type MysqlSimpleHttpHandlers struct {
	DB   *sql.DB
	Tmpl *template.Template
}

func (h *MysqlSimpleHttpHandlers) ListAll(w http.ResponseWriter, r *http.Request) {
	items := []*MysqlSimplePostItem{}

	rows, err := h.DB.Query("SELECT id, title, updated FROM items")
	__err_panic(err)
	for rows.Next() {
		post := &MysqlSimplePostItem{}
		err = rows.Scan(&post.Id, &post.Title, &post.Updated)
		__err_panic(err)
		items = append(items, post)
	}
	// надо закрывать соединение, иначе будет течь
	rows.Close()

	err = h.Tmpl.ExecuteTemplate(w, "index.html", struct {
		Items []*MysqlSimplePostItem
	}{
		Items: items,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MysqlSimpleHttpHandlers) ShowCreateForm(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "create.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MysqlSimpleHttpHandlers) CreateItem(w http.ResponseWriter, r *http.Request) {
	// в целях упрощения примера пропущена валидация
	result, err := h.DB.Exec(
		"INSERT INTO items (`title`, `description`) VALUES (?, ?)",
		r.FormValue("title"),
		r.FormValue("description"),
	)
	__err_panic(err)

	affected, err := result.RowsAffected()
	__err_panic(err)
	lastID, err := result.LastInsertId()
	__err_panic(err)

	fmt.Println("Insert: RowsAffected ", affected, "; LastInsertId ", lastID)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *MysqlSimpleHttpHandlers) ShowUpdateForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	__err_panic(err)

	post := &MysqlSimplePostItem{}
	// QueryRow сам закрывает коннект
	row := h.DB.QueryRow("SELECT id, title, updated, description FROM items WHERE id = ?", id)

	err = row.Scan(&post.Id, &post.Title, &post.Updated, &post.Description)
	__err_panic(err)

	err = h.Tmpl.ExecuteTemplate(w, "edit.html", post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MysqlSimpleHttpHandlers) UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	__err_panic(err)

	// в целях упрощения примера пропущена валидация
	result, err := h.DB.Exec(
		"UPDATE items SET"+
			"`title` = ?"+
			",`description` = ?"+
			",`updated` = ?"+
			"WHERE id = ?",
		r.FormValue("title"),
		r.FormValue("description"),
		"foo",
		id,
	)
	__err_panic(err)

	affected, err := result.RowsAffected()
	__err_panic(err)

	fmt.Println("Update: RowsAffected", affected, "; item id: ", id)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *MysqlSimpleHttpHandlers) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	__err_panic(err)

	result, err := h.DB.Exec(
		"DELETE FROM items WHERE id = ?",
		id,
	)
	__err_panic(err)

	affected, err := result.RowsAffected()
	__err_panic(err)

	fmt.Println("Delete: RowsAffected", affected)

	w.Header().Set("Content-type", "application/json")
	resp := `{"affected": ` + strconv.Itoa(int(affected)) + `}`
	w.Write([]byte(resp))
}

// не используйте такой код в prod // ошибка должна всегда явно обрабатываться
func __err_panic(err error) {
	if err != nil {
		panic(err)
	}
}

// ts returns current timestamp in RFC3339 with milliseconds
func ts() string {
	/*
		https://pkg.go.dev/time#pkg-constants
		https://stackoverflow.com/questions/35479041/how-to-convert-iso-8601-time-in-golang
	*/
	const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
	return time.Now().UTC().Format(RFC3339Milli)
}

// show writes message to standard output. Message combined from prefix msg and slice of arbitrary arguments
func show(msg string, xs ...any) {
	var line = ts() + ": " + msg

	for _, x := range xs {
		// https://pkg.go.dev/fmt
		line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		// line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
}
