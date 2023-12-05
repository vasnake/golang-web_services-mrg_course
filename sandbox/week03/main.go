package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

func jsonDemo() {
	show("program started ...")

	type User struct {
		ID       int
		Username string
		phone    string // private field, can't access from encoding/json
	}

	var jsonStr = `{"id": 42, "username": "rvasily", "phone": "123"}`
	var jsonBytes = []byte(jsonStr) // codec works with bytes, not strings

	u := &User{} // ref to allocated struct
	json.Unmarshal(jsonBytes, u)
	show("Decoded struct x from string s, (x, s): ", u, jsonStr) // load user from json, N.B. empty phone

	u.phone = "987654321"
	result, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	show("Encoded struct x to string s, (x, s): ", u, string(result)) // no phone whatsoever

	show("end of program.")
	/*
	   2023-12-04T10:57:16.672Z: program started ...
	   2023-12-04T10:57:16.672Z: Decoded struct x from string s, (x, s)*main.User(&{42 rvasily }); string({"id": 42, "username": "rvasily", "phone": "123"});
	   2023-12-04T10:57:16.672Z: Encoded struct x to string s, (x, s)*main.User(&{42 rvasily 987654321}); string({"ID":42,"Username":"rvasily"});
	   2023-12-04T10:57:16.672Z: end of program.
	*/
}

func struct_tags() {
	show("program started ...")

	// The encoding of each struct field can be customized by the format string stored under the "json" key in the struct field's tag
	// https://pkg.go.dev/encoding/json#Marshal
	type User struct {
		ID       int `json:"user_id,string"`
		Username string
		Address  string `json:",omitempty"`
		Comnpany string `json:"-"`
	}

	var u = &User{
		ID:       42,
		Username: "rvasily",
		Address:  "test",
		Comnpany: "Mail.Ru Group",
	}

	result, _ := json.Marshal(u)
	show("Encoded struct x to string s, (x, s): ", u, string(result))
	// {"user_id":"42","Username":"rvasily","Address":"test"}

	show("end of program.")
	/*
	   2023-12-04T11:01:34.038Z: program started ...
	   2023-12-04T11:01:34.038Z: Encoded struct x to string s, (x, s)*main.User(&{42 rvasily test Mail.Ru Group}); string({"user_id":"42","Username":"rvasily","Address":"test"});
	   2023-12-04T11:01:34.038Z: end of program.
	*/
}

func dynamicDemo() {
	show("program started ...")

	var jsonStr = `[
		{"id": 17, "username": "iivan", "phone": 0},
		{"id": "17", "address": "none", "company": "Mail.ru"}
	]` // list of map, N.B. mixed id data type

	var jsonBytes = []byte(jsonStr)

	var anyTypeValue interface{}
	json.Unmarshal(jsonBytes, &anyTypeValue)
	show("Decoded json x to an empty interface y, (x, y): ", jsonStr, anyTypeValue)
	// []interface {}([map[id:17 phone:0 username:iivan] map[address:none company:Mail.ru id:17]])

	var universalMap = map[string]interface{}{
		"id":       42,
		"username": "rvasily",
	}
	anyTypeValue = universalMap // cast to an empty interface, just for fun
	result, _ := json.Marshal(anyTypeValue)
	show("Encoded map x to json y, (x, y): ", anyTypeValue, string(result))
	// string({"id":42,"username":"rvasily"});

	show("end of program.")
	/*
	   2023-12-04T11:20:26.550Z: program started ...
	   2023-12-04T11:20:26.550Z: Decoded json x to an empty interface y, (x, y): string([
	                   {"id": 17, "username": "iivan", "phone": 0},
	                   {"id": "17", "address": "none", "company": "Mail.ru"}
	           ]); []interface {}([map[id:17 phone:0 username:iivan] map[address:none company:Mail.ru id:17]]);
	   2023-12-04T11:20:26.550Z: Encoded map x to json y, (x, y): map[string]interface {}(map[id:42 username:rvasily]); string({"id":42,"username":"rvasily"});
	   2023-12-04T11:20:26.550Z: end of program.
	*/
}

