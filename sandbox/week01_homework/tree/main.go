package main

/*
pushd week01_homework/tree
go mod init tree
go mod tidy
pushd ..
go work init
go work use ./tree/
go vet tree
gofmt -w tree
go test -v tree
go run tree . -f
cd tree && docker build -t mailgo_hw1 .
*/

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	/*
		Example output:

		├───project
		│	└───gopher.png (70372b)
		├───static
		│	├───a_lorem
		│	│	├───dolor.txt (empty)
		│	├───css
		│	│	└───body.css (28b)
		...
		│			└───gopher.png (70372b)

		- path should point to a directory,
		- directory items must be sorted,
		- output all dir items in sorted order, w/o distinction file/dir
		- last element prefix is `└───`
		- other elements prefix is `├───`
		- nested elements aligned with one tab `	` for each level
	*/

	path = strings.TrimSpace(path)
	show("path: ", path)
	var x = filepath.Join(path, "project")
	show("joined dir: ", x)

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	items, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return err
	}

	var dirs []string
	var files []string
	for _, item := range items {
		// show("got item: ", item)
		if item.IsDir() {
			dirs = append(dirs, item.Name())
		} else {
			files = append(files, item.Name())
		}
	}
	sort.Strings(dirs)
	slices.Sort(dirs)

	show("dirs: ", dirs)
	show("files: ", files)

	/*
				joined dir: string(project);
				got item: *os.fileStat(
					&{tree 4096 2147484159
						{331655100 63835653016 0x53bf00}
						{72 2251799814208938 1 16895 1000 1000 0 0 4096 4096 0
							{1700062361 140356500}
							{1700056216 331655100}
							{1700056216 331655100}
							[0 0 0]
						}
					}
				);
		dirs: []string([project static zline]);
		files: []string([]);
	*/

	var line string = "├───project"
	fmt.Fprintln(out, line)
	return nil
}

func show(msg string, xs ...any) {
	var line string = msg
	for _, x := range xs {
		line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		// line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
}
