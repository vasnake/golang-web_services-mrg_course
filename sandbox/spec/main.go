package main

import (
	"fmt"
)

func main() {
	show("\nSveiki!")

	integers()
	floats()
	imaginary()
	runeLiterals()
	stringLiterals()
	constants()
	variables()
	types()

	show(`
Viso gero!
`)
}

func integers() {
	show("\nAn integer literal is a sequence of digits representing an integer constant")
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
	show("\nA floating-point literal is a decimal or hexadecimal representation of a floating-point constant")

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
	show("\nAn imaginary literal represents the imaginary part of a complex constant")
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
	show("\nA rune literal represents a rune constant, an integer value identifying a Unicode code point")
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
	show("\nA string literal represents a string constant obtained from concatenating a sequence of characters")
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
	show("\nConstants, numeric and string")
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
	show("\nA variable is a storage location for holding a value. The set of permissible values is determined by the variable's type")
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

	// If a variable has not yet been assigned a value, its value is the `zero value`` for its type.
}

func types() {
	show("\nA type determines a set of values together with operations and methods specific to those values")
	// A type may be denoted by a type name ...
	// A type may also be specified using a type literal, which composes a type from existing types.
	// `Composite` types — array, struct, pointer, function, interface, slice, map, and channel types
	// Predeclared types, defined types, and type parameters are called `named` types

	// To avoid portability issues all numeric types are `defined` types and thus distinct
	// Explicit conversions are required when different numeric types are mixed in an expression or assignment
	// Numeric types: uint8..64, int8..64, float32..64, complex64..128,
	// byte (alias uint8), rune (alias uint32)

	// predeclared integer types with implementation-specific sizes
	// uint     either 32 or 64 bits
	// int      same size as uint
	// uintptr  an unsigned integer large enough to store the uninterpreted bits of a pointer value

	var a any
	var b bool

	var c byte

	var e complex64
	var f complex128

	var h float32
	var i float64

	var j int
	var k int8
	var l int16
	var m int32
	var n int64

	var o rune
	var q uint
	var r uint8
	var s uint16
	var t uint32
	var u uint64
	var v uintptr

	var p string
	var g error

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
		show("\nBoolean types: ", t, f, false)
	}

	var numericTypes = func() {
		// The predeclared architecture-independent numeric type
		show("\nNumeric types ...")

		show("Integer: ", uint8(0xFF), uint16(0xFFFF), uint32(0xFFFFFFFF), uint64(0xFFFFFFFFFFFFFFFF), int8(-0xFF>>1))
		show("Implementation specific integers:", uint(0xFFFFFFFFFFFFFFFF), int(0xFFFFFFFFFFFFFFFF>>1), uintptr(0xFFFFFFFFFFFFFFFF))
		show("Floats: ", float32(1e03), float64(1e308))
		show("Complex: ", complex64(1.2i), complex128(2.1i))
		show("Byte: ", byte(0xFF))
		show("Rune: ", '\xFF', rune(0xFF))
	}

	var stringTypes = func() {
		show("\nA string value is a (possibly empty) sequence of bytes")
		// The number of bytes is called the length of the string and is never negative.
		// Strings are immutable
		// It is illegal to take the address of a string's byte (`&str[i]` is invalid)

		var a string = ""
		show("strings: ", a, "йцукен", len("йцукен"), runesCount("йцуукен")) // strings: string(); string(йцукен); int(12); int(7);
		show("string bytes, runes: ", []byte("йцукен"), []rune("йцукен"))
		// string bytes, runes: []uint8([208 185 209 134 209 131 208 186 208 181 208 189]); []int32([1081 1094 1091 1082 1077 1085])
		show("ASCII string bytes, runes: ", []byte("abc"), []rune("abc"))
		// ASCII string bytes, runes: []uint8([97 98 99]); []int32([97 98 99]);

		var bs string = "йцукен"
		for i, b := range bs { // enumerate runes
			show("string runes: rune first byte idx, rune first byte, rune: ", i, bs[i], b)
		}

		for i := 0; i < len(bs); i++ { // enumerate bytes
			show("string bytes: ", i, bs[i])
		}

		// string mutability
		show("string before, bytes, runes: ", []byte(bs), []rune(bs))
		// bs[0] = byte(42) // denied
		([]byte(bs))[0] = 42 // original string remains the same
		show("string after, bytes, runes: ", []byte(bs), []rune(bs))
		// string before, bytes, runes: []uint8([208 185 209 134 209 131 208 186 208 181 208 189]); []int32([1081 1094 1091 1082 1077 1085]);
		// string after, bytes, runes:  []uint8([208 185 209 134 209 131 208 186 208 181 208 189]); []int32([1081 1094 1091 1082 1077 1085]);
		var strBytes = []byte(bs)
		strBytes[0] = 42 // modified copy
		show("string after, bytes after: ", []byte(bs), strBytes)
		// string after, bytes after: []uint8([208 185 209 134 209 131 208 186 208 181 208 189]); []uint8([42 185 209 134 209 131 208 186 208 181 208 189]);
	}

	var arrayTypes = func() {
		// The length is part of the array's type; it must evaluate to a non-negative constant representable by a value of type int
		show("\nAn array is a numbered sequence of elements of a single type")

		var as [3]rune
		var bs [2]struct{ x, y byte }
		show("arrays: ", as, bs)
		// arrays: [3]int32([0 0 0]); [2]struct { x uint8; y uint8 }([{0 0} {0 0}]);

		// array elements
		for i, b := range bs {
			show("array element, range: ", i, b)
		}

		for i := 0; i < len(bs); i++ {
			show("array elements, for-loop: ", i, bs[i])
		}
	}

	var sliceTypes = func() {
		show("\nA slice is a descriptor for a contiguous segment of an underlying array")
		// The value of an uninitialized slice is nil.
		// A slice therefore shares storage with its array and with other slices of the same array;
		// The capacity of a slice `a` can be discovered using the built-in function `cap(a)`.

		// A slice created with `make` always allocates a new, hidden array
		var sliceA = make([]int, 50, 100)
		// is equivalent to creating a new array and slicing it
		var sliceB = (new([100]int))[0:50]

		var equals = func(a, b []int) bool {
			return len(a) == len(b) && cap(a) == cap(b) && func() bool {
				var res = true
				for i, ae := range a {
					if ae != b[i] {
						break
					}
				}
				return res
			}()
		}(sliceA, sliceB)
		show("Slices a and b are equal: ", equals)

		// with slices of slices (or arrays of slices), the inner lengths may vary dynamically.
		// Moreover, the inner slices must be initialized individually.

		type T = struct{ x, y byte }
		var as []T = make([]T, 2, 4)
		show("slice: ", as)
		show("new slices: ", new([16]byte)[7:9], make([]byte, 2))

		bs := [3]byte{1, 2, 3}
		show("slices of same array: ", bs[0:1], bs[1:2], bs[2:3], bs[0:3], bs[:])
		show("slice capacity: ", cap(bs[0:1]), cap(bs[2:3]))

		var cs [2][3]int
		show("slices 2D: ", cs[:], cs[1][:])
		// slices 2D: [][3]int([[0 0 0] [0 0 0]]); []int([0 0 0]);
	}

	structTypes := func() {
		show("\nA struct is a sequence of named elements, called fields, each of which has a name and a type")
		// Field names may be specified explicitly (IdentifierList) or implicitly (EmbeddedField)

		// A field declared with a type but no explicit field name is called an embedded field

		// Promoted fields act like ordinary fields of a struct except that they cannot be used
		// as field names in composite literals of the struct.

		// The tags are made visible through a reflection interface and take part in type identity for structs but are otherwise ignored.

		// A struct type T may not contain a field of type T, or of a type containing T as a component

		type A = struct{} // empty struct
		type B = struct { // 3 fields
			x, y byte
			a    *[]A
		}
		type C = struct {
			A         // A field declared with a type but no explicit field name is called an `embedded` field
			B         // B.[x,y,a] are promoted to C // Promoted fields act like ordinary fields of a struct
			_ int64   // blank identifier, padding
			c float64 `your:"tag here"` // concatenation of optionally space-separated key:"value" pairs.
			x int64   ``                // empty tag = no tag
		}
		type D = struct {
			microsec  uint64 `protobuf:"1"`
			serverIP6 uint64 `protobuf:"2"`
		}
		show("structs: ", A{}, B{}, C{}, D{})
		// structs: struct {}({}); struct { x uint8; y uint8; a *[]struct {} }({0 0 <nil>}); struct { struct {}; struct { x uint8; y uint8; a *[]struct {} }; _ int64; c float64 "your:\"tag here\""; x int64 }({{} {0 0 <nil>} 0 0 0}); struct { microsec uint64 "protobuf:\"1\""; serverIP6 uint64 "protobuf:\"2\"" }({0 0});

		show("empty struct: ", A{})
		// empty struct: struct {}({});

		show("struct with 3 fields: ", B{})
		// struct with 3 fields: struct { x uint8; y uint8; a *[]struct {} }({0 0 <nil>});

		show("struct with embedded fields A,B (promoted B.*), and padding 64bit, and tagged fields: ", C{}, C{}.B.y)
		// struct with embedded fields A,B (promoted B.*), and padding 64bit, and tagged fields: struct { struct {}; struct { x uint8; y uint8; a *[]struct {} }; _ int64; c float64 "your:\"tag here\""; x int64 }({{} {0 0 <nil>} 0 0 0}); uint8(0);

		show("properly tagged struct: ", D{})
		// properly tagged struct: struct { microsec uint64 "protobuf:\"1\""; serverIP6 uint64 "protobuf:\"2\"" }({0 0});
	}

	pointerTypes := func() {
		show("\nA pointer type denotes the set of all pointers to variables of a given type, called the base type of the pointer")
		// The value of an uninitialized pointer is nil.

		type Point = struct{ x, y byte }
		var a *Point = new(Point) // new returns pointer on an allocated object (with zero value), always
		var b *[4]int             // pointer zero value
		show("pointers: ", a, b)
		// pointers: *struct { x uint8; y uint8 }(&{0 0}); *[4]int(<nil>);
	}

	functionTypes := func() {
		show("\nA function type denotes the set of all functions with the same parameter and result types")
		// The value of an uninitialized variable of function type is nil.

		// The final incoming parameter in a function signature may have a type prefixed with `...`
		// A function with such a parameter is called `variadic` and may be invoked with zero or more arguments for that parameter

		// Named return values in signature
		// Any number of return values

		// examples
		var a = func() {}
		var b = func(x int) int { return x + 42 }
		var c = func(a, _ int, z float32) bool { return float32(a) > z }
		var d = func(a, b int, z float32) bool { return a == b && z > 0 }

		var e = func(prefix string, values ...int) {
			for i, x := range values {
				show(prefix, i, x)
			}
		}

		var f = func(a, b int, z float64, opt ...interface{}) (success bool) {
			success = (float64(a) - float64(b) + z - float64(len(opt))) > 0
			return // success
		}

		var g = func(int, int, float64) (float64, *[]int) { return 42.0, &[]int{1, 2} }
		var h = func(n int) func(p *T) { return func(x *T) { show("n, x: ", n, x) } }

		show("functions a, b, c: ", a, b, c)
		show("functions d, e, f: ", d, e, f)
		show("functions g, h: ", g, h)
		// functions a, b, c: func()(0x47fc80); func(int) int(0x47fca0); func(int, int, float32) bool(0x47fcc0);
		// functions d, e, f: func(int, int, float32) bool(0x47fce0); func(string, ...int)(0x47fd00); func(int, int, float64, ...interface {}) bool(0x47fe20);
		// functions g, h: func(int, int, float64) (float64, *[]int)(0x47fe60); func(int) func(*main.T)(0x47e540);
	}

	var interfaceTypes = func() {
		show("\nAn interface type defines a `type set`")
		// A variable of interface type can store a value of any type that is in the type set of the interface.
		// Such a type is said to `implement` the interface.
		// The value of an uninitialized variable of interface type is nil

		// An interface type is specified by a list of `interface elements`
		//  An interface element is either a method or a type

		// A type T implements an interface I if
		// T is not an interface and is an element of the type set of I; or
		// T is an interface and the type set of T is a subset of the type set of I.

		var basicInterfaces = func() {
			show("Interfaces whose type sets can be defined entirely by a list of methods are called basic interfaces.")
			// The type set defined by such an interface is the set of types which implement all of those methods

			// type illegal interface {
			// 	String() string
			// 	String() string  // illegal: String not unique
			// 	_(x int)         // illegal: method must have non-blank name
			// }

			// A simple File interface. // e.g.
			type file interface {
				Read(b []byte) (n int, err error)
				Write(b []byte) (n int, err error)
				Close() error
			}

			// More than one type may implement an interface
			// type T1 int32
			// type T2 uint32
			// func (f T1) Read(n []byte) (int, error)  { return 42, nil }
			// func (f T1) Write(n []byte) (int, error) { return 42, nil }
			// func (f T1) Close() error                { return nil }
			// func (f T2) Read(n []byte) (int, error)  { return 42, nil }
			// func (f T2) Write(n []byte) (int, error) { return 42, nil }
			// func (f T2) Close() error                { return nil }

			// Any given type may implement several distinct interfaces.
			// For instance, all types implement the empty interface (alias `any`)
			// which stands for the set of all (non-interface) types `interface{}`

		}

		var embeddedInterfaces = func() {
			show("Embedding interface E in T: an interface T may use a (possibly qualified) interface type name E as an interface element")
			// In a slightly more general form
			// an interface T may use a (possibly qualified) interface type name E as an interface element.
			// This is called embedding interface E in T

			// The type set of T is the intersection of the type sets defined by T's explicitly declared methods
			// and the type sets of T’s embedded interfaces
			// When embedding interfaces, methods with the same names must have identical signatures.

			type Reader interface {
				Read(p []byte) (n int, err error)
				Close() error
			}

			type Writer interface {
				Write(p []byte) (n int, err error)
				Close() error
			}

			// ReadWriter's methods are Read, Write, and Close.
			type ReadWriter interface {
				Reader // includes methods of Reader in ReadWriter's method set
				Writer // includes methods of Writer in ReadWriter's method set
			}

			// type ReadCloser interface {
			// 	Reader   // includes methods of Reader in ReadCloser's method set
			// 	Close()  // illegal: signatures of Reader.Close and Close are different
			// }
		}

		var generalInterfaces = func() {
			show("In their most general form, an interface element may also be an arbitrary type term T, or a term of the form ~T specifying the underlying type T, or a union of terms t1|t2|…|tn")
			// Together with method specifications, these elements enable the precise definition of an interface's type set as follows:
			// The type set of the empty interface is the set of all non-interface types.
			// The type set of a non-empty interface is the intersection of the type sets of its interface elements.
			// The type set of a method specification is the set of all non-interface types whose method sets include that method.
			// The type set of a non-interface type term is the set consisting of just that type.
			// The type set of a term of the form ~T is the set of all types whose underlying type is T.
			// The type set of a union of terms t1|t2|…|tn is the union of the type sets of the terms

			// Interfaces that are not `basic` may only be used as type constraints,
			// or as elements of other interfaces used as constraints.
			// They cannot be the types of values or variables, or components of other, non-interface types

			// An interface representing only the type int.
			type intInterface interface {
				int
			}

			// An interface representing all types with underlying type int.
			type allIntsInterface interface {
				~int
			}

			// An interface representing all types with underlying type int that implement the String method.
			type printableInt interface {
				~int
				String() string
			}

			// An interface representing an empty type set: there is no type that is both an int and a string.
			type emptyTypesSet interface {
				int
				string
			}

			// In a term of the form ~T, the underlying type of T must be itself, and T cannot be an interface.
			type MyInt int
			// type illegal interface {
			// 	~[]byte  // the underlying type of []byte is itself
			// 	~MyInt   // illegal: the underlying type of MyInt is not MyInt
			// 	~error   // illegal: error is an interface
			// }

			// Union elements denote unions of type sets:
			// The Float interface represents all floating-point types
			// (including any named types whose underlying types are either float32 or float64).
			type Float interface {
				~float32 | ~float64
			}

			// Given a type parameter P:
			// type illegal interface {
			// 	P                // illegal: P is a type parameter
			// 	int | ~P         // illegal: P is a type parameter
			// 	~int | MyInt     // illegal: the type sets for ~int and MyInt are not disjoint (~int includes MyInt)
			// 	float32 | Float  // overlapping type sets but Float is an interface
			// }

		}

		basicInterfaces()
		embeddedInterfaces()
		generalInterfaces()
	}

	var mapTypes = func() {
		show("\nA map is an unordered group of elements of one type, called the element type, indexed by a set of unique keys of another type, called the key type")
		// The value of an uninitialized map is nil.

		// The comparison operators == and != must be fully defined for operands of the key type;
		// thus the key type must not be a function, map, or slice.
		// If the key type is an interface type, these comparison operators must be defined for the dynamic key values; failure will cause a run-time panic.
		var a map[string]int
		var b map[*T]struct{ x, y float64 }
		var c map[string]interface{}
		show("Maps examples: ", a, b, c)

		// Elements may be added during execution using assignments and retrieved with index expressions;
		// they may be removed with the `delete` and `clear` built-in function

		// A new, empty map value is made using the built-in function `make`, which takes the map type and an optional capacity hint as arguments:
		a = make(map[string]int)
		a = make(map[string]int, 100)

		// A nil map is equivalent to an empty map except that no elements may be added
	}

	var channelTypes = func() {
		show("\n")
	}

	booleanTypes()
	numericTypes()
	stringTypes()
	arrayTypes()
	sliceTypes()
	structTypes()
	pointerTypes()
	functionTypes()
	interfaceTypes()
	mapTypes()
	channelTypes()
}

func show(msg string, xs ...any) {
	var line string = msg
	for _, x := range xs {
		line += fmt.Sprintf("%T(%v); ", x, x)
		// line += fmt.Sprintf("%#v; ", x)
	}
	fmt.Println(line)
}

var runesCount = func(str string) int {
	var runes = []rune(str) // allocation?
	return len(runes)
}
