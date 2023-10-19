package functions

import (
	"bytes"
	"fmt"
	"unsafe"
	// "unicode/utf8"
	// "strconv"
)

// type constraint
type IntegralType interface {
	UnsignedIntegralType | SignedIntegralType
}
type UnsignedIntegralType interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint
}
type SignedIntegralType interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~int
}

// IntBits function returns binary representation for input,
// bytes separated with space symbol.
func IntBits[T IntegralType](x T) string {
	// var res = strconv.FormatInt(i, base)
	// res = fmt.Sprintf("%64b", i)

	// switch it := any(x).(type) {
	// case int:
	// 	numBytes = unsafe.Sizeof(int)
	// default:
	// 	fmt.Printf("I don't know about type %T!\n", v)
	// }

	const maxSize = 8 // bytes
	size := int(unsafe.Sizeof(x))
	mem := unsafe.Pointer(&x)
	inputBytes := *(*[maxSize]byte)(mem) // array 8 byte
	byteSlice := make([]byte, size)

	for i := 0; i < size; i++ {
		byteSlice[i] = inputBytes[size-i-1]
	}

	return string(formatBytes(byteSlice))
}

func formatBytes(data []byte) []byte {
	var buf bytes.Buffer
	for _, b := range data {
		fmt.Fprintf(&buf, "%08b ", b)
	}
	buf.Truncate(buf.Len() - 1) // To remove extra space
	return buf.Bytes()
}

// DotProduct generic function computes dot product for two vectors.
// Vectors should be of the eqal length
func DotProduct[F ~float32 | ~float64](v1, v2 []F) F {
	// the product `x * y` and the addition `s += x * y` are computed with
	// float32 or float64 precision, respectively, depending on the type argument for `F`
	var s F
	for i, x := range v1 {
		y := v2[i]
		s += x * y
	}
	return s
}

// TwoValuesToArray function returns given parameters as array.
func TwoValuesToArray(a, b any) [2]any {
	return [2]any{a, b}
}

// MakeChannel function returns new bidirectional channel, created with options:
// buffer length, list of values put in channel, close channel before returning.
func MakeChannel[T any](bufferLength int, doClose bool, xs ...T) chan T {
	var ch = make(chan T, bufferLength)
	for _, x := range xs {
		ch <- x
	}
	if doClose {
		close(ch)
	}
	return ch
}
