package main

import (
	"bufio"
	"fmt"
	"os"
)

func naiveUnique() {
	// GO111MODULE=off cat data_map.txt | go run uniq || exit
	in := bufio.NewScanner(os.Stdin)
	seenAlready := make(map[string]bool)

	for in.Scan() {
		txt := in.Text()
		fmt.Printf("input line: `%s`", txt)

		if _, found := seenAlready[txt]; found {
			fmt.Printf("\n")
			continue
		}

		seenAlready[txt] = true
		fmt.Printf(" is unique \n")
	}
}

func sortedInputUnique() {
	// GO111MODULE=off cat data_map.txt | sort | go run uniq || exit
	in := bufio.NewScanner(os.Stdin)
	var prev string
	for in.Scan() {
		txt := in.Text()
		fmt.Printf("input line: `%s`", txt)

		if txt == prev {
			fmt.Printf("\n")
			continue
		}

		if txt < prev {
			panic("input is not sorted")
		}

		prev = txt
		fmt.Printf(" is unique \n")
	}
}

func main() {
	//naiveUnique()
	sortedInputUnique()
}
