package main

/*
Web services on Go, week 1, homework, `tree` program.

mkdir -p week01_homework/tree
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
go run tree ./tree/testdata
cd tree && docker build -t mailgo_hw1 .

https://en.wikipedia.org/wiki/Tree_(command)
https://mama.indstate.edu/users/ice/tree/
https://stackoverflow.com/questions/32151776/visualize-tree-in-bash-like-the-output-of-unix-tree

*/

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

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
	- output all dir items in sorted order, w/o distinction file/dir
	- last element prefix is `└───`
	- other elements prefix is `├───`
	- nested elements aligned with one tab `	` for each level
*/

const (
	EOL   = "\n"
	TRUNK = "│"

	BRANCHING_TRUNK        = "├───"
	BRANCHING_TRUNK_SYMBOL = "├"

	LAST_BRANCH        = "└───"
	LAST_BRANCH_SYMBOL = "└"

	TRUNC_TAB = "│\t"
	LAST_TAB  = "\t"
)

func main() {
	// This code is given, I don't think I should touch it
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
	// Function to implement, signature is given, don't touch it.

	var pathSubdir, parentPrefix = "", ""
	return printDirTree(out, path, pathSubdir, parentPrefix, printFiles)
}

func printDirTree(out io.Writer, path, dirName, parentPrefix string, printFiles bool) error {
	if dirName != "" { // On first call dirName = ""
		path = filepath.Join(path, dirName)
	}

	entries, err := readdir(path)
	if err != nil {
		return err
	}

	if !printFiles {
		entries = dropFiles(entries)
	}

	entries = sortByName(entries)

	for idx, entry := range entries {
		isDir, name, size := entry.IsDir(), entry.Name(), entry.Size()
		isLast := (idx + 1) == len(entries)

		prefix, text := formatEntry(name, isDir, size, parentPrefix, isLast)
		fmt.Fprintf(out, "%s%s%s", prefix, text, EOL)

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
		https://pkg.go.dev/fmt

		got:
		├───project
		│       ├───file.txt (19b)
		│       └───gopher.png (70372b)
		├───static
		│       ├───a_lorem
		│       │       ├───dolor.txt (empty)
		│       │       ├───gopher.png (70372b)
		│       │       └───ipsum
		│       │               └───gopher.png (70372b)
		│       ├───css
		│       │       └───body.css (28b)
		│       ├───empty.txt (empty)
		│       ├───html
		│       │       └───index.html (57b)
		│       ├───js
		│       │       └───site.js (10b)
		│       └───z_lorem
		│               ├───dolor.txt (empty)
		│               ├───gopher.png (70372b)
		│               └───ipsum
		│                       └───gopher.png (70372b)
		├───zline
		│       ├───empty.txt (empty)
		│       └───lorem
		│               ├───dolor.txt (empty)
		│               ├───gopher.png (70372b)
		│               └───ipsum
		│                       └───gopher.png (70372b)
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

	var namePrefix = BRANCHING_TRUNK
	if isLast {
		namePrefix = LAST_BRANCH
	}

	if endsWith(parentPrefix, BRANCHING_TRUNK) {
		parentPrefix = replaceTail(parentPrefix, len(BRANCHING_TRUNK), TRUNC_TAB)
	} else if endsWith(parentPrefix, LAST_BRANCH) {
		parentPrefix = replaceTail(parentPrefix, len(LAST_BRANCH), LAST_TAB)
	}

	namePrefix = parentPrefix + namePrefix

	// result
	prefix = namePrefix
	text = formatName(name, isDir, size)
	return
}

func endsWith(text, subtext string) bool {
	var start = len(text) - len(subtext)
	return start >= 0 && text[start:] == subtext
}

func replaceTail(text string, tailLen int, trg string) string {
	var start = len(text) - tailLen
	if start >= 0 {
		return text[0:start] + trg
	}
	return trg
}

func formatName(name string, isDir bool, size int64) string {
	/*
		Result examples
		- `lorem`: directory
		- `dolor.txt (empty)`: empty file
		- `gopher.png (70372b)`: regular file
	*/
	var suffix = ""
	if !isDir {
		var sizeText = "empty"
		if size > 0 {
			sizeText = fmt.Sprintf("%db", size)
		}

		suffix = fmt.Sprintf(" (%s)", sizeText)
	}

	return fmt.Sprintf("%s%s", name, suffix)
}

func dropFiles(entries []fs.FileInfo) []fs.FileInfo {
	var dirsCount uint = 0
	for _, entry := range entries {
		if entry.IsDir() {
			dirsCount += 1
		}
	}

	var dirs = make([]fs.FileInfo, 0, dirsCount)
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry)
		}
	}

	return dirs
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

// readDir is a fast (relatively) but useless (here) function. It is here only for educational purposes.
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
