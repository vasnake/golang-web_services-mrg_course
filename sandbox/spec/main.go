package main

import (
	"fmt"
)

func main() {
	show("Sveiki!")

	integers()
	floats()
	imaginary()
	runeLiterals()
	stringLiterals()
	constants()
	variables()
	types()
}

func integers() {
	show("An integer literal is a sequence of digits representing an integer constant")
	// For readability, an underscore character _ may appear after a base prefix or between successive digits
	show("Decimal: ", 0, 123, 123_456_789)

	// An optional prefix sets a non-decimal base
	show("Binary: ", 0b0, 0b11, 0b_010, 0b_0_0_1)
	show("Octal: ", 0123, 0_123, 0o123, 0o76_54_32_10) // N.B. optional `o`, unlike other bases
	show("Hex: ", 0x0A, 0x0a, 0x_0A, 0x_1234_5678_9abc_ef00, 0x_Bad_Face, 0xBEEF, 0x1E-2 /*0x1e - 2*/)

	// _42         // an identifier, not an integer literal
	// 42_         // invalid: _ must separate successive digits
	// 4__2        // invalid: only one _ at a time
	// 0_xBadFace  // invalid: _ must separate successive digits
}

func floats() {
	show("A floating-point literal is a decimal or hexadecimal representation of a floating-point constant")

	// One of the integer part or the fractional part may be elided;
	// one of the decimal point or the exponent part may be elided.
	// N.B. `0` as prefix will be ignored
	show("Decimal: ", 1., 01.0e+0, 1.e+0, 1.0, 1e0, .1, .1e1, 1.1e-1, 012_345_6.7_89e-01)

	// Hexadecimal floating-point constants make it easy for the compiler to reproduce the exact value.
	// 0x1.fp-2 is (1 + 15/16)•(2^-2) = .484375
	// For readability, an underscore character _ may appear after a base prefix or between successive digits
	show("Hex: ", 0xA.Ap0, 0xA.Fp+1, 0xAp1, 0x.Ap1, 0x_Ap-0_1, 0xAp-1)
	// consists of a 0x or 0X prefix,
	// an integer part (hexadecimal digits),
	// a radix point,
	// a fractional part (hexadecimal digits),
	// and an exponent part (p or P followed by an optional sign and decimal digits).
	// One of the integer part or the fractional part may be elided; the radix point may be elided as well, but the exponent part is required

	// 0x15e-2      // == 0x15e - 2 (integer subtraction)
	// 0x.p1        // invalid: mantissa has no digits
	// 1p-2         // invalid: p exponent requires hexadecimal mantissa
	// 0x1.5e-2     // invalid: hexadecimal mantissa requires p exponent
	// 1_.5         // invalid: _ must separate successive digits
	// 1._5         // invalid: _ must separate successive digits
	// 1.5_e1       // invalid: _ must separate successive digits
	// 1.5e_1       // invalid: _ must separate successive digits
	// 1.5e1_       // invalid: _ must separate successive digits
}

func imaginary() {
	show("An imaginary literal represents the imaginary part of a complex constant")
	// It consists of an integer or floating-point literal followed by the lowercase letter i
	// The value of an imaginary literal is the value of the respective integer or floating-point literal multiplied by the imaginary unit i
	// integer part consisting entirely of decimal digits ... is considered a decimal integer, even if it starts with a leading 0

	show("Decimal: ", 0i, 1i, 987i)
	show("Int: ", 0b11i, 0o123i, 0xFi, 123i)
	show("Float: ", 1.i, 01.0e0i, 1.1e-1i)

	// 0123i  // == 123i for backward-compatibility
	// 0o123i // == 0o123 * 1i == 83i
}

