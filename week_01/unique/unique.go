package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func sortedInputUnique(input io.Reader, output io.Writer) error {
	in := bufio.NewScanner(input)
	var prev string
	for in.Scan() {
		txt := in.Text()
		fmt.Printf("input line: `%s`", txt)

		if txt == prev {
			fmt.Printf(" was seen already\n")
			continue
		}

		if txt < prev {
			return fmt.Errorf("input is not sorted")
		}

		prev = txt
		fmt.Printf(" is unique\n")

		fmt.Fprintln(output, txt)
	}
	return nil
}

func main() {
	// GO111MODULE=off cat data_map.txt | sort | go run uniq || exit
	err := sortedInputUnique(os.Stdin, os.Stdout)
	if err != nil {
		panic(err.Error())
	}
}
