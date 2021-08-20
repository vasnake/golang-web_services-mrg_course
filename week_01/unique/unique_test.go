package main

// GO111MODULE=off go test -v ./unique || exit

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

// file $name_test.go, n.b. `test` suffix

func TestSortedInput(t *testing.T) {
	// test function should start with `Test` prefix

	var inputData = `1
2
2
3`
	var expected = `1
2
3
` // n.b. \n in EOF

	in := bufio.NewReader(strings.NewReader(inputData))
	out := new(bytes.Buffer)

	err := sortedInputUnique(in, out)

	if err != nil || out.String() != expected {
		t.Errorf("TestSortedUnique failed, error: %#v; result: %#v", err, out.String())
	}

}

func TestUnsortedInput(t *testing.T) {

	var inputData = `1
2
1`

	in := bufio.NewReader(strings.NewReader(inputData))
	out := new(bytes.Buffer)

	err := sortedInputUnique(in, out) // should return error

	if err == nil {
		t.Errorf("TestUnsortedInput failed, error: %#v; result: %#v", err, out.String())
	}

}
