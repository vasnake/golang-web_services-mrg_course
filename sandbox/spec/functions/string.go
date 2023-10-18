package functions

import (
	"bytes"
	"fmt"
	"unicode/utf8"
	"unsafe"
	// "strconv"
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

type UnsignedIntegrals interface {
	~byte | ~uint16 | ~uint32 | ~uint64
}
type SignedIntegrals interface {
	~int8 | ~int16 | ~int32 | ~int64
}
type Integrals interface {
	UnsignedIntegrals | SignedIntegrals
}

func IntBits[T Integrals](x T) string {
	// var res = strconv.FormatInt(i, base)
	// res = fmt.Sprintf("%64b", i)

	// switch it := any(x).(type) {
	// case int:
	// 	numBytes = unsafe.Sizeof(int)
	// default:
	// 	fmt.Printf("I don't know about type %T!\n", v)
	// }

	const maxSize = 8
	size := int(unsafe.Sizeof(x))
	mem := unsafe.Pointer(&x)
	inputBytes := *(*[maxSize]byte)(mem)
	byteSlice := make([]byte, size)

	for i := 0; i < int(size); i++ {
		byteSlice[i] = inputBytes[size-i-1]
	}

	return string(fmtBits(byteSlice))
}

func fmtBits(data []byte) []byte {
	var buf bytes.Buffer
	for _, b := range data {
		fmt.Fprintf(&buf, "%08b ", b)
	}
	buf.Truncate(buf.Len() - 1) // To remove extra space
	return buf.Bytes()
}
