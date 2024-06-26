package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"runtime/debug"
	"slices"
	"strings"
	ttmpl "text/template"

	"os"
	"time"
)

func main() {
	// go build handlers_gen/* && ./codegen api.go api_handlers.go

	show("os.Args: ", os.Args)
	// os.Args: []string([./codegen api.go api_handlers.go]);
	if len(os.Args) < 3 {
		usage()
		os.Exit(notEnoughArgumentsErrorCode)
	}

	var inputFileName = os.Args[1]
	var outputFileName = os.Args[2]

	var nodeRef, err = parser.ParseFile(token.NewFileSet(), inputFileName, nil, parser.ParseComments)
	show("node, err:", nodeRef, err)
	if err != nil {
		show("parser.ParseFile failed: ", err)
		os.Exit(parserErrorCode)
	}

	outFileRef, err := os.Create(outputFileName)
	show("out file: ", outFileRef)
	if err != nil {
		show("os.Create failed: ", err)
		os.Exit(createFileErrorCode)
	}

	fmt.Fprintln(outFileRef, headerText, "")

	var taggedStructs = []ApiValidatorStructMeta{}
	var tryAppendS = func(sm *ApiValidatorStructMeta, err error) {
		show("struct parsed, structMeta, error: ", sm, err)
		if err != nil {
			show("parseStruct failed: ", err)
			os.Exit(parseStructErrorCode)
		}
		if sm != nil {
			taggedStructs = append(taggedStructs, *sm)
		}
	}

	var markedFuncs = []ApiGenFuncMeta{}
	var tryAppendF = func(fm *ApiGenFuncMeta, err error) {
		show("func parsed, funcMeta, error: ", fm, err)
		if err != nil {
			show("parseFunc failed: ", err)
			os.Exit(parseFuncErrorCode)
		}
		if fm != nil {
			markedFuncs = append(markedFuncs, *fm)
		}
	}

	for topIdx, topDecl := range nodeRef.Decls {
		/*
			I need:

			1) struct fields, with tags `apivalidator: ... `, e.g:
			type CreateParams struct {
				Login  string `apivalidator:"required,min=10"`
				...
			}

			2) func declaration, with comment `// apigen:api ...`, e.g:
			// apigen:api {"url": "/user/create", "auth": true, "method": "POST"}
			func (srv *OtherApi) Create(ctx context.Context, in OtherCreateParams) (*OtherUser, error) { ... }

		*/
		// show("top-level declaration: ", topIdx, topDecl)
		switch topDecl.(type) {
		case *ast.GenDecl:
			var tgd = topDecl.(*ast.GenDecl)
			// show("got GenDecl, specs: ", topIdx, tgd.Specs)
			for specIdx, spec := range tgd.Specs {
				switch spec.(type) {
				case *ast.TypeSpec:
					var ts = spec.(*ast.TypeSpec)
					// show("got TypeSpec: ", specIdx, ts)
					switch ts.Type.(type) {
					case *ast.StructType:
						var st = ts.Type.(*ast.StructType)
						// show("got StructType: ", specIdx, st)
						var sm, err = parseStruct(ts, st)
						tryAppendS(sm, err)
					default:
						// show("unknown type: ", specIdx, ts.Type)
					}
				default:
					show("unknown spec: ", specIdx, spec)
				}
			} // end decl.specs loop

		case *ast.FuncDecl:
			var tfd = topDecl.(*ast.FuncDecl)
			// show("got FunDecl, name: ", topIdx, tfd.Name)
			var fm, err = parseFunc(tfd)
			tryAppendF(fm, err)

		default:
			show("unknown Decl: ", topIdx, topDecl)
		} // end Decl type switch
	}

	text, err := generateHandlers(markedFuncs, taggedStructs)
	if err != nil {
		show("generateHandlers failed: ", err)
		os.Exit(generateHandlersErrorCode)
	}

	show("writing generated code ...")
	fmt.Fprintln(outFileRef, text, "")
	show("success")
}

