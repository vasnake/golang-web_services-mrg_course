package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	// "log"
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

	fmt.Fprintln(outFileRef, headerText)

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
						show("struct parsed, structMeta: ", sm, err)
						// TODO: generate: struct parser, struct validator
						// actually: add to collection
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
			show("func parsed, funcMeta: ", fm, err)
			// TODO: generate: http handlers (one reciever: one mux and 1..n handlers)
			// actually: add to collection
		default:
			show("unknown Decl: ", topIdx, topDecl)
		} // end Decl type switch
	}

	// TODO: use collections of structs and funcs to generate output code

	panic("not yet")
}

type StructMeta struct {
	name   string
	fields []string
	// each field: tag options
}

func parseStruct(ts *ast.TypeSpec, st *ast.StructType) (*StructMeta, error) {
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
	for fidx, field := range st.Fields.List {
		show("struct field: ", fidx, field, field.Tag)
		if field.Tag != nil && startsWith(field.Tag.Value, apiValidatorTagPrefix) {
			show("lets roll ...", field.Tag.Value) // lets roll ...string(`apivalidator:"min=1,max=50"`);
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
					show("required")
				case "min":
					show("min: ", elems[1])
				case "max":
					show("max: ", elems[1])
				case "paramname":
					show("read field from: ", elems[1])
				case "enum":
					show("enum: ", strings.Split(elems[1], "|"))
				case "default":
					show("default value: ", elems[1])
				default:
					show("unknown tag: ", elems)
				}
			}
		}
	}
	return nil, nil
}

type FuncMeta struct {
	funcName     string
	recieverName string
	paramName    string
	// url, auth, method
}

func parseFunc(fd *ast.FuncDecl) (*FuncMeta, error) {
	show("processFunc: ", fd)
	/*
		parse:
		func declaration, with comment `// apigen:api ...`, e.g:
		// apigen:api {"url": "/user/create", "auth": true, "method": "POST"}
		func (srv *OtherApi) Create(ctx context.Context, in OtherCreateParams) (*OtherUser, error) { ... }
	*/
	return nil, nil
}

func usage() {
	fmt.Println(usageText)
}

const (
	_ = iota
	notEnoughArgumentsErrorCode
	parserErrorCode
	createFileErrorCode

	apiValidatorTagPrefix = "`apivalidator:"

	usageText = `
Program should be executed like so: go build handlers_gen/* && ./codegen api.go api_handlers.go
where:
- api.go: internal API implementation,
- api_handlers.go: filename for generated code, file will be overwritten without warning.
`

	headerText = `package main

import (
	"errors"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
)	
`
)

func startsWith(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
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
