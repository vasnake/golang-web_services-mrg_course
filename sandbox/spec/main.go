package main

import "fmt"

func main() {
	show("Sveiki!")
	integers()
	floats()
	imaginary()
	rune()
}

func integers() {
	show("Int ...")
	show("Decimal: ", 0, 123, 123_456_789)
	show("Binary: ", 0b0, 0b11, 0b_010, 0b_0_0_1)
	show("Octal: ", 0123, 0_123, 0o123, 0o76_54_32_10) // N.B. optional `o`, unlike other bases
	show("Hex: ", 0x0A, 0x0a, 0x_0A, 0x_1234_5678_9abc_ef00, 0x_Bad_Face, 0xBEEF, 0x1E-2 /*0x1e - 2*/)
}

func floats() {
	show("Float ...")
	// N.B. `0` as prefix will be ignored
	show("Decimal: ", 1., 01.0e+0, 1.e+0, 1.0, 1e0, .1, .1e1, 1.1e-1, 012_345_6.7_89e-01)
	// Hexadecimal floating-point constants make it easy for the compiler to reproduce the exact value.
	// 0x1.fp-2 is (1 + 15/16)•(2^-2) = .484375
	show("Hex: ", 0xA.Ap0, 0xA.Fp+1, 0xAp1, 0x.Ap1, 0x_Ap-0_1, 0xAp-1)
}

func imaginary() {
	show("Imaginary ...")
	show("Decimal: ", 0i, 1i, 987i)
	show("Int: ", 0b11i, 0o123i, 0xFi, 123i)
	show("Float: ", 1.i, 01.0e0i, 1.1e-1i)
}

func rune() {
	// Unicode code point, int32
	show("Rune ...")
	show("Characters: ", 'x', '\n', '0', 'a', '\'')
	show("Special chars: ", '\a', '\b', '\f', '\n', '\r', '\t', '\v', '\\', '\'', '"')
	show("Unicode: ", '本', 'Я', '\uABCD', '\U00101234') // N.B. little `u` and big `U`
	show("Byte: ", '\377', '\xFF')
}

func show(msg string, xs ...any) {
	var line string = msg
	for _, x := range xs {
		line += fmt.Sprintf("%T(%v); ", x, x)
		// line += fmt.Sprintf("%#v; ", x)
	}
	fmt.Println(line)
}