func runeLiterals() {
	// Unicode code point, int32
	show("A rune literal represents a rune constant, an integer value identifying a Unicode code point")
	// A rune literal is expressed as one or more characters enclosed in single quotes
	// A single quoted character represents the Unicode value of the character itself,
	// while multi-character sequences beginning with a backslash encode values in various formats

	show("Characters: ", 'x', '\n', '0', 'a', '\'')

	// There are four ways to represent the integer value as a numeric constant:
	// \x followed by exactly two hexadecimal digits;
	// \u followed by exactly four hexadecimal digits;
	// \U followed by exactly eight hexadecimal digits,
	// and a plain backslash \ followed by exactly three octal digits

	show("Special chars: ", '\a', '\b', '\f', '\n', '\r', '\t', '\v', '\\', '\'', '"')
	show("Unicode: ", '本', 'Я', '\uABCD', '\U00101234') // N.B. little `u` and big `U`
	show("Byte: ", '\377', '\xFF')

	// 'aa'         // illegal: too many characters
	// '\k'         // illegal: k is not recognized after a backslash
	// '\xa'        // illegal: too few hexadecimal digits
	// '\0'         // illegal: too few octal digits
	// '\400'       // illegal: octal value over 255
	// '\uDFFF'     // illegal: surrogate half
	// '\U00110000' // illegal: invalid Unicode code point
}

func stringLiterals() {
	show("A string literal represents a string constant obtained from concatenating a sequence of characters")
	// There are two forms: raw string literals and interpreted string literals

	// string composed of the uninterpreted (implicitly UTF-8-encoded) characters
	// Raw string literals are character sequences between back quotes, as in `foo`
	show(
		"Raw: ",
		`foo`,
		`foo\r\nbar`,
		`new line and 2 tabs:
		`,
		` " ' \ `,
	)

	// backslash escapes interpreted as they are in rune literals
	// Interpreted string literals are character sequences between double quotes, as in "bar"
	// The three-digit octal (\nnn) and two-digit hexadecimal (\xnn) escapes represent individual bytes of the resulting string
	show(
		"Interpreted: ",
		"new line and 2 tabs:\n\t\t",
		" \" ' \\ ",
		"\xFF \377",
		"ÿ \u00FF \U000000FF \xc3\xbf", // N.B. little `u` and big `U`
	)

	// These examples all represent the same string:
	// "日本語"                                 // UTF-8 input text
	// `日本語`                                 // UTF-8 input text as a raw literal
	// "\u65e5\u672c\u8a9e"                    // the explicit Unicode code points
	// "\U000065e5\U0000672c\U00008a9e"        // the explicit Unicode code points
	// "\xe6\x97\xa5\xe6\x9c\xac\xe8\xaa\x9e"  // the explicit UTF-8 bytes
}

func constants() {
	show("Constants, numeric and string")
	// There are boolean constants, rune constants, integer constants, floating-point constants, complex constants, and string constants.
	// Rune, integer, floating-point, and complex constants are collectively called numeric constants.

	// Numeric constants represent exact values of arbitrary precision and do not overflow.
	// Consequently, there are no constants denoting the IEEE-754 negative zero, infinity, and not-a-number values.
	// Implementation restriction: Although numeric constants have arbitrary precision in the language,
	// a compiler may implement them using an internal representation with limited precision

	// Constants may be typed or untyped.
	// Literal constants, true, false, iota, and certain constant expressions ... are untyped
	// An untyped constant has a default type
	// The default type of an untyped constant is bool, rune, int, float64, complex128, or string respectively

	const t = true        // default type
	const f bool = 1 == 0 // constant expression
	show("Boolean: ", t, f)

	const (
		r1 = '\xFF'
		r2 = '\u00FF'
	)
	show("Rune: ", r1, r2)

	const i1, i2, i3 = uint(iota), 987654321012345678, len("fooÿ")
	show("Integer: ", i1, i2, i3)

	const f1, f2, f3 = 0.1e1, 0xF.Fp-1 * 5.3, -0.0
	show("Float: ", f1, f2, f3)

	const c1, c2 = 1.2i, 3i / 5i
	show("Complex: ", c1, c2)

	const s1, s2 = "foo", `\r\n`
	show("String: ", s1, s2)
}

