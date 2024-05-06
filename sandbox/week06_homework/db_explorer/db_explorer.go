package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

/*
* GET / - возвращает список все таблиц (которые мы можем использовать в дальнейших запросах)
* GET /$table?limit=5&offset=7 - возвращает список из 5 записей (limit) начиная с 7-й (offset) из таблицы $table.
	limit по-умолчанию 5, offset 0
* GET /$table/$id - возвращает информацию о самой записи или 404
* PUT /$table - создаёт новую запись, данный по записи в теле запроса (POST-параметры)
* POST /$table/$id - обновляет запись, данные приходят в теле запроса (POST-параметры)
* DELETE /$table/$id - удаляет запись


*/

// NewDbExplorer create http handler for db_explorer app
func NewDbExplorer(dbRef *sql.DB) (http.Handler, error) {
	srv := &MysqlExplorerHttpHandlers{
		DB: dbRef,
	}
	r := mux.NewRouter() // gorilla/mux

	r.HandleFunc("/", srv.ListTables).Methods("GET")
	r.HandleFunc("/{table}", srv.ReadTable).Methods("GET")
	r.HandleFunc("/{table}/{id}", srv.ReadRecord).Methods("GET")
	r.HandleFunc("/{table}/", srv.CreateRecord).Methods("PUT")
	r.HandleFunc("/{table}/{id}", srv.UpdateRecord).Methods("POST")
	r.HandleFunc("/{table}/{id}", srv.DeleteRecord).Methods("DELETE")

	return r, nil
}

type MysqlExplorerHttpHandlers struct {
	DB *sql.DB
}

func (srv *MysqlExplorerHttpHandlers) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	/*
	   * DELETE /$table/$id - удаляет запись

	   2024-05-06T07:35:58.612Z: UpdateRecord, table name, recId: "items"; "3";

	   	main_test.go:546: [case 19: [DELETE] /items/3 ] expected http status 200, got 405

	   	Case{
	   		Path:   "/items/3",
	   		Method: http.MethodDelete,
	   		ExpectedRespBody: GenericMap{
	   			"response": GenericMap{
	   				"deleted": 1,
	   			},
	   		},
	   	}, // 19
	*/

	defer recoverPanic(w)

	exitOnError := func(err error, status int, msg string) bool {
		if err != nil {
			show("DeleteRecord, error: ", err)
			writeError(status, msg, w)
			return true
		}
		return false
	}

	// params
	routeVarsMap := MapSS(mux.Vars(r))
	tableName := routeVarsMap.getOrDefault("table", "")
	recordID := routeVarsMap.getOrDefault("id", "")
	show("DeleteRecord, table name, recId: ", tableName, recordID)

	// table
	table, err := srv.getTable(tableName)
	if exitOnError(err, http.StatusInternalServerError, "DeleteRecord, wrong table") {
		return
	}

	// query
	execResult, err := srv.DB.Exec(
		fmt.Sprintf("DELETE FROM %s WHERE %s = ?", tableName, table.Pk),
		recordID,
	)
	if exitOnError(err, http.StatusInternalServerError, "`DELETE` failed") {
		return
	}

	affectedCount, err := execResult.RowsAffected()
	if exitOnError(err, http.StatusInternalServerError, "RowsAffected failed") {
		return
	}

	// result
	respData := &GenericMap{
		"response": GenericMap{
			"deleted": affectedCount,
		},
	}
	writeSuccess(http.StatusOK, respData, w)
}

