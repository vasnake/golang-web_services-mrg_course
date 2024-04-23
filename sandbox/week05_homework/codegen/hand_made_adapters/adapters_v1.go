package hand_made_adapters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"
	"time"
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

// ServeHTTP implements http.Handler.
func (srv *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	show("\n\nOtherApi.ServeHTTP ...")
	// ts := httptest.NewServer(NewOtherApi())
	/*
		// apigen:api {"url": "/user/create", "auth": true, "method": "POST"}
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
	case "/user/create":
		srv.handleCreate(w, r)
	default:
		show("OtherApi unknown url:", r.URL.Path)
		writeError(http.StatusNotFound, "unknown method", w)
	}
}

// ServeHTTP implements http.Handler
func (srv *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	show("\n\nMyApi.ServeHTTP ...")
	// ts := httptest.NewServer(NewMyApi())
	/*
		// apigen:api {"url": "/user/profile", "auth": false}
		func (srv *MyApi) Profile(ctx context.Context, in ProfileParams) (*User, error) {...}

		// apigen:api {"url": "/user/create", "auth": true, "method": "POST"}
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
	case "/user/profile":
		srv.handlerProfile(w, r)
	case "/user/create":
		srv.handleCreate(w, r)
	default:
		show("MyApi unknown url:", r.URL.Path)
		writeError(http.StatusNotFound, "unknown method", w)
	}
}

func (srv *OtherApi) handleCreate(w http.ResponseWriter, r *http.Request) {
	// apigen:api {"url": "/user/create", "auth": true, "method": "POST"}
	// func (srv *OtherApi) Create(ctx context.Context, in OtherCreateParams) (*OtherUser, error) {...}
	method := "POST" // diff
	if method != "" && method != r.Method {
		show("handleCreate wrong method, required x, got y: (x, y): ", method, r.Method)
		writeError(http.StatusNotAcceptable, "bad method", w)
		return
	}

	authRequired := true // diff
	ok := isAuthenticated(r, authRequired)
	if !ok {
		show("handleCreate auth required: ", r.Header)
		writeError(http.StatusForbidden, "unauthorized", w)
		return
	}

	r.ParseForm()
	params := new(OtherCreateParams) // diff
	err := params.fillFromForm(r.Form)
	if err != nil {
		show("handleCreate invalid params: ", r.Form)
		writeError(http.StatusBadRequest, err.Error(), w)
		return
	}

	res, err := srv.Create(r.Context(), *params) // diff
	if err != nil {
		show("handleCreate error from srv: ", err)
		writeSrvError(err, w)
		return
	}

	writeSuccess(http.StatusOK, res, w)
}

func (srv *MyApi) handleCreate(w http.ResponseWriter, r *http.Request) {
	// apigen:api {"url": "/user/create", "auth": true, "method": "POST"}
	// func (srv *MyApi) Create(ctx context.Context, in CreateParams) (*NewUser, error) {...}
	method := "POST" // diff
	if method != "" && method != r.Method {
		show("handleCreate wrong method, required x, got y: (x, y): ", method, r.Method)
		writeError(http.StatusNotAcceptable, "bad method", w)
		return
	}

	authRequired := true // diff
	ok := isAuthenticated(r, authRequired)
	if !ok {
		show("handleCreate auth required: ", r.Header)
		writeError(http.StatusForbidden, "unauthorized", w)
		return
	}

	r.ParseForm()
	params := new(CreateParams) // diff
	err := params.fillFromForm(r.Form)
	if err != nil {
		show("handleCreate invalid params: ", r.Form)
		writeError(http.StatusBadRequest, err.Error(), w)
		return
	}

	res, err := srv.Create(r.Context(), *params) // diff
	if err != nil {
		show("handleCreate error from srv: ", err)
		writeSrvError(err, w)
		return
	}

	writeSuccess(http.StatusOK, res, w)
}

func (srv *MyApi) handlerProfile(w http.ResponseWriter, r *http.Request) {
	method := ""
	if method != "" && method != r.Method {
		show("handlerProfile wrong method, required x, got y: (x, y): ", method, r.Method)
		writeError(http.StatusNotAcceptable, "bad method", w)
		return
	}

	authRequired := false
	ok := isAuthenticated(r, authRequired)
	if !ok {
		show("handlerProfile auth required: ", r.Header)
		writeError(http.StatusForbidden, "unauthorized", w)
		return
	}

	r.ParseForm()
	params := new(ProfileParams)
	err := params.fillFromForm(r.Form)
	if err != nil {
		show("handlerProfile invalid params: ", r.Form)
		writeError(http.StatusBadRequest, err.Error(), w)
		return
	}

	res, err := srv.Profile(r.Context(), *params)
	if err != nil {
		show("handlerProfile error from srv: ", err)
		writeSrvError(err, w)
		return
	}

	writeSuccess(http.StatusOK, res, w)
}

func (ocpref *OtherCreateParams) fillFromForm(params url.Values) error {
	/*
		Username string `apivalidator:"required,min=3"`
		Name     string `apivalidator:"paramname=account_name"`
		Class    string `apivalidator:"enum=warrior|sorcerer|rouge,default=warrior"`
		Level    int    `apivalidator:"min=1,max=50"`
	*/
	var err error = nil
	var rules []ValidatorRule
	var tmpVal string

	// Username string `apivalidator:"required,min=3"`
	tmpVal = getOrDefault(params, "username", "")
	rules = []ValidatorRule{
		{ValueType: "str", RuleName: "required"},
		{ValueType: "str", RuleName: "min", RuleValue: "3"},
	}
	err = validate(tmpVal, "username", rules)
	if err != nil {
		return err
	}
	ocpref.Username = tmpVal

	// Name     string `apivalidator:"paramname=account_name"`
	tmpVal = getOrDefault(params, "account_name", "")
	rules = []ValidatorRule{}
	err = validate(tmpVal, "name", rules)
	if err != nil {
		return err
	}
	ocpref.Name = tmpVal

	// Class    string `apivalidator:"enum=warrior|sorcerer|rouge,default=warrior"`
	tmpVal = getOrDefault(params, "class", "warrior")
	rules = []ValidatorRule{
		{ValueType: "str", RuleName: "enum", RuleValue: "warrior|sorcerer|rouge"},
	}
	err = validate(tmpVal, "class", rules)
	if err != nil {
		return err
	}
	ocpref.Class = tmpVal

	// Level    int    `apivalidator:"min=1,max=50"`
	tmpVal = getOrDefault(params, "level", "")
	ocpref.Level, err = strconv.Atoi(tmpVal)
	if err != nil {
		return fmt.Errorf("level must be int")
	}
	rules = []ValidatorRule{
		{ValueType: "int", RuleName: "min", RuleValue: "1"},
		{ValueType: "int", RuleName: "max", RuleValue: "50"},
	}
	err = validate(tmpVal, "level", rules)
	if err != nil {
		return err
	}

	return nil
}