func reflect_1() {
	show("program started ...")

	type User struct {
		ID       int
		RealName string `unpack:"-"`
		Login    string
		Flags    int
	}

	var PrintReflect = func(anyVal interface{}) error {
		reflectValue := reflect.ValueOf(anyVal).Elem() // it's a key to all the magic
		show("Given value x has y fields, (x, y): ", anyVal, reflectValue.NumField())

		for i := 0; i < reflectValue.NumField(); i++ {
			valueField := reflectValue.Field(i)
			typeField := reflectValue.Type().Field(i)

			show(
				"Field (name, type, value, tag): ",
				typeField.Name,
				typeField.Type.Kind(),
				valueField,
				typeField.Tag,
			)
		}

		return nil
	}

	userRef := &User{
		ID:       42,
		RealName: "rvasily",
		Flags:    32,
	}

	err := PrintReflect(userRef) // examine the output
	if err != nil {
		panic(err)
	}

	show("end of program.")
	/*
	   2023-12-04T12:24:45.037Z: program started ...
	   2023-12-04T12:24:45.037Z: Given value x has y fields, (x, y): *main.User(&{42 rvasily  32}); int(4);
	   2023-12-04T12:24:45.038Z: Field (name, type, value, tag): string(ID); reflect.Kind(int); reflect.Value(42); reflect.StructTag();
	   2023-12-04T12:24:45.038Z: Field (name, type, value, tag): string(RealName); reflect.Kind(string); reflect.Value(rvasily); reflect.StructTag(unpack:"-");
	   2023-12-04T12:24:45.038Z: Field (name, type, value, tag): string(Login); reflect.Kind(string); reflect.Value(); reflect.StructTag();
	   2023-12-04T12:24:45.038Z: Field (name, type, value, tag): string(Flags); reflect.Kind(int); reflect.Value(32); reflect.StructTag();
	   2023-12-04T12:24:45.038Z: end of program.
	*/
}

func reflect_2() {
	show("program started ...")

	type User struct {
		ID       int
		RealName string `unpack:"-"`
		Login    string
		Flags    int
	}

	var UnpackReflect = func(targetStructRef interface{}, sourceBytes []byte) error {
		// restore int or string fields of given struct, from given bytes

		bytesReader := bytes.NewReader(sourceBytes)
		reflectValue := reflect.ValueOf(targetStructRef).Elem()
		var bytesOrder = binary.LittleEndian // hidden knowledge

		for i := 0; i < reflectValue.NumField(); i++ {
			valueField := reflectValue.Field(i)
			typeField := reflectValue.Type().Field(i)

			if typeField.Tag.Get("unpack") == "-" {
				// skip if tag says so
				continue
			}

			switch typeField.Type.Kind() {

			case reflect.Int:
				var value uint32 // hidden knowlwdge
				binary.Read(bytesReader, bytesOrder, &value)
				valueField.Set(reflect.ValueOf(int(value)))

			case reflect.String:
				var strLen uint32
				binary.Read(bytesReader, bytesOrder, &strLen)
				strBytes := make([]byte, strLen)
				binary.Read(bytesReader, bytesOrder, &strBytes)
				valueField.SetString(string(strBytes))

			default:
				return fmt.Errorf("bad type: %v for field %v", typeField.Type.Kind(), typeField.Name)

			}

		}

		return nil
	}

	/*
		perl -E '$b = pack("L L/a* L", 1_123_456, "v.romanov", 16);
			print map { ord.", "  } split("", $b); '
	*/

	data := []byte{
		128, 36, 17, 0, // int

		9, 0, 0, 0, // str len
		118, 46, 114, 111, 109, 97, 110, 111, 118, // str bytes

		16, 0, 0, 0, // int
	}

	userRef := new(User)
	err := UnpackReflect(userRef, data)
	if err != nil {
		panic(err)
	}
	show("Unpacked struct: ", userRef)

	show("end of program.")
	/*
	   2023-12-04T12:47:43.434Z: program started ...
	   2023-12-04T12:47:43.434Z: Unpacked struct: *main.User(&{1123456  v.romanov 16});
	   2023-12-04T12:47:43.434Z: end of program.
	*/
}