func generateHandlers(funcs []ApiGenFuncMeta, structs []ApiValidatorStructMeta) (string, error) {
	show("generateHandlers: ", funcs, structs)
	var buffer = new(strings.Builder)

	/*
		for each func reciever: create ServeHTTP function with as many route handlers as there are routes
		func (srv *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
			...
			case "/user/create":
				srv.handlerCreate(w, r)
			...
		}
	*/
	var recievers = make([]string, 0, 3)
	for _, fm := range funcs {
		recievers = append(recievers, fm.RecieverName)
	}
	// show("all recievers: ", recievers)
	recievers = distinct(recievers)
	// show("distinct recievers: ", recievers)

	for _, rcv := range recievers {
		var rcvRoutes = filterByReciever(funcs, rcv)
		// show("reciever routes: ", rcv, rcvRoutes)
		var text, err = renderServeHTTPTemplate(rcv, rcvRoutes)
		if err != nil {
			return "", fmt.Errorf("generateHandlers failed: %v", err)
		}
		buffer.WriteString(text)
	}

	/*
		for each func: create handler
		func (srv *OtherApi) handlerCreate(w http.ResponseWriter, r *http.Request) { ... }
	*/
	for _, fm := range funcs {
		var text, err = renderHandlerMethodTemplate(fm)
		if err != nil {
			return "", fmt.Errorf("generateHandlers failed: %v", err)
		}
		buffer.WriteString(text)
	}

	/*
		for each struct: create form parser
		func (ocpref *OtherCreateParams) fillFrom(params url.Values) error { ... }
		and data validator
		func (ocpref *OtherCreateParams) validate() error { ... }
	*/
	for _, st := range structs {
		var text, err = renderFillFormTemplate(st)
		if err != nil {
			return "", fmt.Errorf("generateHandlers failed: %v", err)
		}
		buffer.WriteString(text)

		text, err = renderValidateTemplate(st)
		if err != nil {
			return "", fmt.Errorf("generateHandlers failed: %v", err)
		}
		buffer.WriteString(text)
	}

	return buffer.String(), nil
}

func filterByReciever(xs []ApiGenFuncMeta, rcv string) []ApiGenFuncMeta {
	var ys = make([]ApiGenFuncMeta, 0, len(xs))
	for _, x := range xs {
		if x.RecieverName == rcv {
			ys = append(ys, x)
		}
	}
	return ys
}

func renderFillFormTemplate(st ApiValidatorStructMeta) (string, error) {
	return renderTemplate("fillForm", fillFormTemplate, st)
}

func renderValidateTemplate(st ApiValidatorStructMeta) (string, error) {
	return renderTemplate("validate", validateTemplate, st)
}

func renderHandlerMethodTemplate(fm ApiGenFuncMeta) (string, error) {
	return renderTemplate("handlerMethod", handlerMethodTemplate, fm)
}

func renderServeHTTPTemplate(reciever string, funcs []ApiGenFuncMeta) (string, error) {
	type serveHTTPTemplateData struct {
		Reciever     string
		RouteHanlers []ApiGenFuncMeta
	}
	return renderTemplate("serveHTTP", serveHTTPTemplate, serveHTTPTemplateData{
		Reciever:     reciever,
		RouteHanlers: funcs,
	})
}

func renderTemplate(name, tmpl string, data any) (string, error) {
	var template = ttmpl.New(name)
	var err error

	template, err = template.Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("renderTemplate, failed template.Parse. %v", err)
	}

	var buffer = new(strings.Builder)

	err = template.Execute(buffer, data)
	if err != nil {
		return "", fmt.Errorf("renderTemplate, failed template.Execute. %v", err)
	}

	buffer.WriteString("\n")
	return buffer.String(), nil
}

