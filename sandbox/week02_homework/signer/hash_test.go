package main

import (
	"fmt"
	"testing"
	"time"
)

type testDataTuple struct {
	input, expected string
}

func TestMultiHash(t *testing.T) {
	var testData = []testDataTuple{
		{
			input:    "4108050209~502633748",
			expected: "29568666068035183841425683795340791879727309630931025356555",
		},
		{
			input:    "2212294583~709660146",
			expected: "4958044192186797981418233587017209679042592862002427381542",
		},
	}

	var maxDuration = time.Millisecond * 1500

	var errors = testHash(computeMultiHash, testData, maxDuration, t)

	for _, error := range errors {
		t.Errorf("%s", error)
	}
}

func TestSingleHash(t *testing.T) {
	var testData = []testDataTuple{
		{
			input:    "0",
			expected: "4108050209~502633748",
		},
		{
			input:    "1",
			expected: "2212294583~709660146",
		},
	}

	var maxDuration = time.Millisecond * 1500

	var errors = testHash(computeSingleHash, testData, maxDuration, t)

	for _, error := range errors {
		t.Errorf("%s", error)
	}
}

func TestCombineResults(t *testing.T) {
	var testData = []string{
		"4958044192186797981418233587017209679042592862002427381542",
		"29568666068035183841425683795340791879727309630931025356555",
	}
	var expected = "29568666068035183841425683795340791879727309630931025356555_4958044192186797981418233587017209679042592862002427381542"

	var result = computeCombineResults(testData)

	if result != expected {
		t.Errorf("Unexpected result, expected %[2]v, got %[1]v", result, expected)
	}
}

func testHash(hashFunc func(string) string, testData []testDataTuple, maxDuration time.Duration, t *testing.T) []string {
	var result = make([]string, 0, len(testData))

	for idx, rec := range testData {
		var start = time.Now()
		var hash = hashFunc(rec.input)
		var duration = time.Since(start)

		if hash != rec.expected {
			result = append(result, fmt.Sprintf("Unexpected result, expected %[2]v, got %[1]v, recNo %[3]d", hash, rec.expected, idx))
			// t.Errorf("Not expected result, expected %[2]v, got %[1]v, recNo %[3]d", hash, rec.expected, idx)
		}

		if duration > maxDuration {
			result = append(result, fmt.Sprintf("Computation took too long, %s > %s", duration, maxDuration))
			// t.Errorf("Computation took too long, %s > %s", duration, maxDuration)
		}
	}

	return result
}