func variables() {
	show("A variable is a storage location for holding a value. The set of permissible values is determined by the variable's type")
	// A variable declaration or, ... the signature of a function declaration or function literal reserves storage for a named variable.

	// A variable's value is retrieved by referring to the variable in an expression;
	// it is the most recent value assigned to the variable.
	// If a variable has not yet been assigned a value, its value is the zero value for its type

	var a int
	var b int = 42
	var c = func(x int) (int, error) { return 42, nil }
	d := "xz"

	// Calling the built-in function `new` or taking the address of a composite literal allocates storage for a variable at run time.
	// Such an anonymous variable is referred to via a (possibly implicit) pointer indirection.

	type Point3D struct{ x, y, z float64 }
	origin := &Point3D{} // composite literal

	show("Values: ", a, b, c, d, origin, *origin)
	// Values: int(0); int(42); func(int) (int, error)(0x47ed20); string(xz); *main.Point3D(&{0 0 0}); main.Point3D({0 0 0});

	// static vs dynamic types

	// Static type of a variable is the type given in its declaration,
	// the type provided in the new call or composite literal,
	// or the type of an element of a structured variable.

	// Variables of interface type also have a distinct dynamic type,
	// which is the (non-interface) type of the value assigned to the variable at run time
	// (unless the value is the predeclared identifier nil, which has no type).
	// The dynamic type may vary during execution
	// but values stored in interface variables are always assignable to the static type of the variable.

	type T struct{ a, b byte }
	var v *T                       // v has value nil, static type *T
	var x interface{}              // x is nil and has static type interface{}
	show("Dynamic type 1: ", x, v) // Dynamic type 1: <nil>(<nil>); *main.T(<nil>);

	x = 42                      // x has value 42 and dynamic type int
	show("Dynamic type 2: ", x) // Dynamic type 2: int(42);

	x = v                       // x has value (*T)(nil) and dynamic type *T
	show("Dynamic type 3: ", x) // Dynamic type 3: *main.T(<nil>);
}

