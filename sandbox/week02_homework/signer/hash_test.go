package main

import (
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

	var maxDuration = time.Second * 2

	for idx, rec := range testData {
		var start = time.Now()
		var hash = computeMultiHash(rec.input)
		var duration = time.Since(start)

		if hash != rec.expected {
			t.Errorf("Not expected result, expected %[2]v, got %[1]v, recNo %[3]d", hash, rec.expected, idx)
		}

		if duration > maxDuration {
			t.Errorf("Computation took too long, %s > %s", duration, maxDuration)
		}
	}
}
