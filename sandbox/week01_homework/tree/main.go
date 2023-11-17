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

	var parentPrefix = ""

	// skip it, it is unimportant
	path = strings.TrimSpace(path)
	show("path: ", path)
	var x = filepath.Join(path, "project")
	show("joined dir: ", x)

	entries, err := readdir(path)
	if err != nil {
		return err
	}

	entries = sortByName(entries)

	for idx, entry := range entries {
		isDir, name, size := entry.IsDir(), entry.Name(), entry.Size()
		isLast := (idx + 1) == len(entries)

		prefix, text := formatEntry(name, isDir, size, parentPrefix, isLast)
		fmt.Fprintf(out, "%s%s\n", prefix, text)

		if isDir {
			err = printDirTree(out, path, name, prefix, printFiles)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func printDirTree(out io.Writer, path, dirName, parentPrefix string, printFiles bool) error {
	path = filepath.Join(path, dirName)

	entries, err := readdir(path)
	if err != nil {
		return err
	}

	entries = sortByName(entries)

	for idx, entry := range entries {
		isDir, name, size := entry.IsDir(), entry.Name(), entry.Size()
		isLast := (idx + 1) == len(entries)

		prefix, text := formatEntry(name, isDir, size, parentPrefix, isLast)
		fmt.Fprintf(out, "%s%s\n", prefix, text)

		if isDir {
			err = printDirTree(out, path, name, prefix, printFiles)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func formatEntry(name string, isDir bool, size int64, parentPrefix string, isLast bool) (prefix, text string) {
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

		https://pkg.go.dev/fmt

		got:
		├───project
		├───├───file.txt (19b)
		├───└───gopher.png (70372b)
		├───static
		├───├───a_lorem
		├───├───├───dolor.txt (empty)
		├───├───├───gopher.png (70372b)
		├───├───└───ipsum
		├───├───└───└───gopher.png (70372b)
		├───├───css
		├───├───└───body.css (28b)
		├───├───empty.txt (empty)
		├───├───html
		├───├───└───index.html (57b)
		├───├───js
		├───├───└───site.js (10b)
		├───└───z_lorem
		├───└───├───dolor.txt (empty)
		├───└───├───gopher.png (70372b)
		├───└───└───ipsum
		├───└───└───└───gopher.png (70372b)
		├───zline
		├───├───empty.txt (empty)
		├───└───lorem
		├───└───├───dolor.txt (empty)
		├───└───├───gopher.png (70372b)
		├───└───└───ipsum
		├───└───└───└───gopher.png (70372b)
		└───zzfile.txt (empty)

		expected:
		├───project
		│	├───file.txt (19b)
		│	└───gopher.png (70372b)
		├───static
		│	├───a_lorem
		│	│	├───dolor.txt (empty)
		│	│	├───gopher.png (70372b)
		│	│	└───ipsum
		│	│		└───gopher.png (70372b)
		│	├───css
		│	│	└───body.css (28b)
		│	├───empty.txt (empty)
		│	├───html
		│	│	└───index.html (57b)
		│	├───js
		│	│	└───site.js (10b)
		│	└───z_lorem
		│		├───dolor.txt (empty)
		│		├───gopher.png (70372b)
		│		└───ipsum
		│			└───gopher.png (70372b)
		├───zline
		│	├───empty.txt (empty)
		│	└───lorem
		│		├───dolor.txt (empty)
		│		├───gopher.png (70372b)
		│		└───ipsum
		│			└───gopher.png (70372b)
		└───zzfile.txt (empty)

	*/

	var namePrefix = "├───"
	if isLast {
		namePrefix = "└───"
	}

	namePrefix = parentPrefix + namePrefix

	var namePostfix = ""
	if !isDir {
		var sizeText = "empty"
		if size > 0 {
			sizeText = fmt.Sprintf("%db", size)
		}

		namePostfix = fmt.Sprintf(" (%s)", sizeText)
	}

	prefix = namePrefix
	text = fmt.Sprintf("%s%s", name, namePostfix)
	return
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

	entries, err := file.ReadDir(0) // DirEntry nave no size attribute, we will need to read from disk again, it is stupid
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
