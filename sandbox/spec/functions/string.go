package functions

import (
	"fmt"
	"unicode/utf8"
)

// RuneCount function returns number of runes in string
var RuneCount = func(str string) (int, error) {
	rcs := utf8.RuneCountInString(str) // NO allocation
	rcb := utf8.RuneCount([]byte(str)) // allocation?

	runes := []rune(str) // allocation?
	rcr := len(runes)

	if rcb == rcs && rcs == rcr {
		return rcs, nil
	}
	return 0, fmt.Errorf("Invalid string: %#v", str)
}