func (srv *MysqlExplorerHttpHandlers) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	/*
		* POST /$table/$id - обновляет запись, данные приходят в теле запроса (POST-параметры)
		r.HandleFunc("/{table}/{id}", srv.UpdateRecord).Methods("POST")

		Case{
			Path:   "/items/3",
			Method: http.MethodPost,
			RequestBody: GenericMap{
				"updated": "autotests",
			},
			ExpectedRespBody: GenericMap{
				"response": GenericMap{
					"updated": 1,
				},
			},
		},

	*/

	defer recoverPanic(w)

	exitOnError := func(err error) bool {
		if err != nil {
			show("UpdateRecord, error: ", err)
			writeError(http.StatusInternalServerError, "UpdateRecord failed", w)
			return true
		}
		return false
	}

	var wrongId = func(selectedId string, newId any) bool {
		if newId == nil {
			return false
		}
		if fmt.Sprintf("%s", selectedId) == fmt.Sprintf("%s", newId) {
			return false
		}
		return true
	}

	var invalid = func(v any, colName string, tab Table) bool {
		c, err := tab.getColumn(colName)
		if err != nil {
			return true
		}

		if c.Type.IsValidValue(v) {
			return false
		}

		return true
	}

	// params
	routeVarsMap := MapSS(mux.Vars(r))
	tableName := routeVarsMap.getOrDefault("table", "")
	recordID := routeVarsMap.getOrDefault("id", "")
	show("UpdateRecord, table name, recId: ", tableName, recordID)

	// table
	table, err := srv.getTable(tableName)
	if exitOnError(err) {
		return
	}

	// record
	bodyBytes, err := io.ReadAll(r.Body)
	if exitOnError(err) {
		return
	}
	record := TableRecord{}
	err = json.Unmarshal(bodyBytes, &record)
	if exitOnError(err) {
		return
	}

	// validate
	if wrongId(recordID, record[table.Pk]) {
		writeError(http.StatusBadRequest, "field "+table.Pk+" have invalid type", w)
		return
	}

	// query
	var updateCols, updateVals = make([]string, 0, len(record)), make([]any, 0, len(record))
	for k, v := range record {
		updateCols = append(updateCols, fmt.Sprintf("%s = ?", k))
		updateVals = append(updateVals, v)
		// validate
		if invalid(v, k, table) {
			writeError(http.StatusBadRequest, "field "+k+" have invalid type", w)
			return
		}
	}
	updateQuery := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s = ?",
		tableName,
		strings.Join(updateCols, ", "),
		table.Pk,
	)

	// sql
	execResult, err := srv.DB.Exec(
		updateQuery,
		append(updateVals, recordID)...,
	)
	if exitOnError(err) {
		return
	}
	affectedCount, err := execResult.RowsAffected()
	if exitOnError(err) {
		return
	}

	// result
	respData := &GenericMap{
		"response": GenericMap{
			"updated": affectedCount,
		},
	}
	writeSuccess(http.StatusOK, respData, w)
}

