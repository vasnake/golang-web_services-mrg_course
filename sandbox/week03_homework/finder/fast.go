package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// FastSearch should be (at least) 20% better (resource-wise, faster) than SlowSearch.
//
// Check: `go test -v $module; go test -bench . -benchmem $module`
func FastSearch(out io.Writer) {
	/*
		Task outline:
		For each line from file:
		- transform json => Record, where Record is (list-of-browsers, user-name, user-email)
		- loop over browsers: update set-of-unique-browsers, set MSIE-and-Android flag
		- if MSIE-and-Android: update result (print to out)
		print-to-out: len(set-of-unique-browsers)
	*/

	var browsersSet = make(map[string]struct{}, 128)

	var processBrowsersList = func(browsers []string) bool {
		var isAndroid, isMSIE = false, false

		for _, browser := range browsers {
			// filter condition
			cond1 := strings.Contains(browser, "Android")
			cond2 := strings.Contains(browser, "MSIE")
			if cond1 {
				isAndroid = true
			}
			if cond2 {
				isMSIE = true
			}
			if cond1 || cond2 {
				browsersSet[browser] = struct{}{}
			}
		}

		return isAndroid && isMSIE
	}

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var userInfo = UserInfoRecord{}
	var currentLine = bufio.NewScanner(file)
	var lineNo int = -1 // start from `0`
	fmt.Fprintln(out, "found users:")

	for currentLine.Scan() {
		lineNo += 1
		if err := json.Unmarshal(currentLine.Bytes(), &userInfo); err != nil {
			panic(err)
		}

		if processBrowsersList(userInfo.Browsers) {
			// filtered record, update output
			fmt.Fprintf(out, "[%d] %s <%s>\n", lineNo, userInfo.Name, strings.ReplaceAll(userInfo.Email, "@", " [at] "))
		}
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(browsersSet))
}

type UserInfoRecord struct {
	Browsers []string `json:"browsers"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
}