func unpackDemo() {
	// go build gen/* && ./codegen.exe pack/packer.go  pack/marshaller.go

	show("program started ...")

	// codegen program

	var codegenMain = func(inputFileName, outputFileName string) {
		// generate `Unpack` method for marked by `cgen: binpack` struct
		/*
		   // lets generate code for this struct
		   // cgen: binpack
		   type User struct {
		   	ID       int
		   	RealName string `cgen:"-"`
		   	Login    string
		   	Flags    int
		   }
		*/
		type templateData struct {
			FieldName string
		}

		var (
			intTemplate = template.Must(template.New("intTpl").Parse(`
			// {{.FieldName}}
			var {{.FieldName}}Raw uint32
			binary.Read(r, binary.LittleEndian, &{{.FieldName}}Raw)
			in.{{.FieldName}} = int({{.FieldName}}Raw)
		`))

			strTemplate = template.Must(template.New("strTpl").Parse(`
			// {{.FieldName}}
			var {{.FieldName}}LenRaw uint32
			binary.Read(r, binary.LittleEndian, &{{.FieldName}}LenRaw)
			{{.FieldName}}Raw := make([]byte, {{.FieldName}}LenRaw)
			binary.Read(r, binary.LittleEndian, &{{.FieldName}}Raw)
			in.{{.FieldName}} = string({{.FieldName}}Raw)
		`))
		)

		fset := token.NewFileSet()

		// input
		node, err := parser.ParseFile(fset, inputFileName, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}

		// output
		out, err := os.Create(outputFileName)
		if err != nil {
			log.Fatal(err)
		}

		// build output, line by line
		fmt.Fprintln(out, `package `+node.Name.Name)
		fmt.Fprintln(out) // empty line
		fmt.Fprintln(out, `import "encoding/binary"`)
		fmt.Fprintln(out, `import "bytes"`)
		fmt.Fprintln(out) // empty line

		for _, f := range node.Decls {
			g, ok := f.(*ast.GenDecl)
			if !ok {
				fmt.Printf("SKIP %T is not *ast.GenDecl\n", f)
				continue
			}

		SPECS_LOOP:
			for _, spec := range g.Specs {
				currType, ok := spec.(*ast.TypeSpec)
				if !ok {
					fmt.Printf("SKIP %T is not ast.TypeSpec\n", spec)
					continue
				}

				currStruct, ok := currType.Type.(*ast.StructType)
				if !ok {
					fmt.Printf("SKIP %T is not ast.StructType\n", currStruct)
					continue
				}

				if g.Doc == nil {
					fmt.Printf("SKIP struct %#v doesnt have comments\n", currType.Name.Name)
					continue
				}

				needCodegen := false
				for _, comment := range g.Doc.List {
					needCodegen = needCodegen || strings.HasPrefix(comment.Text, "// cgen: binpack")
				}
				if !needCodegen {
					fmt.Printf("SKIP struct %#v doesnt have cgen mark\n", currType.Name.Name)
					continue SPECS_LOOP
				}

				fmt.Printf("process struct %s\n", currType.Name.Name)
				fmt.Printf("\tgenerating Unpack method\n")

				fmt.Fprintln(out, "func (in *"+currType.Name.Name+") Unpack(data []byte) error {")
				fmt.Fprintln(out, "	r := bytes.NewReader(data)")

			FIELDS_LOOP:
				for _, field := range currStruct.Fields.List {

					if field.Tag != nil {
						tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
						if tag.Get("cgen") == "-" {
							continue FIELDS_LOOP
						}
					}

					fieldName := field.Names[0].Name
					fieldType := field.Type.(*ast.Ident).Name

					fmt.Printf("\tgenerating code for field %s.%s\n", currType.Name.Name, fieldName)

					switch fieldType {
					case "int":
						intTemplate.Execute(out, templateData{fieldName})
					case "string":
						strTemplate.Execute(out, templateData{fieldName})
					default:
						log.Fatalln("unsupported", fieldType)
					}
				}

				fmt.Fprintln(out, "	return nil")
				fmt.Fprintln(out, "}") // end of Unpack func
				fmt.Fprintln(out)      // empty line
			}
		}
	}

	// apply codegen to source file with struct declaration

	var srcFileName, trgFileName string
	srcFileName = "./week03/user_struct.go"
	trgFileName = "./week03/user_struct_unpack.go"
	codegenMain(srcFileName, trgFileName)

	// Using generated code in application

	/*
		perl -E '$b = pack("L L/a* L", 1_123_456, "v.romanov", 16);
			print map { ord.", "  } split("", $b); '
	*/
	packedBytes := []byte{
		128, 36, 17, 0, // int

		9, 0, 0, 0, // string len
		118, 46, 114, 111, 109, 97, 110, 111, 118, // string bytes

		16, 0, 0, 0, // int
	}

	user := User{}
	user.Unpack(packedBytes) // apply generated code
	show("Unpacked user: ", user)

	show("end of program.")
	/*
	   2023-12-04T15:33:09.621Z: program started ...
	   SKIP *ast.ImportSpec is not ast.TypeSpec
	   process struct User
	           generating Unpack method
	           generating code for field User.ID
	           generating code for field User.Login
	           generating code for field User.Flags
	   SKIP struct "Avatar" doesnt have comments
	   SKIP *ast.ValueSpec is not ast.TypeSpec
	   SKIP *ast.FuncDecl is not *ast.GenDecl
	   2023-12-04T15:33:09.657Z: Unpacked user: main.User({1123456  v.romanov 16});
	   2023-12-04T15:33:09.657Z: end of program.
	*/
}