func parseStruct(ts *ast.TypeSpec, st *ast.StructType) (*ApiValidatorStructMeta, error) {
	show("processStruct: ", ts, st)
	/*
		parse:
		struct fields, with tags `apivalidator: ... `, e.g:
		type CreateParams struct {
			Login  string `apivalidator:"required,min=10"`
			...
		}
		Login string `apivalidator:"required"`
		Login  string `apivalidator:"required,min=10"`
		Name   string `apivalidator:"paramname=full_name"`
		Status string `apivalidator:"enum=user|moderator|admin,default=user"`
		Age    int    `apivalidator:"min=0,max=128"`
		Username string `apivalidator:"required,min=3"`
		Name     string `apivalidator:"paramname=account_name"`
		Class    string `apivalidator:"enum=warrior|sorcerer|rouge,default=warrior"`
		Level    int    `apivalidator:"min=1,max=50"`
	*/
	// collect: struct_name, field_name, field_type, field_tags (list of parser_validator_rules)
	// if no-fields-with-tag: empty result (not error)
	var structName = ts.Name.Name
	// show("struct name: ", structName)
	var taggedFields = []ApiValidatorFieldMeta{}

	for _, field := range st.Fields.List {
		// show("struct field, tag: ", field, field.Tag)

		if field.Tag != nil && startsWith(field.Tag.Value, apiValidatorTagPrefix) {
			show("field with tag, process ...", field.Tag.Value) // lets roll ...string(`apivalidator:"min=1,max=50"`);
			var err error
			var fieldMeta = NewApiValidatorFieldMeta()
			fieldMeta.FieldName = field.Names[0].Name
			fieldMeta.FieldType, err = decodeFieldTypeFromExpr(field.Type)
			if err != nil {
				return nil, fmt.Errorf("Field type decode problem: %#v; %v", field, err)
			}

			var tagsLine = field.Tag.Value[len(apiValidatorTagPrefix)+1 : len(field.Tag.Value)-2]
			// show("tag stripped: ", tv) // tag stripped: string(min=1,max=50);
			var tagsList = strings.Split(tagsLine, ",")
			// show("list of tags: ", strings.Join(ts, `","`)) // list of tags: string(min=1","max=50);
			// tag `required` w/o value, other: with values

			for _, kv := range tagsList {
				// show("tag pair: ", kv)
				var elems = strings.Split(kv, "=")
				switch elems[0] {
				case "required":
					// show("required")
					fieldMeta.Tag.Required = true
				case "min":
					// show("min: ", elems[1])
					fieldMeta.Tag.Min = elems[1]
				case "max":
					// show("max: ", elems[1])
					fieldMeta.Tag.Max = elems[1]
				case "paramname":
					// show("read field from: ", elems[1])
					fieldMeta.Tag.Paramname = elems[1]
				case "enum":
					// show("enum: ", strings.Split(elems[1], "|"))
					fieldMeta.Tag.Enum = strings.Split(elems[1], "|")
				case "default":
					// show("default value: ", elems[1])
					fieldMeta.Tag.DefaultValue = elems[1]
				default:
					// show("unknown tag: ", elems)
					return nil, fmt.Errorf("Unknown tag: %v, %v", tagsLine, elems[0])
				}
			} // end of field tags iterator

			show("field with tag: ", fieldMeta)
			taggedFields = append(taggedFields, *fieldMeta)
		} // end if field 'have the tag'
	} // end fields iterator

	if len(taggedFields) > 0 {
		return NewApiValidatorStructMeta(structName, taggedFields), nil
	}
	return nil, nil
}

