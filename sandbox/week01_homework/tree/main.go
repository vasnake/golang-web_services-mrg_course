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
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	// "slices"
	// "sort"
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

	var treeLevel uint = 0

	// skip it, it is unimportant
	path = strings.TrimSpace(path)
	show("path: ", path)
	var x = filepath.Join(path, "project")
	show("joined dir: ", x)

	// open (and close) dir
	// f, err := os.Open(path)
	// if err != nil {
	// 	return err
	// }

	// get content from dir
	// items, err := f.Readdir(-1)
	// f.Close()
	// if err != nil {
	// 	return err
	// }

	entries, err := readdir(path)
	if err != nil {
		return err
	}

	entries = sortByName(entries)

	// var dirs []string
	// var files []string

	for idx, entry := range entries {
		isDir, name, size := entry.IsDir(), entry.Name(), entry.Size()
		isLast := (idx + 1) == len(entries)
		show("file (index, isDir, name, sizeBytes): ", idx, isDir, name, size)
		outLine := formatEntry(name, isDir, size, treeLevel, isLast)
		fmt.Fprintln(out, outLine)
	}

	// sort.Strings(dirs)
	// slices.Sort(dirs)
	// show("dirs: ", dirs)
	// show("files: ", files)

	/*
		joined dir: string(project);
		dirs: []string([project static zline]);
		files: []string([]);
	*/

	return nil
}

func formatEntry(name string, isDir bool, size int64, treeLevel uint, isLast bool) string {
	defer panic("not implemented yet")
	return name
}

func readdir(path string) ([]fs.FileInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	entries, err := file.Readdir(0)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func readDir(path string) ([]fs.DirEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	entries, err := file.ReadDir(0) // DirEntry nave no size, need to read from disk again, it is stupid
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func sortByName(entries []fs.FileInfo) []fs.FileInfo {
	// cmp(a, b) should return
	// a negative number when a < b,
	// a positive number when a > b
	// and zero when a == b.
	slices.SortFunc(entries, compareByName)
	return entries
}

var compareByName = func(a, b fs.FileInfo) int {
	return strings.Compare(a.Name(), b.Name())
}

func show(msg string, xs ...any) {
	var line string = msg
	for _, x := range xs {
		line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		// line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
}