func unpack_testBench() {
	show("program started ...")

	show(`
	go test -bench . week03
BenchmarkGenerated-8     3871654               303.0 ns/op
BenchmarkReflect-8       1778101               671.4 ns/op

	go test -bench . -benchmem week03
BenchmarkGenerated-8     3889042               302.2 ns/op           152 B/op          8 allocs/op
BenchmarkReflect-8       1770966               668.3 ns/op           320 B/op         14 allocs/op

BenchmarkGenerated-8     4248238               269.5 ns/op           104 B/op          7 allocs/op
BenchmarkReflect-8       1851026               643.7 ns/op           272 B/op         13 allocs/op

	go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1 unpack_bench_test.go

	go tool pprof main.test.exe cpu.out
	go tool pprof main.test.exe mem.out

	go get github.com/uber/go-torch
	go-torch main.test.exe cpu.out
`)

	show("end of program.")
}

func main() {
	// jsonDemo()
	// struct_tags()
	// dynamicDemo()
	// reflect_1()
	// reflect_2()
	// unpackDemo()
	unpack_testBench()

	// var err = fmt.Errorf("While doing %s: %v", "main", "not implemented")
	// panic(err)
}

func demoTemplate() {
	show("program started ...")
	show("end of program.")
}

func show(msg string, xs ...any) {
	var line = ts() + ": " + msg

	for _, x := range xs {
		// https://pkg.go.dev/fmt
		line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		// line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
}

// ts return current timestamp in RFC3339 with milliseconds
func ts() string {
	/*
		https://pkg.go.dev/time#pkg-constants
		https://stackoverflow.com/questions/35479041/how-to-convert-iso-8601-time-in-golang
	*/
	const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
	return time.Now().UTC().Format(RFC3339Milli)
}
