package functions

// RunesCount function returns number of runes in string
var RunesCount = func(str string) int {
	var runes = []rune(str) // allocation?
	return len(runes)
}