func parseFunc(fd *ast.FuncDecl) (*ApiGenFuncMeta, error) {
	show("processFunc: ", fd)
	/*
		parse:
		func declaration, with comment `// apigen:api ...`, e.g:
		// apigen:api {"url": "/user/create", "auth": true, "method": "POST"}
		func (srv *OtherApi) Create(ctx context.Context, in OtherCreateParams) (*OtherUser, error) { ... }

		decode: funcName, recvTypeName, paramTypeName; apigenTrio(url, auth, method)
	*/
	var comment = fd.Doc.Text()
	if startsWith(comment, apiGenTagPrefix) {
		show("parseFunc, apigen marker found, processing func: ", fd.Name, comment)

		var jsonSpec = comment[len(apiGenTagPrefix):]
		// show("api json: ", jsonSpec)
		var specMap specMap
		var err = json.Unmarshal([]byte(jsonSpec), &specMap)
		if err != nil {
			return nil, fmt.Errorf("parseFunc failed, invalid apigen json: %v. %v", jsonSpec, err)
		}
		// show("api spec map: ", specMap)

		var funcMeta = NewApiGenFuncMeta().fillFromSpec(specMap)
		if funcMeta == nil {
			return nil, fmt.Errorf("parseFunc failed, problems with spec comment. %v", comment)
		}
		show("api meta: ", funcMeta)

		funcMeta.FuncName = fd.Name.Name

		funcMeta.RecieverName = func(r *ast.FieldList) string { // TODO: refactor typeName decoder
			if r == nil || len(r.List) == 0 {
				return ""
			}
			rt, err := decodeAnyTypeFromExpr(r.List[0].Type)
			if err != nil {
				show("decode reciever type failed: ", err, r.List[0].Type)
				return ""
			}
			return rt
		}(fd.Recv)

		funcMeta.ParamName = func(p *ast.FieldList) string {
			for i, f := range p.List { // TODO: access by index `1`
				typeName, err := decodeAnyTypeFromExpr(f.Type)
				if err != nil && i == 1 { // don't care about context param
					show("decode parameter type failed: ", err, f.Type)
					return ""
				}
				if i == 1 {
					return typeName
				} // skip context parameter
			}
			return ""
		}(fd.Type.Params)

		if funcMeta.RecieverName == "" {
			return nil, fmt.Errorf("parseFunc failed, invalid recieverName. %v", funcMeta)
		}
		if funcMeta.ParamName == "" {
			return nil, fmt.Errorf("parseFunc failed, invalid paramName. %v", funcMeta)
		}
		return funcMeta, nil
	} // end if found apigen comment

	return nil, nil
}

type ApiGenFuncMeta struct {
	FuncName     string // empty by default
	RecieverName string // empty by default
	ParamName    string // empty by default
	Url          string // empty by default
	HttpMethod   string // empty by default
	Auth         bool   // false by default
}

func NewApiGenFuncMeta() *ApiGenFuncMeta {
	return &ApiGenFuncMeta{
		Auth: false,
	}
}

func (fm *ApiGenFuncMeta) fillFromSpec(spec specMap) *ApiGenFuncMeta {
	defer func() { // type assertion problems
		if err := recover(); err != nil {
			debug.PrintStack()
			show("fillFromSpec, recover from error: ", err)
			fm = nil
		}
		if fm.Url == "" {
			show("fillFromSpec, url is empty")
			fm = nil
		}
	}()

	// apigen:api {"url": "/user/create", "auth": true, "method": "POST"}
	fm.Url = (spec.getOrDefault("url", "")).(string)
	fm.HttpMethod = (spec.getOrDefault("method", "")).(string)
	fm.Auth = (spec.getOrDefault("auth", false)).(bool)

	return fm
}

type ApiValidatorStructMeta struct {
	StructName   string
	TaggedFields []ApiValidatorFieldMeta
}

func NewApiValidatorStructMeta(name string, fields []ApiValidatorFieldMeta) *ApiValidatorStructMeta {
	return &ApiValidatorStructMeta{
		StructName:   name,
		TaggedFields: fields,
	}
}

type ApiValidatorFieldMeta struct {
	FieldName string
	FieldType string // int, string
	Tag       ApiValidatorTags
}

