package main

import (
	// "encoding/json"
	"errors"
	// "fmt"
	"net/http"
	"net/url"
	"runtime/debug"
	// "slices"
	"strconv"
	// "strings"
	// "time"
	// _ "codegen/hand_made_adapters"
)

/*
Чтобы понять, что должен делать кодоген, надо предварительно реализовать логику руками.
Done.

Что надо делать, в принципе:
в апи есть реализация методов на структуре;
эти методы делают работу по обработке запросов;
наша задача: реализовать адаптеры для этих методов,
чтобы снаружи можно было прикрутить веб-интерфейс (http модуль).
Тесты задают спеку на этот веб-интерфейс.
Теги и комментарии кодегена задают спеку адаптеров.

- `ServeHTTP` - принимает все методы из мультиплексора, если нашлось - вызывает `handler$methodName`, если нет - говорит `404`
- `handler$methodName` - обёртка над методом структуры `$methodName` - осуществляет все проверки, выводит ошибки или результат в формате `JSON`
- `$methodName` - непосредственно метод структуры ... Его генерировать не нужно, он уже есть.

Пример тега к полю структуры, определяет парсинг и/или вализацию значений:
Class    string `apivalidator:"enum=warrior|sorcerer|rouge,default=warrior"`
встречающиеся поля тега:
- required
- min
- max
- paramname
- enum
- default

Комментарий-метка http хендлера:
префикс `apigen:api` за которым следует json с полями url, auth, method. method опционален, если не указан, то любой.
Пример:
// apigen:api {"url": "/user/create", "auth": true, "method": "POST"}
func ...

Подробности в ридми.
*/

// ServeHTTP implements http.Handler
func (srv *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	show("\n\nOtherApi.ServeHTTP ...")
	/*
		handlers:
		apigen:api {"url": "/user/create", "auth": true, "method": "POST"}
		func (srv *OtherApi) Create(ctx context.Context, in OtherCreateParams) (*OtherUser, error) {...}
	*/

	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			show("ServeHTTP recover from error: ", err)
			writeError(http.StatusInternalServerError, "Internal server error", w)
		}
	}()

	switch r.URL.Path {
	case "/user/create": // {"url": "/user/create", "auth": true, "method": "POST"}
		srv.handlerCreate(w, r)
	default:
		show("OtherApi unknown url:", r.URL.Path)
		writeError(http.StatusNotFound, "unknown method", w)
	}
}

// ServeHTTP implements http.Handler
func (srv *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	show("\n\nMyApi.ServeHTTP ...")
	/*
		handlers:
		apigen:api {"url": "/user/profile", "auth": false}
		func (srv *MyApi) Profile(ctx context.Context, in ProfileParams) (*User, error) {...}

		apigen:api {"url": "/user/create", "auth": true, "method": "POST"}
		func (srv *MyApi) Create(ctx context.Context, in CreateParams) (*NewUser, error) {...}
	*/
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			show("ServeHTTP recover from error: ", err)
			writeError(http.StatusInternalServerError, "Internal server error", w)
		}
	}()

	switch r.URL.Path {
	case "/user/profile": // {"url": "/user/profile", "auth": false}
		srv.handlerProfile(w, r)
	case "/user/create": // {"url": "/user/create", "auth": true, "method": "POST"}
		srv.handlerCreate(w, r)
	default:
		show("MyApi unknown url:", r.URL.Path)
		writeError(http.StatusNotFound, "unknown method", w)
	}
}

func (srv *OtherApi) handlerCreate(w http.ResponseWriter, r *http.Request) {
	// apigen:api {"url": "/user/create", "auth": true, "method": "POST"} // method is optional
	// func (srv *OtherApi) Create(ctx context.Context, in OtherCreateParams) (*OtherUser, error) {...}

	// method check
	if r.Method != "POST" {
		show("handlerCreate wrong method, required POST, got: ", r.Method)
		writeError(http.StatusNotAcceptable, "bad method", w)
		return
	}

	// auth check
	if !isAuthenticated(r) {
		show("handlerCreate auth required: ", r.Header)
		writeError(http.StatusForbidden, "unauthorized", w)
		return
	}

	// parse params
	r.ParseForm()
	paramsRef := new(OtherCreateParams)
	err := paramsRef.fillFrom(r.Form)
	if err != nil {
		show("handlerCreate can't parse params: ", r.Form)
		writeError(http.StatusBadRequest, err.Error(), w)
		return
	}

	// validate params
	err = paramsRef.validate()
	if err != nil {
		show("handlerCreate invalid params: ", paramsRef)
		writeError(http.StatusBadRequest, err.Error(), w)
		return
	}

	// process
	resultRef, err := srv.Create(r.Context(), *paramsRef)
	if err != nil {
		show("handlerCreate error from srv: ", err)
		writeSrvError(err, w)
		return
	}

	// response
	writeSuccess(http.StatusOK, resultRef, w)
}