func (ppref *ProfileParams) fillFromForm(params url.Values) error {
	var err error = nil
	var rules []ValidatorRule
	var tmpVal string

	// Login string `apivalidator:"required"`
	tmpVal = getOrDefault(params, "login", "")
	rules = []ValidatorRule{
		{ValueType: "str", RuleName: "required"},
	}
	err = validate(tmpVal, "login", rules)
	if err != nil {
		return err
	}
	ppref.Login = tmpVal

	return nil
}

func (cpref *CreateParams) fillFromForm(params url.Values) error {
	/*
		Login  string `apivalidator:"required,min=10"`
		Name   string `apivalidator:"paramname=full_name"`
		Status string `apivalidator:"enum=user|moderator|admin,default=user"`
		Age    int    `apivalidator:"min=0,max=128"`

		// parsing
		* `paramname` - если указано - то брать из параметра с этим именем, иначе `lowercase` от имени
		* `default` - если указано и приходит пустое значение (значение по-умолчанию) - устанавливать то что написано указано в `default`
	*/
	var err error = nil
	var rules []ValidatorRule
	var tmpVal string

	// Login  string `apivalidator:"required,min=10"`
	tmpVal = getOrDefault(params, "login", "")
	rules = []ValidatorRule{
		{ValueType: "str", RuleName: "required"},
		{ValueType: "str", RuleName: "min", RuleValue: "10"},
	}
	err = validate(tmpVal, "login", rules)
	if err != nil {
		return err
	}
	cpref.Login = tmpVal

	// Name   string `apivalidator:"paramname=full_name"`
	tmpVal = getOrDefault(params, "full_name", "")
	rules = []ValidatorRule{}
	err = validate(tmpVal, "name", rules)
	if err != nil {
		return err
	}
	cpref.Name = tmpVal

	// Status string `apivalidator:"enum=user|moderator|admin,default=user"`
	tmpVal = getOrDefault(params, "status", "user")
	rules = []ValidatorRule{
		{ValueType: "str", RuleName: "enum", RuleValue: "user|moderator|admin"},
	}
	err = validate(tmpVal, "status", rules)
	if err != nil {
		return err
	}
	cpref.Status = tmpVal

	// Age    int    `apivalidator:"min=0,max=128"`
	tmpVal = getOrDefault(params, "age", "")
	cpref.Age, err = strconv.Atoi(tmpVal)
	if err != nil {
		return fmt.Errorf("age must be int")
	}
	rules = []ValidatorRule{
		{ValueType: "int", RuleName: "min", RuleValue: "0"},
		{ValueType: "int", RuleName: "max", RuleValue: "128"},
	}
	err = validate(tmpVal, "age", rules)
	if err != nil {
		return err
	}

	return nil
}

func writeSrvError(err error, w http.ResponseWriter) {
	switch err.(type) {
	case ApiError:
		writeError((err.(ApiError)).HTTPStatus, err.Error(), w)
	default:
		writeError(http.StatusInternalServerError, err.Error(), w)
	}
}