func NewApiValidatorFieldMeta() *ApiValidatorFieldMeta {
	return &ApiValidatorFieldMeta{
		Tag: *NewApiValidatorTags(),
	}
}
func (fm ApiValidatorFieldMeta) ParamNameOrFieldName() string {
	if fm.Tag.Paramname == "" {
		return strings.ToLower(fm.FieldName)
	}
	return strings.ToLower(fm.Tag.Paramname)
}
func (fm ApiValidatorFieldMeta) DefaultValue() string {
	return fm.Tag.DefaultValue
}
func (fm ApiValidatorFieldMeta) Min() bool {
	return fm.Tag.Min != ""
}
func (fm ApiValidatorFieldMeta) Max() bool {
	return fm.Tag.Max != ""
}
func (fm ApiValidatorFieldMeta) Enum() bool {
	return len(fm.Tag.Enum) > 0
}
func (fm ApiValidatorFieldMeta) EnumGoRepr() string {
	// return `[]string{"user", "moderator", "admin"}`
	return fmt.Sprintf(`[]string{"%s"}`, strings.Join(fm.Tag.Enum, `", "`))
}
func (fm ApiValidatorFieldMeta) EnumListRepr() string {
	// return "[user, moderator, admin]"
	return fmt.Sprintf(`[%s]`, strings.Join(fm.Tag.Enum, `, `))
}

type ApiValidatorTags struct {
	Required     bool     // false by default
	Min          string   // empty by default
	Max          string   // empty by default
	Paramname    string   // empty by default
	Enum         []string // empty by default
	DefaultValue string   // empty by default
}

func NewApiValidatorTags() *ApiValidatorTags {
	return &ApiValidatorTags{
		Required:     false,
		Min:          "",
		Max:          "",
		Paramname:    "",
		Enum:         []string{},
		DefaultValue: "",
	}
}

func usage() {
	fmt.Println("", usageText, "")
}