func (srv *MyApi) handlerCreate(w http.ResponseWriter, r *http.Request) {
	// apigen:api {"url": "/user/create", "auth": true, "method": "POST"}
	// func (srv *MyApi) Create(ctx context.Context, in CreateParams) (*NewUser, error) {...}

	// method check
	if r.Method != "POST" {
		show("handlerCreate wrong method, required POST, got: ", r.Method)
		writeError(http.StatusNotAcceptable, "bad method", w)
		return
	}

	// auth check
	if !isAuthenticated(r) {
		show("handlerCreate auth required: ", r.Header)
		writeError(http.StatusForbidden, "unauthorized", w)
		return
	}

	// parse params
	r.ParseForm()
	paramsRef := new(CreateParams)
	err := paramsRef.fillFrom(r.Form)
	if err != nil {
		show("handlerCreate can't parse params: ", r.Form)
		writeError(http.StatusBadRequest, err.Error(), w)
		return
	}

	// validate params
	err = paramsRef.validate()
	if err != nil {
		show("handlerCreate invalid params: ", paramsRef)
		writeError(http.StatusBadRequest, err.Error(), w)
		return
	}

	// process
	resultRef, err := srv.Create(r.Context(), *paramsRef)
	if err != nil {
		show("handlerCreate error from srv: ", err)
		writeSrvError(err, w)
		return
	}

	// response
	writeSuccess(http.StatusOK, resultRef, w)
}

func (srv *MyApi) handlerProfile(w http.ResponseWriter, r *http.Request) {
	// apigen:api {"url": "/user/profile", "auth": false}
	// func (srv *MyApi) Profile(ctx context.Context, in ProfileParams) (*User, error) {...}

	// method check
	// no method

	// auth check
	// no auth

	// parse params
	r.ParseForm()
	paramsRef := new(ProfileParams)
	err := paramsRef.fillFrom(r.Form)
	if err != nil {
		show("handlerProfile can't parse params: ", r.Form)
		writeError(http.StatusBadRequest, err.Error(), w)
		return
	}

	// validate params
	err = paramsRef.validate()
	if err != nil {
		show("handlerProfile invalid params: ", paramsRef)
		writeError(http.StatusBadRequest, err.Error(), w)
		return
	}

	// process
	resultRef, err := srv.Profile(r.Context(), *paramsRef)
	if err != nil {
		show("handlerProfile error from srv: ", err)
		writeSrvError(err, w)
		return
	}

	// response
	writeSuccess(http.StatusOK, resultRef, w)
}

func (ppref *ProfileParams) fillFrom(params url.Values) error {
	// Login string `apivalidator:"required"`
	ppref.Login = getOrDefault(params, "login", "")
	return nil
}

func (ocpref *OtherCreateParams) fillFrom(params url.Values) error {
	/*
		Username string `apivalidator:"required,min=3"`
		Name     string `apivalidator:"paramname=account_name"`
		Class    string `apivalidator:"enum=warrior|sorcerer|rouge,default=warrior"`
		Level    int    `apivalidator:"min=1,max=50"`
	*/
	var err error

	// Username string `apivalidator:"required,min=3"`
	ocpref.Username = getOrDefault(params, "username", "")

	// Name     string `apivalidator:"paramname=account_name"`
	ocpref.Name = getOrDefault(params, "account_name", "")

	// Class    string `apivalidator:"enum=warrior|sorcerer|rouge,default=warrior"`
	ocpref.Class = getOrDefault(params, "class", "warrior")

	// Level    int    `apivalidator:"min=1,max=50"`
	ocpref.Level, err = strconv.Atoi(getOrDefault(params, "level", ""))
	if err != nil {
		return levelIntError
	}

	return nil
}