func writeSuccess(status int, obj interface{}, w http.ResponseWriter) {
	bs, err := json.Marshal(ApiSuccessResponse{"", obj})
	writeResponse(status, err, bs, w)
}

func writeError(status int, message string, w http.ResponseWriter) {
	bs, err := json.Marshal(ApiErrorResponse{message})
	writeResponse(status, err, bs, w)
}

func writeResponse(status int, err error, msg []byte, w http.ResponseWriter) {
	if err != nil {
		panic(err) // TODO: replace panic with StatusInternalServerError
	}
	w.WriteHeader(status)
	w.Write(msg)
}

type ApiSuccessResponse struct {
	Error    string      `json:"error"`
	Response interface{} `json:"response"`
}

type ApiErrorResponse struct {
	Error string `json:"error"`
}

func isAuthenticated(r *http.Request, requiredAuth bool) bool {
	return !requiredAuth || r.Header.Get("X-Auth") == "100500"
}

func getOrDefault(values url.Values, key string, defaultValue string) string {
	items, ok := values[strings.ToLower(key)]
	if !ok {
		return defaultValue
	}
	if len(items) == 0 {
		return defaultValue
	}
	return items[0] // TODO: or find first not empty
}

type ValidatorRule struct {
	ValueType string
	RuleName  string
	RuleValue string
}

func validate(value, name string, rules []ValidatorRule) error {
	/*
		Нам доступны следующие метки валидатора-заполнятора `apivalidator`:
		* `required` - поле не должно быть пустым (не должно иметь значение по-умолчанию)
		* `paramname` - если указано - то брать из параметра с этим именем, иначе `lowercase` от имени
		* `enum` - "одно из"
		* `default` - если указано и приходит пустое значение (значение по-умолчанию) - устанавливать то что написано указано в `default`
		* `min` - >= X для типа `int`, для строк `len(str)` >=
		* `max` - <= X для типа `int`
	*/
	parseIntOrPanic := func(v string) int {
		i64, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			show("parseIntOrPanic, error: ", v, err)
			panic("parseIntOrPanic, no can do: " + v)
		}
		return int(i64)
	}
	errorRequired := func(name string) error {
		return fmt.Errorf("`" + name + "`: value required")
	}
	errorStringMin := func(name, limit string) error {
		return fmt.Errorf("%s len must be >= %s", name, limit)
		// Expected: main.CR{"error":"login len must be >= 10"}
	}
	errorIntMin := func(name, limit string) error {
		return fmt.Errorf("%s must be >= %s", name, limit)
		// Expected: main.CR{"error":"age must be >= 0"}
	}
	errorStringMax := func(name string) error {
		return fmt.Errorf("`" + name + "`: str value too long")
	}
	errorIntMax := func(name, limit string) error {
		return fmt.Errorf("%s must be <= %s", name, limit)
		// Expected: main.CR{"error":"age must be <= 128"}
	}
	errorNotIn := func(name string, enum []string) error {
		return fmt.Errorf("%s must be one of [%s]", name, strings.Join(enum, ", "))
		// Expected: main.CR{"error":"status must be one of [user, moderator, admin]"}
	}

	for ruleIdx, r := range rules {
		show("validate, rule idx, rule; name, value: ", ruleIdx, r, name, value)
		switch r.RuleName {

		case "required":
			if value == "" {
				return errorRequired(name)
			}

		case "min":
			if r.ValueType == "str" && len(value) < parseIntOrPanic(r.RuleValue) {
				return errorStringMin(name, r.RuleValue)
			}
			if r.ValueType == "int" && parseIntOrPanic(value) < parseIntOrPanic(r.RuleValue) {
				return errorIntMin(name, r.RuleValue)
			}

		case "max":
			if r.ValueType == "str" && len(value) > parseIntOrPanic(r.RuleValue) {
				return errorStringMax(name)
			}
			if r.ValueType == "int" && parseIntOrPanic(value) > parseIntOrPanic(r.RuleValue) {
				return errorIntMax(name, r.RuleValue)
			}

		case "enum":
			// Status string `apivalidator:"enum=user|moderator|admin"`
			if !contains(value, split(r.RuleValue, "|")) {
				return errorNotIn(name, split(r.RuleValue, "|"))
			}

		default:
			show("validate, unknown rule: ", r)
		} // end switch rule name
	}
	return nil
}

func contains(str string, lst []string) bool {
	return slices.Contains(lst, str)
}

func split(str, sep string) []string {
	return strings.Split(str, sep)
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

// ts returns current timestamp in RFC3339 with milliseconds
func ts() string {
	/*
		https://pkg.go.dev/time#pkg-constants
		https://stackoverflow.com/questions/35479041/how-to-convert-iso-8601-time-in-golang
	*/
	const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
	return time.Now().UTC().Format(RFC3339Milli)
}