func types() {
	show("A type determines a set of values together with operations and methods specific to those values")

	var a any
	var b bool
	var c byte
	// var d comparable
	var e complex64
	var f complex128
	var g error
	var h float32
	var i float64
	var j int
	var k int8
	var l int16
	var m int32
	var n int64
	var o rune
	var p string
	var q uint
	var r uint8
	var s uint16
	var t uint32
	var u uint64
	var v uintptr
	show("Predeclared type names: ", a, b, c, e, f, g, h, i, j, k, l, m, n, o, p, q, r, s, t, u, v)

	// Composite types — array, struct, pointer, function, interface, slice, map, and channel types — may be constructed using type literals.
	type T struct{ a, b byte }
	show("Composite type: ", T{})

	// Predeclared types, defined types, and type parameters are called named types.
	// An alias denotes a named type if the type given in the alias declaration is a named type.
	type twoBytes = T
	show("Alias type: ", twoBytes{})

	var booleanTypes = func() {
		// predeclared constants true and false
		var t bool = true
		var f bool
		show("Boolean types: ", t, f, false)
	}

	var numericTypes = func() {
		// The predeclared architecture-independent numeric type
		show("Numeric types ...")
		show("Integer: ", uint8(0xFF), uint16(0xFFFF), uint32(0xFFFFFFFF), uint64(0xFFFFFFFFFFFFFFFF), int8(-0xFF>>1))
		show("Implementation specific integers:", uint(0xFFFFFFFFFFFFFFFF), int(0xFFFFFFFFFFFFFFFF>>1), uintptr(0xFFFFFFFFFFFFFFFF))
		show("Floats: ", float32(1e03), float64(1e308))
		show("Complex: ", complex64(1.2i), complex128(2.1i))
		show("Byte: ", byte(0xFF))
		show("Rune: ", '\xFF', rune(0xFF))
	}

	var stringTypes = func() {
		// A string value is a (possibly empty) sequence of bytes.
		// The number of bytes is called the length of the string and is never negative.
		// Strings are immutable
		show("String types ...")
		var a string = ""
		show("strings: ", a, "йцукен", len("йцукен"))

		var bs string = "йцукен"
		for i, b := range bs { // enumerate runes
			show("string runes: fist byte idx, first byte, rune: ", i, bs[i], b)
		}

		for i := 0; i < len(bs); i++ { // enumerate bytes
			show("string bytes: ", i, bs[i])
		}
	}

	var arrayTypes = func() {
		// The length is part of the array's type; it must evaluate to a non-negative constant representable by a value of type int
		show("Array types ...")
		var as [3]rune
		var bs [2]struct{ x, y byte }
		show("arrays: ", as, bs)

		for i, b := range bs {
			show("array element: ", i, b)
		}
		for i := 0; i < len(bs); i++ {
			show("array elements: ", i, bs[i], bs[i].x-bs[i].y)
		}
	}

	var sliceTypes = func() {
		// A slice is a descriptor for a contiguous segment of an underlying array
		// The value of an uninitialized slice is nil.
		// A slice therefore shares storage with its array and with other slices of the same array;
		// The capacity of a slice a can be discovered using the built-in function cap(a).
		// A slice created with make always allocates a new, hidden array
		// with slices of slices (or arrays of slices), the inner lengths may vary dynamically.
		// Moreover, the inner slices must be initialized individually.
		show("Slice types ...")
		type T = struct{ x, y byte }
		var as []T = make([]T, 2, 4)
		show("slice: ", as)
		show("new slices: ", new([16]byte)[7:9], make([]byte, 2))

		bs := [3]byte{1, 2, 3}
		show("array slices: ", bs[0:1], bs[1:2], bs[2:3], bs[0:3], bs[:])
		show("slice capacity: ", cap(bs[0:1]), cap(bs[2:3]))

		var cs [2][3]int
		show("slices 2D: ", cs[:], cs[1][:])
	}

	structTypes := func() {
		// A struct is a sequence of named elements, called fields, each of which has a name and a type
		// Field names may be specified explicitly (IdentifierList) or implicitly (EmbeddedField)
		// A field declared with a type but no explicit field name is called an embedded field
		// Promoted fields act like ordinary fields of a struct except that they cannot be used
		// as field names in composite literals of the struct.
		// The tags are made visible through a reflection interface and take part in type identity for structs but are otherwise ignored.
		show("Struct types ...")
		type A = struct{}
		type B = struct {
			x, y byte
			a    *[]A
		}
		type C = struct {
			A
			B         // B.[x,y,a] are promoted to C
			_ int64   // padding
			c float64 `your:"tag here"` // concatenation of optionally space-separated key:"value" pairs.
			x int64   ``                // empty tag = no tag
		}
		type D = struct {
			microsec  uint64 `protobuf:"1"`
			serverIP6 uint64 `protobuf:"2"`
		}
		show("structs: ", A{}, B{}, C{}, D{})
		show("empty struct: ", A{})
		show("struct with 3 fields: ", B{})
		show("struct with embedded fields A,B (promoted B.*), and padding 64bit, and tagged fields: ", C{}, C{}.B.y)
		show("properly tagged struct: ", D{})
	}

	pointerTypes := func() {
		// set of all pointers to variables of a given type, called the base type of the pointer.
		// The value of an uninitialized pointer is nil.
		show("Pointer types ...")
		type Point = struct{ x, y byte }
		var a *Point = new(Point)
		var b *[4]int
		show("pointers: ", a, b)
	}

	booleanTypes()
	numericTypes()
	stringTypes()
	arrayTypes()
	sliceTypes()
	structTypes()
	pointerTypes()

}

func show(msg string, xs ...any) {
	var line string = msg
	for _, x := range xs {
		line += fmt.Sprintf("%T(%v); ", x, x)
		// line += fmt.Sprintf("%#v; ", x)
	}
	fmt.Println(line)
}