func (cpref *CreateParams) fillFrom(params url.Values) error {
	/*
		Login  string `apivalidator:"required,min=10"`
		Name   string `apivalidator:"paramname=full_name"`
		Status string `apivalidator:"enum=user|moderator|admin,default=user"`
		Age    int    `apivalidator:"min=0,max=128"`

		// parsing
		* `paramname` - если указано - то брать из параметра с этим именем, иначе `lowercase` от имени
		* `default` - если указано и приходит пустое значение (значение по-умолчанию) - устанавливать то что написано указано в `default`
	*/
	var err error

	// Login  string `apivalidator:"required,min=10"`
	cpref.Login = getOrDefault(params, "login", "")

	// Name   string `apivalidator:"paramname=full_name"`
	cpref.Name = getOrDefault(params, "full_name", "")

	// Status string `apivalidator:"enum=user|moderator|admin,default=user"`
	cpref.Status = getOrDefault(params, "status", "user")

	// Age    int    `apivalidator:"min=0,max=128"`
	cpref.Age, err = strconv.Atoi(getOrDefault(params, "age", ""))
	if err != nil {
		return ageIntError
	}

	return nil
}

func (ppref *ProfileParams) validate() error {
	// Login string `apivalidator:"required"`
	if ppref.Login == "" { // required
		return loginRequiredError
	}

	return nil
}

func (cpref *CreateParams) validate() error {
	/*
		Login  string `apivalidator:"required,min=10"`
		Name   string `apivalidator:"paramname=full_name"`
		Status string `apivalidator:"enum=user|moderator|admin,default=user"`
		Age    int    `apivalidator:"min=0,max=128"`
	*/
	// Login  string `apivalidator:"required,min=10"`
	if cpref.Login == "" { // required
		return loginRequiredError
	}
	if len(cpref.Login) < 10 { // min
		return loginMinLenError
	}

	// Name   string `apivalidator:"paramname=full_name"`
	// no validation rules

	// Status string `apivalidator:"enum=user|moderator|admin,default=user"`
	if !contains(cpref.Status, []string{"user", "moderator", "admin"}) { // enum
		return statusEnumError
	}

	// Age    int    `apivalidator:"min=0,max=128"`
	if cpref.Age < 0 { // min
		return ageMinError
	}
	if cpref.Age > 128 { // max
		return ageMaxError
	}

	return nil
}

func (ocpref *OtherCreateParams) validate() error {
	/*
		Username string `apivalidator:"required,min=3"`
		Name     string `apivalidator:"paramname=account_name"`
		Class    string `apivalidator:"enum=warrior|sorcerer|rouge,default=warrior"`
		Level    int    `apivalidator:"min=1,max=50"`
	*/

	// Username string `apivalidator:"required,min=3"`
	if ocpref.Username == "" { // required
		return usernameRequiredError
	}
	if len(ocpref.Username) < 3 { // min
		return usernameMinLenError
	}

	// Name     string `apivalidator:"paramname=account_name"`
	// no validation rules

	// Class    string `apivalidator:"enum=warrior|sorcerer|rouge,default=warrior"`
	if !contains(ocpref.Class, []string{"warrior", "sorcerer", "rouge"}) { // enum
		return classEnumError
	}

	// Level    int    `apivalidator:"min=1,max=50"`
	if ocpref.Level < 1 { // min
		return levelMinError
	}
	if ocpref.Level > 50 { // max
		return levelMaxError
	}

	return nil
}

var (
	loginRequiredError = errors.New("`login`: value required")
	loginMinLenError   = errors.New("login len must be >= 10")
	statusEnumError    = errors.New("status must be one of [user, moderator, admin]")
	ageIntError        = errors.New("age must be int")
	ageMinError        = errors.New("age must be >= 0")
	ageMaxError        = errors.New("age must be <= 128")

	usernameRequiredError = errors.New("`username`: value required")
	usernameMinLenError   = errors.New("username len must be >= 3")
	classEnumError        = errors.New("class must be one of [warrior, sorcerer, rouge]")
	levelIntError         = errors.New("level must be int")
	levelMinError         = errors.New("level must be >= 1")
	levelMaxError         = errors.New("level must be <= 50")
)
