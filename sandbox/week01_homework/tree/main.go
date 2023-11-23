package main

/*
Course `Web services on Go`, week 1, homework, `tree` program.
See: week_01\materials.zip\week_1\99_hw\tree

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
	EOL             = "\n"
	BRANCHING_TRUNK = "├───"
	LAST_BRANCH     = "└───"
	TRUNC_TAB       = "│\t"
	LAST_TAB        = "\t"
	EMPTY_FILE      = "empty"

	USE_RECURSION_ENV_KEY = "RECURSIVE_TREE"
	USE_RECURSION_ENV_VAL = "YES"
)

func main() {
	// This code is given, I don't think I should touch it
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage: go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

// dirTree: top-level `tree` program implementation.
// It is write `path` dir listing to `out`. If `prinFiles` is set, files is listed along with directories.
func dirTree(out io.Writer, path string, printFiles bool) error {
	// Function to implement, signature is given, don't touch it.

	var useRecursionEvv = os.Getenv(USE_RECURSION_ENV_KEY)
	show("Recursion option: (key, expectedVal, actualVal): ", USE_RECURSION_ENV_KEY, USE_RECURSION_ENV_VAL, useRecursionEvv)
	var useRecursion = strings.ToUpper(useRecursionEvv) == strings.ToUpper(USE_RECURSION_ENV_VAL)

	var parentPrefix = ""
	if useRecursion {
		show("Print tree using recursion ...")
		return printDirTreeRecur(out, path, parentPrefix, printFiles)
	} else {
		show("Print tree using stack, no recursion ...")
		return printDirTreeStack(out, path, printFiles)
	}
}

func printDirTreeStack(out io.Writer, path string, printFiles bool) error {

	type dirEntry struct {
		name         string
		isDir        bool
		size         int64
		isLast       bool
		parentPrefix string
		parentPath   string
	}

	var getEntries = func(dirPath, parentPrefix string) ([]dirEntry, error) {
		entries, err := readdir(dirPath)
		if err != nil {
			return nil, err
		}

		if !printFiles {
			entries = dropFiles(entries)
		}

		entries = sortByName(entries)

		var myEntries = make([]dirEntry, len(entries))
		for idx := range entries {
			myEntries[idx] = dirEntry{
				name:         entries[idx].Name(),
				isDir:        entries[idx].IsDir(),
				size:         entries[idx].Size(),
				isLast:       idx == (len(entries) - 1),
				parentPrefix: parentPrefix,
				parentPath:   dirPath,
			}
		}

		return myEntries, nil
	}

	entries, err := getEntries(path, "") // root prefix
	if err != nil {
		return err
	}

	for len(entries) > 0 {
		// pop item from stack
		var entry = entries[0]
		entries = entries[1:]

		prefix, text := formatEntry(entry.name, entry.isDir, entry.size, entry.parentPrefix, entry.isLast)
		fmt.Fprint(out, prefix+text+EOL)

		if entry.isDir {
			// push new items to stack
			newEntries, err := getEntries(filepath.Join(entry.parentPath, entry.name), prefix)
			if err != nil {
				return err
			}
			entries = append(newEntries, entries...)
		}
	}

	return nil // no errors
}

// printDirTreeRecur: recursive implementation of a `tree` program. Parameters:
// `out`: result, where to write directory tree. `path`: directory to process.
// `parentPrefix`: text representing tree leaf or branch, w/o file actual data.
// `printFiles`: define option for writing files (not directories) to result.
func printDirTreeRecur(out io.Writer, path, parentPrefix string, printFiles bool) error {
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
		// fmt.Fprintf(out, "%s%s%s", prefix, text, EOL) // could be a little faster
		fmt.Fprint(out, prefix+text+EOL)

		if isDir {
			err = printDirTreeRecur(out, filepath.Join(path, name), prefix, printFiles)
			if err != nil {
				return err
			}
		}
	}

	return nil // no errors
}

func formatEntry(name string, isDir bool, size int64, parentPrefix string, isLast bool) (prefix, text string) {
	/*
		https://pkg.go.dev/fmt

		Complete text, set of entries example (expected):
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

	if endsWith(parentPrefix, BRANCHING_TRUNK) {
		prefix = replaceTail(parentPrefix, len(BRANCHING_TRUNK), TRUNC_TAB)
	} else if endsWith(parentPrefix, LAST_BRANCH) {
		prefix = replaceTail(parentPrefix, len(LAST_BRANCH), LAST_TAB)
	} else {
		prefix = parentPrefix
	}

	if isLast {
		prefix += LAST_BRANCH
	} else {
		prefix += BRANCHING_TRUNK
	}

	text = formatName(name, isDir, size)
	return // prefix, text
}

func endsWith(text, subtext string) bool {
	var start = len(text) - len(subtext)
	return start >= 0 && text[start:] == subtext
}

func replaceTail(text string, tailLen int, newTail string) (result string) {
	var start = len(text) - tailLen
	if start >= 0 {
		result = text[0:start] + newTail
	} else {
		result = newTail
	}
	return
}

func formatName(name string, isDir bool, size int64) string {
	/*
		Result examples
		- `lorem`: directory
		- `dolor.txt (empty)`: empty file
		- `gopher.png (70372b)`: regular file, size in bytes
	*/
	var suffix = "" // if `name` is a directory
	if !isDir {
		var sizeText = EMPTY_FILE // if file is empty
		if size > 0 {
			sizeText = fmt.Sprintf("%db", size)
		}

		suffix = fmt.Sprintf(" (%s)", sizeText)
	}

	return name + suffix
}

func dropFiles(entries []fs.FileInfo) []fs.FileInfo {
	// I think two slice enumerations should be more effective than x memory reallocations in case when result slice size is unknown
	var dirsCount uint = 0 // result slice size
	for idx := range entries {
		if entries[idx].IsDir() {
			dirsCount += 1
		}
	}

	var dirs = make([]fs.FileInfo, dirsCount)
	var dirsIdx uint = 0
	for entriesIdx := range entries {
		if entries[entriesIdx].IsDir() {
			dirs[dirsIdx] = entries[entriesIdx]
			dirsIdx += 1
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

	entries, err := file.Readdir(0) // 0: read all entries
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// readDir is fast (relatively) but useless (here) function. It is here only for educational purposes.
// DirEntry nave no `size` attribute, we will need to read from disk again, it is stupid.
func readDir(path string) ([]fs.DirEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	entries, err := file.ReadDir(0)
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
