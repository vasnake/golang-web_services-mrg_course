package main

/*
# go test -v $module_name
=== RUN   TestSortedInput
input line: `1` is unique
input line: `2` is unique
input line: `2` was seen already
input line: `3` is unique
--- PASS: TestSortedInput (0.00s)
=== RUN   TestUnsortedInput
input line: `1` is unique
input line: `2` is unique
input line: `1`--- PASS: TestUnsortedInput (0.00s)
PASS
ok      week01  0.005s
*/

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

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