const (
	_ = iota
	notEnoughArgumentsErrorCode
	parserErrorCode
	createFileErrorCode
	parseStructErrorCode
	parseFuncErrorCode
	generateHandlersErrorCode

	apiValidatorTagPrefix = "`apivalidator:"
	apiGenTagPrefix       = "apigen:api"

	usageText = `Program should be executed like so:
go build handlers_gen/* && ./codegen api.go api_handlers.go
where:
- api.go: internal API implementation,
- api_handlers.go: filename for generated code, file will be overwritten without warning.`

	headerText = `package main

import (
	"errors"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
)`

	serveHTTPTemplate = `// ServeHTTP implements http.Handler
func (srv {{ .Reciever }} ) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			writeError(http.StatusInternalServerError, "Internal server error", w)
		}
	}()

	switch r.URL.Path {
{{ range .RouteHanlers }}
	case "{{ .Url }}":
		srv.handler{{ .FuncName }}(w, r)
{{ end }}
	default:
		writeError(http.StatusNotFound, "unknown method", w)
	}
}`

	handlerMethodTemplate = `// handler{{ .FuncName }} implements http.Handler for '{{ .FuncName }}' method
func (srv {{ .RecieverName }} ) handler{{ .FuncName }}(w http.ResponseWriter, r *http.Request) {
{{ $length := len .HttpMethod }} {{ if gt $length 0 }}
	if r.Method != "{{ .HttpMethod }}" {
		writeError(http.StatusNotAcceptable, "bad method", w)
		return
	}
{{ end }}
{{ if .Auth }}
	if !isAuthenticated(r) {
		writeError(http.StatusForbidden, "unauthorized", w)
		return
	}
{{- end }}
	r.ParseForm()
	paramsRef := new({{ .ParamName }})
	err := paramsRef.fillFrom(r.Form)
	if err != nil {
		writeError(http.StatusBadRequest, err.Error(), w)
		return
	}

	err = paramsRef.validate()
	if err != nil {
		writeError(http.StatusBadRequest, err.Error(), w)
		return
	}

	resultRef, err := srv.{{ .FuncName }}(r.Context(), *paramsRef)
	if err != nil {
		writeSrvError(err, w)
		return
	}

	writeSuccess(http.StatusOK, resultRef, w)
}`

	fillFormTemplate = `// fillFrom write data from 'params' to 'pref'
func (pref *{{ .StructName }}) fillFrom(params url.Values) error {
	var err error = nil
{{ range .TaggedFields }}
{{- if eq .FieldType "string" }}
	pref.{{ .FieldName }} = getOrDefault(params, "{{ .ParamNameOrFieldName }}", "{{ .DefaultValue }}")
{{- else }}
	pref.{{ .FieldName }}, err = strconv.Atoi(getOrDefault(params, "{{ .ParamNameOrFieldName }}", "{{ .DefaultValue }}"))
	if err != nil {
		return errors.New("{{ .ParamNameOrFieldName }} must be int")
	}
{{ end }}
{{ end }}
	return err
}`

	validateTemplate = `// validate check data against set of rules
func (cpref *{{ .StructName }}) validate() error {
{{ range .TaggedFields }}
{{- if eq .FieldType "string" }}
{{ if .Tag.Required }}
	if cpref.{{ .FieldName }} == "" { // required
		return errors.New("{{ .ParamNameOrFieldName }}: value required")
	}
{{ end }}
{{ if .Min }}
	if len(cpref.{{ .FieldName }}) < {{ .Tag.Min }} { // min string
		return errors.New("{{ .ParamNameOrFieldName }} len must be >= {{ .Tag.Min }}")
	}
{{ end }}
{{ if .Max }}
	if len(cpref.{{ .FieldName }}) > {{ .Tag.Max }} {
		return errors.New("{{ .ParamNameOrFieldName }} len must be <= {{ .Tag.Max }}")
	}
{{ end }}
{{ if .Enum }}
	if !contains(cpref.{{ .FieldName }}, {{ .EnumGoRepr }}) { // enum
		return errors.New("{{ .ParamNameOrFieldName }} must be one of {{ .EnumListRepr }}")
	}
{{ end }}
{{- else }}
{{ if .Min }}
	if cpref.{{ .FieldName }} < {{ .Tag.Min }} { // min int
		return errors.New("{{ .ParamNameOrFieldName }} must be >= {{ .Tag.Min }}")
	}
{{ end }}
{{ if .Max }}
	if cpref.{{ .FieldName }} > {{ .Tag.Max }} {
		return errors.New("{{ .ParamNameOrFieldName }} must be <= {{ .Tag.Max }}")
	}
{{ end }}
{{- end }}
{{- end }}
	return nil
}`
)

type specMap map[string]any

func (xs specMap) getOrDefault(key string, dflt any) any {
	var v, exist = xs[key]
	if !exist {
		return dflt
	}
	return v
}

func startsWith(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func decodeFieldTypeFromExpr(expr ast.Expr) (string, error) {
	return decodeTypeFromExpr(expr, []string{"int", "string"})
}

func decodeAnyTypeFromExpr(expr ast.Expr) (string, error) {
	return decodeTypeFromExpr(expr, []string{})
}

func decodeTypeFromExpr(expr ast.Expr, check []string) (string, error) {
	var exprStr = types.ExprString(expr)
	// show("decodeTypeFromExpr: ", expr, exprStr)

	if exprStr == "" {
		return "", fmt.Errorf("decodeTypeFromExpr, failed conversion from Expr to string. Expr: %v", expr)
	}

	if len(check) == 0 || slices.Contains(check, exprStr) {
		return exprStr, nil
	}

	return exprStr, fmt.Errorf("decodeTypeFromExpr, unknown type: %s", exprStr)
}

func distinct(xs []string) []string {
	var ys = make([]string, 0, len(xs))
	slices.Sort(xs)
	for i, _ := range xs {
		if i == 0 || xs[i] != ys[len(ys)-1] {
			ys = append(ys, xs[i])
		}
	}
	return ys
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