func (srv *MysqlExplorerHttpHandlers) CreateRecord(w http.ResponseWriter, r *http.Request) {
	/*
				* PUT /$table - создаёт новую запись, данный по записи в теле запроса (POST-параметры)

				r.HandleFunc("/{table}/", srv.CreateRecord).Methods("PUT")

				Path:   "/items/",
				Method: http.MethodPut,
				RequestBody: GenericMap{
					"id":          42, // auto increment primary key игнорируется при вставке
					"title":       "db_crud",
					"description": "",
				},
				ExpectedRespBody: GenericMap{
					"response": GenericMap{
						"id": 3,
					},
				},

				Case{
					Path:   "/users/",
					Method: http.MethodPut,
					RequestBody: GenericMap{
						"user_id":    2,
						"login":      "qwerty'",
						"password":   "love\"",
						"unkn_field": "love",
					},
					ExpectedRespBody: GenericMap{
						"response": GenericMap{
							"user_id": 2,
						},
					},
				}, // 26
		        Actual : map[string]interface {}{"response":map[string]interface {}{"id":4}}
		        Expected: map[string]interface {}{"response":map[string]interface {}{"user_id":2}}
	*/

	defer recoverPanic(w)

	exitOnError := func(err error, status int, msg string) bool {
		if err != nil {
			show("CreateRecord, error: ", err)
			writeError(status, msg, w)
			return true
		}
		return false
	}

	// params
	routeVarsMap := MapSS(mux.Vars(r))
	tableName := routeVarsMap.getOrDefault("table", "")
	show("CreateRecord, table name: ", tableName)

	// table
	table, err := srv.getTable(tableName)
	if exitOnError(err, http.StatusInternalServerError, "wrong table") {
		return
	}
	// show("CreateRecord, table legit: ", table.Name)

	// record
	bodyBytes, err := io.ReadAll(r.Body)
	if exitOnError(err, http.StatusInternalServerError, "can't read body bytes") {
		return
	}
	record := TableRecord{}
	err = json.Unmarshal(bodyBytes, &record)
	if exitOnError(err, http.StatusInternalServerError, "can't unmarshal body") {
		return
	}

	// "Field 'email' doesn't have a default value"};
	for _, c := range table.Columns {
		_, isIn := record[c.Field]
		if !isIn && !c.Null {
			record[c.Field] = c.Type.NewVar()
		}
	}

	// query
	var insertCols, insertVals = make([]string, 0, len(record)), make([]any, 0, len(record))
	for k, v := range record {
		if k == table.Pk {
			continue // skip key, it's autoincrement
		}
		_, err := table.getColumn(k)
		if err != nil {
			continue // skip unknown columns
		}
		// show("column legit: ", c.Field)

		insertCols = append(insertCols, fmt.Sprintf("`%s`", k))
		insertVals = append(insertVals, v)
	}
	if len(insertCols) < 1 {
		writeError(http.StatusInternalServerError, "Nothing to insert", w)
		return
	}

	placeholders := strings.Repeat("?, ", len(insertCols))
	insertQuery := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(insertCols, ", "),
		placeholders[:len(placeholders)-2],
	)

	// sql
	execResult, err := srv.DB.Exec(
		insertQuery,
		insertVals...,
	)
	if exitOnError(err, http.StatusInternalServerError, "DB.Exec failed") {
		return
	}

	// result
	lastId, err := execResult.LastInsertId()
	panicOnError("LastInsertId failed", err)
	resultRef := &GenericMap{
		"response": GenericMap{
			table.Pk: lastId,
		},
	}
	writeSuccess(http.StatusOK, resultRef, w)
}

func (srv *MysqlExplorerHttpHandlers) ReadRecord(w http.ResponseWriter, r *http.Request) {
	/*
		* GET /$table/$id - возвращает информацию о самой записи или 404

		Path: "/items/1",
		ExpectedRespBody: GenericMap{
			"response": GenericMap{
				"record": GenericMap{
					"id":          1,
					"title":       "database/sql",
					"description": "Рассказать про базы данных",
					"updated":     "rvasily",
				},
			},
		},
	*/

	defer recoverPanic(w)

	exitOnError := func(err error) bool {
		if err != nil {
			show("ReadRecord, error: ", err)
			writeError(http.StatusNotFound, "record not found", w)
			return true
		}
		return false
	}

	// params
	routeVarsMap := MapSS(mux.Vars(r))
	tableName := routeVarsMap.getOrDefault("table", "")
	recordID := routeVarsMap.getOrDefault("id", "")
	show("ReadRecord, table name, recId: ", tableName, recordID)

	// table
	var table, err = srv.getTable(tableName)
	if exitOnError(err) {
		return
	}
	var keyName = table.Pk

	// records
	row := srv.DB.QueryRow(fmt.Sprintf("SELECT * FROM %s where %s = ?", tableName, keyName), recordID)
	values := table.NewRow()
	err = row.Scan(values...)
	if exitOnError(err) {
		return
	}

	resultRef := &GenericMap{
		"response": GenericMap{
			"record": table.NewRecord(values), // map fieldName:fieldValue
		},
	}

	writeSuccess(http.StatusOK, resultRef, w)
}

func (srv *MysqlExplorerHttpHandlers) ReadTable(w http.ResponseWriter, r *http.Request) {
	/*
		* GET /$table?limit=5&offset=7 - возвращает список из 5 записей (limit) начиная с 7-й (offset) из таблицы $table.
			limit по-умолчанию 5, offset 0

	*/
	defer recoverPanic(w)

	exitOnError := func(err error) bool {
		if err != nil {
			show("ReadTable, error: ", err)
			writeError(http.StatusNotFound, "unknown table", w)
			return true
		}
		return false
	}

	// params
	routeVarsMap := MapSS(mux.Vars(r))
	tableName := routeVarsMap.getOrDefault("table", "")

	// если пришло не число на вход - берём дефолтное значене для лимита-оффсета
	limit := getOrDefaultInt(r.URL.Query(), "limit", "5")
	offset := getOrDefaultInt(r.URL.Query(), "offset", "0")
	show("ReadTable, name: ", tableName, fmt.Sprintf("limit %v, offset %v", limit, offset))

	// table
	var table, err = srv.getTable(tableName)
	if exitOnError(err) {
		return
	}

	// records
	rows, err := srv.DB.Query(fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ?", tableName), limit, offset)
	if exitOnError(err) {
		return
	}
	defer rows.Close()

	var records = make([]TableRecord, 0, 16)
	for rows.Next() {
		row := table.NewRow() // values
		err := rows.Scan(row...)
		if exitOnError(err) {
			return
		}
		records = append(records, table.NewRecord(row)) // map fieldName:fieldValue
	}
	// show("ReadTable, records: ", records)

	resultRef := &GenericMap{
		"response": GenericMap{
			"records": records,
		},
	}

	writeSuccess(http.StatusOK, resultRef, w)
}

func (srv *MysqlExplorerHttpHandlers) ListTables(w http.ResponseWriter, r *http.Request) {
	defer recoverPanic(w)

	exitOnError := func(err error) bool {
		if err != nil {
			show("ListTables, error: ", err)
			writeError(http.StatusInternalServerError, "probably DB access failed", w)
			return true
		}
		return false
	}

	allTables, err := srv.GetTableNames()
	if exitOnError(err) {
		return
	}
	show("ListTables: ", allTables)

	resultRef := &GenericMap{
		"response": GenericMap{
			"tables": allTables,
		},
	}

	writeSuccess(http.StatusOK, resultRef, w)
}

func writeSuccess(status int, obj interface{}, w http.ResponseWriter) {
	// bytes, err := json.Marshal(ApiSuccessResponse{"", obj})
	bytes, err := json.Marshal(obj)
	writeResponse(status, err, bytes, w)
}

func writeError(status int, message string, w http.ResponseWriter) {
	bytes, err := json.Marshal(ApiErrorResponse{message})
	writeResponse(status, err, bytes, w)
}

func writeResponse(status int, err error, msg []byte, w http.ResponseWriter) {
	panicOnError("writeResponse, got an error: ", err) // TODO: replace panic with StatusInternalServerError
	w.WriteHeader(status)
	w.Write(msg)
}

// recoverPanic concocted for using in `defer recoverPanic ...` in http handlers
func recoverPanic(w http.ResponseWriter) {
	if err := recover(); err != nil {
		debug.PrintStack()
		show("recover from error: ", err)
		writeError(http.StatusInternalServerError, "Internal server error", w)
	}
}

type ApiSuccessResponse struct {
	Error    string      `json:"error"`
	Response interface{} `json:"response"`
}

type ApiErrorResponse struct {
	Error string `json:"error"`
}

type GenericMap map[string]interface{}

func getOrDefault(values url.Values, key string, defaultValue string) string {
	items, ok := values[key]
	if !ok {
		return defaultValue
	}

	if len(items) == 0 {
		return defaultValue
	}

	return items[0] // TODO: or find first not empty
}

func getOrDefaultInt(values url.Values, key string, defaultValue string) string {
	items, ok := values[key]
	if !ok {
		return defaultValue
	}

	if len(items) == 0 {
		return defaultValue
	}

	v := items[0] // TODO: or find first not empty
	_, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}

	return v
}

type MapSS map[string]string

func (m MapSS) getOrDefault(key, dflt string) string {
	if v, isIn := m[key]; isIn {
		return v
	}
	return dflt
}

func contains(str string, lst []string) bool {
	return slices.Contains(lst, str)
}

func split(str, sep string) []string {
	return strings.Split(str, sep)
}

// panicOnError throw the panic with given error and msg prefix, if err != nil
func panicOnError(msg string, err error) {
	if err != nil {
		panic(msg + ": " + err.Error())
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
		// line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
}
