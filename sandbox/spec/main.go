package main

import (
	"fmt"
	sfs "spec/functions"
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
	propertiesOfTypesAndValues()
	blocks()
	declarationsAndScope()
	expressions()

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
		show("strings: ", a, "йцукен", len("йцукен"), sfs.RunesCount("йцуукен")) // strings: string(); string(йцукен); int(12); int(7);
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
		show("\nA channel provides a mechanism for concurrently executing functions to communicate by sending and receiving values of a specified element type")
		// The value of an uninitialized channel is nil.

		var a chan T         // can be used to send and receive values of type T
		var b chan<- float64 // can only be used to send float64s
		var c <-chan int     // can only be used to receive ints
		show("Channels: ", a, b, c)

		// A channel may be constrained only to send or only to receive by assignment or explicit conversion.
		// The optional <- operator specifies the channel direction, send or receive
		// If a direction is given, the channel is `directional`, otherwise it is `bidirectional`

		// The <- operator associates with the leftmost chan possible:
		var d chan<- chan int   // same as chan<- (chan int)
		var g chan (<-chan int) // you have to use parentheses
		var e chan<- <-chan int // same as chan<- (<-chan int)
		var f <-chan <-chan int // same as <-chan (<-chan int)
		show("Channel <- operator: ", d, e, f, g)

		// A new, initialized channel value can be made using the built-in function make,
		// which takes the channel type and an optional capacity as arguments:
		var h = make(chan int, 100)
		show("buffered channel: ", h, len(h), cap(h)) // buffered channel: chan int(0xc00012e000); int(0); int(100);
		// The capacity, in number of elements, sets the size of the buffer in the channel.
		// If the capacity is zero or absent, the channel is unbuffered and communication succeeds only when both a sender and receiver are ready.
		// Otherwise, the channel is buffered and communication succeeds without blocking ...

		// A channel may be closed with the built-in function `close`
		// ... the built-in function close records that no more values will be sent on the channel.
		// It is an error if ch is a receive-only channel.
		// Sending to or closing a closed channel causes a run-time panic.
		// After calling close, ... receive operations will return the zero value for the channel's type without blocking.
		close(h)
		show("closed channel: ", h)

		// A single channel may be used in:
		// send statements, receive operations, calls to the built-in functions cap and len
		// by any number of goroutines without further synchronization.
		// Channels act as first-in-first-out queues

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

func propertiesOfTypesAndValues() {
	show("\nProperties of types and values ...")
	// Underlying types

	var underlyingTypes = func() {
		show("Each type T has an underlying type")
		// If T is one of the predeclared boolean, numeric, or string types, or a type literal, the corresponding underlying type is T itself.
		// Otherwise, T's underlying type is the underlying type of the type to which T refers in its declaration.
		// For a type parameter that is the underlying type of its type constraint, which is always an interface.

		// The underlying type of string, A1, A2, B1, and B2 is string.
		// The underlying type of []B1, B3, and B4 is []B1.
		// The underlying type of P is interface{}
		type (
			A1 = string
			A2 = A1
		)
		type (
			B1 string
			B2 B1
			B3 []B1
			B4 B3
		)
		// func f[P any](x P) { return }
	}

	var coreTypes = func() {
		show("By definition, a core type is never a defined type, type parameter, or interface type")
		// Each non-interface type T has a core type, which is the same as the underlying type of T

		// An interface T has a core type if one of the following conditions is satisfied
		// 1) There is a single type U which is the underlying type of all types in the type set of T
		// or 2) the type set of T contains only channel types with identical element type E, and all directional channels have the same direction

		type Celsius float32
		type Kelvin float32
		type data byte
		type myString interface{ string }

		// Examples of interfaces with core types:
		type a interface{ int }                     // int
		type b interface{ Celsius | Kelvin }        // float32
		type c interface{ ~chan int }               // chan int
		type d interface{ ~chan int | ~chan<- int } // chan<- int
		type e interface {
			~[]*data
			String() string
		} // []*data

		// Examples of interfaces without core types:
		type f interface{}                           // no single underlying type
		type g interface{ Celsius | float64 }        // no single underlying type
		type h interface{ chan int | chan<- string } // channels have different element types
		type i interface{ <-chan int | chan<- int }  // directional channels have different directions

		// Some operations (slice expressions, append and copy) rely on a slightly more loose form of core types which accept byte slices and strings
		// Note that `bytestring` is not a real type; it cannot be used to declare variables or compose other types.
		// It exists solely to describe the behavior of some operations that read from a sequence of bytes, which may be a byte slice or a string.
		type j interface{ []byte | string }    // bytestring
		type k interface{ ~[]byte | myString } // bytestring
	}

	var typeIdentity = func() {
		show("A named type is always different from any other type")
		// Otherwise, two types are identical if their underlying type literals are structurally equivalent
		// Two struct types are identical ... Non-exported field names from different packages are always different

		// Two function types are identical if they have the same number of parameters and result values,
		// corresponding parameter and result types are identical, and either both functions are variadic or neither is.
		// Parameter and result names are not required to match

		// Given the declarations (what about named types?)
		type (
			A0 = []string
			A1 = A0
			A2 = struct{ a, b int }
			A3 = int
			A4 = func(A3, float64) *A0
			A5 = func(x int, _ float64) *[]string

			B0 A0
			C0 = B0

			B1 []string
			// B0 and B1 are different because they are new types created by distinct type definitions

			B2 struct{ a, b int }
			B3 struct{ a, c int }

			B4 func(int, float64) *B0
			// B4 vs A5: func(int, float64) *B0 and func(x int, y float64) *[]string are different because B0 is different from []string

			B5 func(x int, y float64) *A1

			D0[P1, P2 any] struct {
				x P1
				y P2
			}
			// P1 and P2 are different because they are different type parameters
			E0 = D0[int, string]
			// D0[int, string] and struct{ x int; y string } are different because the former is an instantiated defined type
			// while the latter is a type literal
		)
		// these types are identical:
		// A0, A1, and []string
		// A2 and struct{ a, b int }
		// A3 and int
		// A4, func(int, float64) *[]string, and A5

		// B0 and C0
		// D0[int, string] and E0
		// []int and []int
		// struct{ a, b *B5 } and struct{ a, b *B5 }
		// func(x int, y float64) *[]string, func(int, float64) (result *[]string), and A5

	}

	var assignability = func() {
		show(`A value x of type V is assignable to a variable of type T ("x is assignable to T") if one of the following conditions applies`)
		// V and T are identical.
		// V and T have identical underlying types but are not type parameters and at least one of V or T is not a named type.
		// V and T are channel types with identical element types, V is a bidirectional channel, and at least one of V or T is not a named type.
		// T is an interface type, but not a type parameter, and x implements T.
		// x is the predeclared identifier nil and T is a pointer, function, slice, map, channel, or interface type, but not a type parameter.
		// x is an untyped constant representable by a value of type T.

		// Additionally, if ... V or T are type parameters, x is assignable to a variable of type T if one of the following conditions applies:
		// x is the predeclared identifier nil, T is a type parameter, and x is assignable to each type in T's type set.
		// V is not a named type, T is a type parameter, and x is assignable to each type in T's type set.
		// V is a type parameter and T is not a named type, and values of each type in V's type set are assignable to T.
	}

	var representability = func() {
		show("A constant x is representable by a value of type T, where T is not a type parameter, if one of the following conditions applies")
		// If T is a type parameter, x is representable by a value of type T if x is representable by a value of each type in T's type set

		// x                   T           x is representable by a value of T because:
		// 'a'                 byte        97 is in the set of byte values
		// 97                  rune        rune is an alias for int32, and 97 is in the set of 32-bit integers
		// "foo"               string      "foo" is in the set of string values
		// 1024                int16       1024 is in the set of 16-bit integers
		// 42.0                byte        42 is in the set of unsigned 8-bit integers
		// 1e10                uint64      10000000000 is in the set of unsigned 64-bit integers
		// 2.718281828459045   float32     2.718281828459045 rounds to 2.7182817 which is in the set of float32 values
		// -1e-1000            float64     -1e-1000 rounds to IEEE -0.0 which is further simplified to 0.0
		// 0i                  int         0 is an integer value
		// (42 + 0i)           float32     42.0 (with zero imaginary part) is in the set of float32 values

		// x                   T           x is not representable by a value of T because
		// 0                   bool        0 is not in the set of boolean values
		// 'a'                 string      'a' is a rune, it is not in the set of string values
		// 1024                byte        1024 is not in the set of unsigned 8-bit integers
		// -1                  uint16      -1 is not in the set of unsigned 16-bit integers
		// 1.1                 int         1.1 is not an integer value
		// 42i                 float32     (0 + 42i) is not in the set of float32 values
		// 1e1000              float64     1e1000 overflows to IEEE +Inf after rounding
	}

	var methodSets = func() {
		show("The method set of a type determines the methods that can be called on an operand of that type")
		// Every type has a (possibly empty) method set associated with it:
		// - The method set of a defined type T consists of all methods declared with receiver type T.
		// - The method set of a pointer to a defined type T (where T is neither a pointer nor an interface) is
		// the set of all methods declared with receiver *T or T.
		// - The method set of an interface type is the intersection of the method sets of each type in the interface's type set
		// (the resulting method set is usually just the set of declared methods in the interface).

		// urther rules apply to structs ... containing embedded fields
	}

	underlyingTypes()
	coreTypes()
	typeIdentity()
	assignability()
	representability()
	methodSets()
}

func blocks() {
	show("\nA block is a possibly empty sequence of declarations and statements within matching brace brackets")
	// In addition to explicit blocks in the source code, there are implicit blocks:
	// - The universe block encompasses all Go source text.
	// - Each package has a package block containing all Go source text for that package.
	// - Each file has a file block containing all Go source text in that file.
	// - Each "if", "for", and "switch" statement is considered to be in its own implicit block.
	// - Each clause in a "switch" or "select" statement acts as an implicit block.

	// Blocks nest and influence scoping
	{
		var x = 42
		show("x 1: ", x)
		{
			var x = 24
			show("x 2: ", x)
		}
		show("x 3: ", x)
		// x 1: int(42);
		// x 2: int(24);
		// x 3: int(42);
	}

}

func declarationsAndScope() {
	show("\nA declaration binds a non-blank identifier to a constant, type, type parameter, variable, function, label, or package")
	// Every identifier in a program must be declared.
	// No identifier may be declared twice in the same block, and no identifier may be declared in both the file and package block

	//  In the package block, the identifier `init` may only be used for init function declarations

	// The package clause is not a declaration; the package name does not appear in any scope.
	// Its purpose is to identify the files belonging to the same package and to specify the default package name for import declarations.

	// Go is lexically scoped using blocks:
	//-  The scope of a `predeclared` identifier is the `universe` block.
	//-  The scope of an identifier denoting a constant, type, variable, or function (but not method) declared at top level
	// (outside any function) is the package block.
	//-  The scope of the package name of an imported package is the file block of the file containing the import declaration.
	//-  The scope of an identifier denoting a method receiver, function parameter, or result variable is the function body.
	//-  The scope of an identifier denoting a type parameter of a function or declared by a method receiver begins
	// after the name of the function and ends at the end of the function body.
	//-  The scope of an identifier denoting a type parameter of a type begins after the name of the type and ends at the end of the TypeSpec.
	//-  The scope of a constant or variable identifier declared inside a function begins
	// at the end of the ConstSpec or VarSpec (ShortVarDecl for short variable declarations) and ends
	// at the end of the innermost containing block.
	//-  The scope of a type identifier declared inside a function begins
	// at the identifier in the TypeSpec and ends at the end of the innermost containing block.

	var labelScopes = func() {
		show(`Labels are declared by labeled statements and are used in the "break", "continue", and "goto" statements`)
		// In contrast to other identifiers, labels are not block scoped and do not conflict with identifiers that are not labels.
		// The scope of a label is the body of the function in which it is declared and excludes the body of any nested function
	}

	var blankIdentifier = func() {
		show("The blank identifier is represented by the underscore character _")
		// It serves as an anonymous placeholder instead of a regular (non-blank) identifier
		// and has special meaning in declarations, as an operand, and in assignment statements.
	}

	var predeclaredIdentifiers = func() {
		show("The following identifiers are implicitly declared in the universe block")
		// Types:
		// 		any bool byte comparable
		// 		complex64 complex128 error float32 float64
		// 		int int8 int16 int32 int64 rune string
		// 		uint uint8 uint16 uint32 uint64 uintptr

		// 	Constants:
		// 		true false iota

		// 	Zero value:
		// 		nil

		// 	Functions:
		// 		append cap clear close complex copy delete imag len
		// 		make max min new panic print println real recover

	}

	var exportedIdentifiers = func() {
		show("An identifier may be exported to permit access to it from another package")
		// An identifier is exported if both:
		//- the first character of the identifier's name is a Unicode uppercase letter (Unicode character category Lu); and
		//- the identifier is declared in the package block or it is a field name or method name.
	}

	var constantDeclarations = func() {
		show("A constant declaration binds a list of identifiers (the names of the constants) to the values of a list of constant expressions")
		// If the type is omitted, the constants take the individual types of the corresponding expressions
		const Pi float64 = 3.14159265358979323846
		const zero = 0.0 // untyped floating-point constant
		const (
			size int64 = 1024
			eof        = -1 // untyped integer constant
		)
		const a, b, c = 3, 4, "foo" // a = 3, b = 4, c = "foo", untyped integer and string constants
		const u, v float32 = 0, 3   // u = 0.0, v = 3.0

		// Within a parenthesized const declaration list the expression list may be omitted from any but the first ConstSpe
		// Such an empty list is equivalent to the textual substitution of the first preceding non-empty expression list and its type if any.
		// Omitting the list of expressions is therefore equivalent to repeating the previous list
		// Together with the iota constant generator this mechanism permits light-weight declaration of sequential values:
		const (
			Sunday = iota
			Monday
			Tuesday
			Wednesday
			Thursday
			Friday
			Partyday
			numberOfDays // this constant is not exported
		)
		show("iota generator: ", Sunday, Monday, Tuesday)
	}

	var iotaX = func() {
		show("Within a constant declaration, the predeclared identifier iota represents successive untyped integer constants.")
		// Its value is the index of the respective ConstSpec in that constant declaration, starting at zero
		const (
			c0 = iota // c0 == 0
			c1 = iota // c1 == 1
			c2 = iota // c2 == 2
		)
		const (
			a = 1 << iota // a == 1  (iota == 0)
			b = 1 << iota // b == 2  (iota == 1)
			c = 3         // c == 3  (iota == 2, unused)
			d = 1 << iota // d == 8  (iota == 3)
		)
		const (
			u         = iota * 42 // u == 0     (untyped integer constant)
			v float64 = iota * 42 // v == 42.0  (float64 constant)
			w         = iota * 42 // w == 84    (untyped integer constant)
		)
		const x = iota // x == 0
		const y = iota // y == 0

		// By definition, multiple uses of iota in the same ConstSpec all have the same value:
		const (
			bit0, mask0 = 1 << iota, 1<<iota - 1 // bit0 == 1, mask0 == 0  (iota == 0)
			bit1, mask1                          // bit1 == 2, mask1 == 1  (iota == 1)
			_, _                                 //                        (iota == 2, unused)
			bit3, mask3                          // bit3 == 8, mask3 == 7  (iota == 3)
		)
		// This last example exploits the implicit repetition of the last non-empty expression list
	}

	var typeDeclarations = func() {
		show("A type declaration binds an identifier, the type name, to a type")
		// Type declarations come in two forms: alias declarations and type definitions

		show("Type definitions")
		// A type definition creates a new, distinct type with the same underlying type and operations as the given type
		// and binds an identifier, the type name, to it
		// The new type is called a `defined type`. It is different from any other type, including the type it is created from
		type (
			Point struct{ x, y float64 } // Point and struct{ x, y float64 } are different types
			polar Point                  // polar and Point denote different types
			Node  struct{ value any }
		)

		show("Alias declarations")
		// An alias declaration binds an identifier to the given type
		// Within the scope of the identifier, it serves as an alias for the type
		type (
			nodeList = []*Node // nodeList and []*Node are identical types
			Polar    = polar   // Polar and polar denote identical types
		)

		// A defined type may have methods associated with it
		type (
			Block interface {
				BlockSize() int
				Encrypt(src, dst []byte)
				Decrypt(src, dst []byte)
			}
		)
		// It does not inherit any methods bound to the given type, but the method set of an interface type
		// or of elements of a composite type remains unchanged, e.g:

		// A Mutex is a data type with two methods, Lock and Unlock.
		type Mutex struct { /* Mutex fields */
		}
		// func (m *Mutex) Lock()   { /* Lock implementation */ }
		// func (m *Mutex) Unlock() { /* Unlock implementation */ }

		// NewMutex has the same composition as Mutex but its method set is empty
		type NewMutex Mutex

		// The method set of *PrintableMutex contains the methods
		// Lock and Unlock bound to its embedded field Mutex.
		type PrintableMutex struct{ Mutex }

		// e.g. EOF

		// Type definitions may be used to define different boolean, numeric, or string types and associate methods with them
		type TimeZone int
		const (
			EST TimeZone = -(5 + iota)
			CST
			MST
			PST
		)
		// func (tz TimeZone) String() string { return fmt.Sprintf("GMT%+dh", tz) }

		// If the type definition specifies type parameters, the type name denotes a generic type
		type List[T any] struct {
			next  *List[T]
			value T
		}
		// A generic type may also have methods associated with it
		// func (l *List[T]) Len() int  { ??? }
	}

	var typeParameterDeclarations = func() {
		show("A type parameter list declares the type parameters of a generic function or type declaration")
		// The type parameter list looks like an ordinary function parameter list
		// The type parameter is replaced with a type argument upon instantiation of the generic function or type

		// [P any]
		type a[P any] struct{ x P }
		// [S interface{ ~[]byte|string }]
		type b[S interface{ ~[]byte | string }] struct{ x S }
		// [S ~[]E, E any]
		type c[S ~[]E, E any] struct{ x S }
		// [P Constraint[int]]
		// [_ any]

		// Just as each ordinary function parameter has a parameter type,
		// each type parameter has a corresponding (meta-)type which is called its type constraint.

		// A parsing ambiguity arises when the type parameter list for a generic type declares a single type parameter P with a constraint C
		// such that the text `P C` forms a valid expression
		// type T[P *C] …
		// type T[P (C)] …
		// type T[P *C|Q] …
		// …
		// To resolve the ambiguity, embed the constraint in an interface or use a trailing comma:
		// type T[P interface{*C}] …
		// type T[P *C,] …

		// no recursion
		// type T1[P T1[P]] …                    // illegal: T1 refers to itself
		// type T2[P interface{ T2[int] }] …     // illegal: T2 refers to itself
		// type T3[P interface{ m(T3[int])}] …   // illegal: T3 refers to itself
		// type T4[P T5[P]] …                    // illegal: T4 refers to T5 and
		// type T5[P T4[P]] …                    // illegal: T5 refers to T4
		// type T6[P int] struct{ f *T6[P] }     // ok: reference to T6 is not in type parameter list

		var typeConstraints = func() {
			show("A type constraint is an interface")
			// that defines the set of permissible type arguments for the respective type parameter
			// and controls the operations supported by values of that type parameter.

			// in a type parameter list the enclosing interface{ … } may be omitted for convenience:
			// [T []P]                      // = [T interface{[]P}]
			// [T ~int]                     // = [T interface{~int}]
			// [T int|string]               // = [T interface{int|string}]
			// type Constraint ~int         // illegal: ~int is not in a type parameter list

			// The `comparable` interface and interfaces that (directly or indirectly) embed `comparable` may only be used as type constraints
			// The predeclared interface type `comparable` denotes the set of all non-interface types that are strictly comparable
			type a[P comparable] struct{ x P }
			type b[P interface{ comparable }] struct{ x P }

			// Satisfying a type constraint

			// comparing operands of type parameter type may panic at run-time:
			// A type argument T satisfies a type constraint C if ... T implements C
			// As an exception, a `strictly comparable` type constraint may also be satisfied by a `comparable`
			// (not necessarily strictly comparable) type argument
		}

		typeConstraints()
	}

	var variableDeclarations = func() {
		show("\nA variable declaration creates one or more variables")
		// binds corresponding identifiers to them, and gives each a type and an initial value
		// e.g.
		var i int
		var U, V, W float64
		var k = 0
		var x, y float32 = -1, -2
		var (
			I       int
			u, v, s = 2.0, 3.0, "bar"
		)
		// var re, im = complexSqrt(-1)
		// var _, found = entries[name]  // map lookup; only interested in "found"
		show("vars: ", i, U, V, W, k, x, y, u, v, s, I)

		// The predeclared value nil cannot be used to initialize a variable with no explicit type
		// var d = math.Sin(0.5)  // d is float64
		// var i = 42             // i is int
		// var t, ok = x.(T)      // t is T, ok is bool
		// var n = nil            // illegal
	}

	var shortVariableDeclarations = func() {
		show("\nA short variable declaration ... is shorthand for a regular variable declaration with initializer expressions but no types")
		// Short variable declarations may appear only inside functions

		// e.g.
		i, j := 0, 10
		f := func() int { return 7 }
		ch := make(chan int)
		// r, w, _ := os.Pipe()  // os.Pipe() returns a connected pair of Files and an error, if any
		// _, y, _ := coord(p)   // coord() returns three values; only interested in y coordinate
		show("vars: ", i, j, f, ch)

		// a short variable declaration may redeclare variables
		// provided they were originally declared earlier in the same block (or the parameter lists if the block is the function body)
		// with the same type, and at least one of the non-blank variables is new
		i, k := 1, 11
		show("redeclared i, add new k: ", i, k)
	}

	var functionDeclarations = func() {
		show("\nA function declaration binds an identifier, the function name, to a function")
		// If the function's signature declares result parameters, the function body's statement list must end in a terminating statement

		// A generic function must be instantiated before it can be called or used as a value.

		// A function declaration without type parameters may omit the body.
		// Such a declaration provides the signature for a function implemented outside Go, such as an assembly routine
	}

	var methodDeclarations = func() {
		show("A method is a function with a receiver")
		// A method declaration binds an identifier, the method name, to a method,
		// and associates the method with the receiver's base type (n.b. not pointer)

		//  Its type must be a defined type T or a pointer to a defined type T, ... T is called the receiver base type
		// A receiver base type cannot be a pointer or interface type and it must be defined in the same package as the method
		// The method is said to be bound to its receiver base type and the method name is visible only within selectors for type T or *T

		// Given defined type Point the declarations
		// func (p *Point) Length() float64 {
		// 	return math.Sqrt(p.x * p.x + p.y * p.y)
		// }
		// func (p *Point) Scale(factor float64) {
		// 	p.x *= factor
		// 	p.y *= factor
		// }
		// bind the methods Length and Scale, with receiver type *Point, to the base type Point.

		// If the receiver base type is a generic type, the receiver specification must declare corresponding type parameters
		// type Pair[A, B any] struct { a A; b B }
		// func (p Pair[A, B]) Swap() Pair[B, A] { return Pair[B, A]{p.b, p.a} } // receiver declares A, B
		// func (p Pair[First, _]) First() First { return p.a }                  // receiver declares First, corresponds to A in Pair
	}

	labelScopes()
	blankIdentifier()
	predeclaredIdentifiers()
	exportedIdentifiers()

	// Uniqueness of identifiers
	show("Two identifiers are different if they are spelled differently, or if they appear in different packages and are not exported.")

	constantDeclarations()
	iotaX()
	typeDeclarations()
	typeParameterDeclarations()
	variableDeclarations()
	shortVariableDeclarations()
	functionDeclarations()
	methodDeclarations()

}

func expressions() {
	show("\nAn expression specifies the computation of a value by applying operators and functions to operands")

	// Operands
	// Operands denote the elementary values in an expression
	// An operand may be a literal, a ... identifier denoting a constant, variable, or function,
	// or a parenthesized expression.

	// Qualified identifiers
	// A qualified identifier is an identifier qualified with a package name prefix
	// The identifier must be exported and declared in the package block of that package

	var compositeLiterals = func() {
		show("Composite literals construct new composite values each time they are evaluated")
		// The LiteralType's core type T must be a struct, array, slice, or map type
		// The key is interpreted as a field name for struct literals, an index for array and slice literals, and a key for map literals

		// Given the declarations
		type Point struct{ x, y float64 }
		type Point3D struct{ x, y, z float64 }
		type Line struct{ p, q Point3D }
		// one may write
		origin := Point3D{}                           // zero value for Point3D
		line := Line{origin, Point3D{y: -4, z: 12.3}} // zero value for line.q.x
		show("zero values: ", origin, line, Line{})
		// zero values: main.Point3D({0 0 0}); main.Line({{0 0 0} {0 -4 12.3}}); main.Line({{0 0 0} {0 0 0}});

		// For array and slice literals the following rules apply:
		//- Each element has an associated integer index marking its position in the array.
		//- An element with a key uses the key as its index.
		// The key must be a non-negative constant representable by a value of type int; and if it is typed it must be of integer type.
		//- An element without a key uses the previous element's index plus one.
		// If the first element has no key, its index is zero.

		// Taking the address of a composite literal generates a pointer to a unique variable initialized with the literal's value.
		var pointer *Point3D = &Point3D{y: 1000}
		var zvPointer *Point3D
		show("pointer: ", pointer, zvPointer)
		// pointer: *main.Point3D(&{0 1000 0}); *main.Point3D(<nil>);

		// Note that the zero value for a slice or map type is not the same as an initialized but empty value of the same type.
		// Consequently, taking the address of an empty slice or map composite literal does not have the same effect as
		// allocating a new slice or map value with new.
		p1 := &[]int{}   // p1 points to an initialized, empty slice with value []int{} and length 0
		p2 := new([]int) // p2 points to an uninitialized slice with value nil and length 0
		show("empty slices: ", p1, p2)
		// empty slices: *[]int(&[]); *[]int(&[]);

		// Array ... The notation `...` specifies an array length equal to the maximum element index plus one.
		buffer := [10]string{}            // len(buffer) == 10
		intSet := [6]int{1, 2, 3, 5}      // len(intSet) == 6
		days := [...]string{"Sat", "Sun"} // len(days) == 2
		show("arrays: ", buffer, intSet, days)
		// arrays: [10]string([         ]); [6]int([1 2 3 5 0 0]); [2]string([Sat Sun]);

		// A slice literal describes the entire underlying array literal.
		// Thus the length and capacity of a slice literal are the maximum element index plus one
		xs := []int{3, 7, 11}
		// and is shorthand for a slice operation applied to an array:
		tmp := [3]int{3, 7, 11}
		xs = tmp[0:3]
		show("slice: ", xs)

		// Within a composite literal of array, slice, or map type T,
		// elements or map keys that are themselves composite literals
		// may elide the respective literal type if it is identical to the element or key type of T
		var a = [...]Point{{1.5, -3.5}, {0, 0}}  // same as [...]Point{Point{1.5, -3.5}, Point{0, 0}}
		var b = [][]int{{1, 2, 3}, {4, 5}}       // same as [][]int{[]int{1, 2, 3}, []int{4, 5}}
		var c = [][]Point{{{0, 1}, {1, 2}}}      // same as [][]Point{[]Point{Point{0, 1}, Point{1, 2}}}
		var d = map[string]Point{"orig": {0, 0}} // same as map[string]Point{"orig": Point{0, 0}}
		var e = map[Point]string{{0, 0}: "orig"} // same as map[Point]string{Point{0, 0}: "orig"}

		type PPoint *Point
		var f = [2]*Point{{1.5, -3.5}, {}} // same as [2]*Point{&Point{1.5, -3.5}, &Point{}}
		var g = [2]PPoint{{1.5, -3.5}, {}} // same as [2]PPoint{PPoint(&Point{1.5, -3.5}), PPoint(&Point{})}
		show("elide literal type: ", a, b, c, d, e, f, g)
		// elide literal type: [2]main.Point([{1.5 -3.5} {0 0}]); [][]int([[1 2 3] [4 5]]); [][]main.Point([[{0 1} {1 2}]]); map[string]main.Point(map[orig:{0 0}]); map[main.Point]string(map[{0 0}:orig]); [2]*main.Point([0xc000114d30 0xc000114d40]); [2]main.PPoint([0xc000114d50 0xc000114d60]);

		// To resolve the ambiguity, the composite literal must appear within parentheses.
		// if x == (T{a,b,c}[i]) { … }
		// if (x == T{a,b,c}[i]) { … }
		// A parsing ambiguity arises when a composite literal using the TypeName form of the LiteralType appears as an operand
		// between the keyword and the opening brace of the block of an "if", "for", or "switch" statement ...

		// some examples:
		// vowels[ch] is true if ch is a vowel
		vowels := [128]bool{'a': true, 'e': true, 'i': true, 'o': true, 'u': true, 'y': true}
		show("vowels: ", vowels)
	}
	compositeLiterals()

	// Function literals
	// A function literal represents an anonymous function.
	// Function literals cannot declare type parameters.
	// A function literal can be assigned to a variable or invoked directly
	a := func(a, b int, z float64) bool { return a*b < int(z) }
	b := func(a, b int, z float64) bool { return a*b < int(z) }(1, 2, 3)
	show("function literals: ", a, b)
	// Function literals are closures: they may refer to variables defined in a surrounding function

	// Primary expressions
	// Primary expressions are the operands for unary and binary expressions
	// x
	// 2
	// (s + ".txt")
	// f(3.1415, true)
	// Point{1, 2}
	// m["foo"]
	// s[i : j + 1]
	// obj.color
	// f.p[i].x()

	var selectors = func() {
		show("Selectors, in `x.f` the identifier f is called the (field or method) selector")
		// For a primary expression x that is not a package name, the selector expression `x.f` denotes the field or method f of the value
		// If x is a package name, see the section on qualified identifiers
		// The number of embedded fields traversed to reach f is called its `depth in T`

		// As an exception, if the type of x is a defined pointer type and `(*x).f` is a valid selector expression denoting a field
		// (but not a method), `x.f` is shorthand for `(*x).f`

		// For example, given the declarations:
		// type T0 struct { x int }
		// (func(*T0) M0)()
		// type T1 struct { y int }
		// func(T1) M1()
		// type T2 struct {
		// 	z int
		// 	T1
		// 	*T0
		// }
		// func(*T2) M2()
		// type Q *T2
		// var t T2  // with t.T0 != nil
		// var p *T2 // with p != nil and (*p).T0 != nil
		// var q Q = p

		// one may write:
		// t.z // t.z
		// t.y // t.T1.y
		// t.x // (*t.T0).x
		// p.z // (*p).z
		// p.y // (*p).T1.y
		// p.x // (*(*p).T0).x
		// q.x // (*(*q).T0).x        (*q).x is a valid field selector
		// p.M0() // ((*p).T0).M0()      M0 expects *T0 receiver
		// p.M1() // ((*p).T1).M1()      M1 expects T1 receiver
		// p.M2() // p.M2()              M2 expects *T2 receiver
		// t.M2() // (&t).M2()           M2 expects *T2 receiver, see section on Calls

		// but the following is invalid:
		// q.M0()       // (*q).M0 is valid but not a field selector
	}
	selectors()

	var methodExpressions = func() {
		show("Method expressions, `T.M` is a function that is callable as a regular function with the same arguments as `M`")
		// prefixed by an additional argument that is the receiver of the method
		/*
			Consider a struct type T with two methods
				type T struct {
					a int
				}
				func (tv T) Mv(a int) int          { return 0 } // value receiver
				func (tp *T) Mp(f float32) float32 { return 1 } // pointer receiver
				var t T

			The expression
			T.Mv
			yields a function equivalent to Mv but with an explicit receiver as its first argument; it has signature
			func(tv T, a int) int

			so these five invocations are equivalent:
			t.Mv(7)
			T.Mv(t, 7)
			(T).Mv(t, 7)
			f1 := T.Mv; f1(t, 7)
			f2 := (T).Mv; f2(t, 7)

			Similarly, the expression
			(*T).Mp
			yields a function value representing Mp with signature
			func(tp *T, f float32) float32

			For a method with a value receiver, one can derive a function with an explicit pointer receiver, so
			(*T).Mv
			yields a function value representing Mv with signature
			func(tv *T, a int) int
			the method does not overwrite the value whose address is passed in the function call

			var f1 = (*T).Mv
			show("Value reciever, t.a before: ", t.a) // 0
			f1(&t, 42)
			show("t.a after: ", t.a) // 0
			var f2 = (*T).Mp
			f2(&t, 42) // t.Mp(42)
			show("t.a after: ", t.a) // 42

			... a value-receiver function for a pointer-receiver method, is illegal
		*/
	}
	methodExpressions()

	var methodValues = func() {
		show("Method values, x.M is called a method value")
		// If the expression `x` has static type `T` and `M` is in the method set of type `T`
		// The expression x is evaluated and saved during the evaluation of the method value; the saved copy is then used as the receiver

		/* e.g.
		   type S struct { *T }
		   type T int
		   func (t T) M() { print(t) }

		   t := new(T) // *T, reference to t
		   s := S{T: t}
		   f := t.M                    // receiver *t is evaluated and stored in f
		   g := s.M                    // receiver *(s.T) is evaluated and stored in g
		   *t = 42                     // does not affect stored receivers in f and g
		*/

		/* consider:
		type T struct {
			a int
		}
		func (tv T) Mv(a int) int          { tv.a = a; return tv.a } // value receiver
		func (tp *T) Mp(f float32) float32 { return 1 }              // pointer receiver
		var t T
		var pt *T
		func makeT() T { return T{} }

		// The expression
		t.Mv // yields a function value of type
		func(int) int

		// These two invocations are equivalent:
		t.Mv(7)
		f := t.Mv; f(7)

		// a reference to a non-interface method with a value receiver using a pointer
		// will automatically dereference that pointer:
		// pt.Mv is equivalent to (*pt).Mv.

		// a reference to a non-interface method with a pointer receiver using an addressable value
		// will automatically take the address of that value:
		// t.Mp is equivalent to (&t).Mp

		// e.g.
		f := t.Mv; f(7)   // like t.Mv(7)
		f := pt.Mp; f(7)  // like pt.Mp(7)
		f := pt.Mv; f(7)  // like (*pt).Mv(7)
		f := t.Mp; f(7)   // like (&t).Mp(7)
		f := makeT().Mp   // invalid: result of makeT() is not addressable
		*/

		//  it is also legal to create a method value from a value of interface type
	}
	methodValues()

	var indexExpressions = func() {
		show("Index expressions, A primary expression of the form `a[x]`")
		// denotes the element of the array, pointer to array, slice, string or map a indexed by x.
		// The value x is called the index or map key

		// array, slice, string:
		// index x is in range `0 <= x < len(a)`, otherwise it is out of range

		// For a of pointer to array type: `a[x]` is shorthand for `(*a)[x]`

		// string:
		// `a[x]` is the non-constant byte value at index `x` and the type of `a[x]` is `byte`
		// a[x] may not be assigned to

		// map:
		// x's type must be assignable to the key type of `M`
		// if the map is `nil` or does not contain such an entry, `a[x]` is the zero value for the element type of `M`
		var m map[string]int = nil
		var x = "bar"
		show("map: ", m)
		show("map[foo] element: ", m["foo"])

		// An index expression on a map a of type map[K]V used in an assignment statement or initialization of the special form
		// yields an additional untyped boolean value
		var v, keyInMap = m[x]
		v, keyInMap2 := m[x]
		v, keyInMap = m[x]
		show("map index expression: ", v, keyInMap, keyInMap2)
		// map index expression: int(0); bool(false); bool(false);

		// Assigning to an element of a nil map causes a run-time panic
		/*
			var m = map[string]int{} // empty map vs
			var m map[string]int = nil // nil map
			var m map[string]int // nil map
		*/
	}
	indexExpressions()

	var sliceExpressions = func() {
		show("Slice expressions construct a substring or slice from a string, array, pointer to array, or slice")
		// There are two variants: a simple form that specifies a low and high bound,
		// and a full form that also specifies a bound on the capacity
		var a = "foo bar"
		var low, high = 2, 5

		// Simple slice expressions

		// The primary expression `a[low : high]`
		// constructs a substring or slice.
		// The core type of a must be a string, array, pointer to array, slice, or a bytestring
		// The result has indices starting at 0 and length equal to high - low
		var b = a[low:high]
		show("slice: ", a, low, high, b, len(b))
		// slice: string(foo bar); int(2); int(5); string(o b); int(3);

		// For convenience, any of the indices may be omitted.
		// A missing low index defaults to zero; a missing high index defaults to the length
		b = a[2:] // same as a[2 : len(a)]
		b = a[:3] // same as a[0 : 3]
		b = a[:]  // same as a[0 : len(a)]

		// If a is a pointer to an array, `a[low : high]` is shorthand for `(*a)[low : high]`

		// For arrays or strings, the indices are in range `0 <= low <= high <= len(a)`, otherwise they are out of range
		// For slices, the upper index bound is the slice capacity `cap(a)` rather than the length

		// If the sliced operand is an array, it must be addressable

		// If the sliced operand of a valid slice expression is a `nil` slice, the result is a `nil` slice.
		func() {
			var a [10]int // array, ten zeros
			s1 := a[3:7]  // underlying array of s1 is array a; &s1[2] == &a[5]
			s2 := s1[1:4] // underlying array of s2 is underlying array of s1 which is array a; &s2[1] == &a[5]
			s2[1] = 42    // s2[1] == s1[2] == a[5] == 42; they all refer to the same underlying array element
			show("array: ", a)
			// array: [10]int([0 0 0 0 0 42 0 0 0 0]);

			var s []int // nil slice
			s3 := s[:0] // s3 == nil
			show("slice: ", s3, s3 == nil, s == nil)
			// slice: []int([]); bool(true); bool(true);
		}()

		// Full slice expressions

		func() {
			var a, low, high, max = [10]int{}, 2, 5, 7
			// The primary expression `a[low : high : max]`
			var b = a[low:high:max]
			// constructs a slice of the same type, and with the same length and elements as the simple slice expression `a[low : high]`
			// Additionally, it controls the resulting slice's capacity by setting it to `max - low`
			// Only the first index may be omitted; it defaults to 0
			// The indices are in range `0 <= low <= high <= max <= cap(a)`, otherwise they are out of range
			show("array: ", b, len(b), cap(b)) // array: []int([0 0 0]); int(3); int(5);
		}()
	}
	sliceExpressions()

	var typeAssertions = func() {
		show("Type assertions, The notation `x.(T)` is called a type assertion")
		// the primary expression
		// x.(T)
		// asserts that `x` is not `nil` and that the value stored in `x` is of type `T`.

		// If the type assertion holds, the value of the expression is the value stored in x and its type is T
		// If the type assertion is false, a run-time panic occurs

		var x interface{} = 7 // x has dynamic type int and value 7
		i := x.(int)          // i has type int and value 7
		_ = i == 7

		// type I interface { m() }
		// func f(y I) {
		// 	s := y.(string)        // illegal: string does not implement I (missing method m)
		// 	r := y.(io.Reader)     // r has type io.Reader and the dynamic type of y must implement both I and io.Reader
		// }

		// special form

		// A type assertion used in an assignment statement or initialization of the special form
		// yields an additional untyped boolean value.
		// The value of ok is true if the assertion holds.
		// Otherwise it is false and the value of v is the zero value for type T.
		// No run-time panic occurs in this case
		type T = int // type alias, types identical
		var v, ok = x.(T)
		v, ok = x.(T)
		v, ok2 := x.(T)
		// var v, ok interface{} = x.(T) // dynamic types of v and ok are T and bool
		show("asserted: ", v, ok, ok2)
		// asserted: int(7); bool(true); bool(true);
	}
	typeAssertions()

	var calls = func() {
		show("Calls, `f(a1, a2, … an)` calls f with arguments a1, a2, … an")
		// Given an expression `f` with a core type `F` of function type

		// arguments must be single-valued expressions assignable to the parameter types of F and
		// are evaluated before the function is called

		/*
			// method call vs function call from package
			math.Atan2(x, y)  // function call
			var pt *Point
			pt.Scale(3.5)     // method call with receiver pt
		*/

		// If f denotes a generic function, it must be instantiated before it can be called or used as a function value

		// In a function call, the function value and arguments are evaluated in the usual order.
		// After they are evaluated, the parameters of the call are passed by value to the function
		// and the called function begins execution.
		// The return parameters of the function are passed by value back to the caller when the function returns

		// Calling a nil function value causes a run-time panic

		// As a special case, if the return values of a `g` are individually assignable to the parameters of `f`,
		// then the call `f(g(parameters_of_g))` will invoke `f`
		// after binding the return values of `g` to the parameters of `f`
		var Split = func(s string, pos int) (string, string) {
			return s[0:pos], s[pos:]
		}
		var Join = func(s, t string) string {
			return s + t
		}
		if Join(Split("value", len("value")/2)) != "value" {
			panic("test fails")
		}

		// Methods: There is no distinct method type and there are no method literals
	}
	calls()

	var passingVariadicArguments = func() {
		show("Passing arguments to ... parameters")
		// If `f` is variadic with a final parameter `p` of type `...T`, then within `f` the type of `p` is equivalent to type `[]T`.
		// If `f` is invoked with no actual arguments for `p`, the value passed to `p` is `nil`

		// new slice/array for every call
		var f = func(xs ...int) {
			show("one argument as slice: ", xs == nil, xs, len(xs), cap(xs))
		}
		f()        // one argument as slice: bool(true); []int([]); int(0); int(0);
		f(1)       // one argument as slice: bool(false); []int([1]); int(1); int(1);
		f(1, 2, 3) // one argument as slice: bool(false); []int([1 2 3]); int(3); int(3);

		// If the final argument is assignable to a slice type `[]T` and is followed by `...`,
		// it is passed unchanged as the value for a `...T` parameter. In this case no new slice is created
		var ys = []int{3, 7}
		f(ys...) // one argument as slice: bool(false); []int([3 7]); int(2); int(2);
	}
	passingVariadicArguments()

	var instantiations = func() {
		show("A generic function or type is instantiated by substituting type arguments for the type parameters")
		// Instantiating a type results in a new non-generic named type;
		// instantiating a function produces a new non-generic function

		/*
			type parameter list    type arguments    after substitution
			[P any]                int               int satisfies any
			[S ~[]E, E any]        []int, int        []int satisfies ~[]int, int satisfies any
			[P io.Writer]          string            illegal: string doesn't satisfy io.Writer
			[P comparable]         any               any satisfies (but does not implement) comparable
		*/

		// For a generic type, all type arguments must always be provided explicitly.

		// When using a generic function, type arguments may be provided explicitly,
		// or they may be partially or completely inferred from the context
		/*
			// sum returns the sum (concatenation, for strings) of its arguments.
			func sum[T ~int | ~float64 | ~string](x... T) T { ... }

			x := sum                       // illegal: the type of x is unknown
			intSum := sum[int]             // intSum has type func(x... int) int
			a := intSum(2, 3)              // a has value 5 of type int
			b := sum[float64](2.0, 3)      // b has value 5.0 of type float64
			c := sum(b, -1)                // c has value 4.0 of type float64

			type sumFunc func(x... string) string
			var f sumFunc = sum            // same as var f sumFunc = sum[string]
			f = sum                        // same as f = sum[string]
		*/

		// A partial type argument list cannot be empty; at least the first argument must be present
		/*
			func apply[S ~[]E, E any](s S, f func(E) E) S { … }

			f0 := apply[]                  // illegal: type argument list cannot be empty
			f1 := apply[[]int]             // type argument for S explicitly provided, type argument for E inferred
			f2 := apply[[]string, string]  // both type arguments explicitly provided

			var bytes []byte
			r := apply(bytes, func(byte) byte { … })  // both type arguments inferred from the function arguments
		*/
	}
	instantiations()

	var typeInference = func() {
		show("Type inference, A use of a generic function may omit some or all type arguments if they can be inferred from the context")
		// Type inference uses the type relationships between pairs of types for inference
		// a function argument must be `assignable` to its respective function parameter

		// Each such pair of matched types corresponds to a `type equation` ...
		// Inferring the missing type arguments means solving the resulting set of type equations

		// Type equations are always solved for the bound type parameters only
		// The types of function arguments may contain type parameters from other functions (such as a generic function enclosing a function call).
		// Those type parameters may also appear in type equations but they are not `bound` in that context.
		// type parameters are called bound type parameters:
		// Given a set of type equations, the type parameters to solve for
		// are the type parameters of the functions that need to be instantiated and for which no explicit type arguments is provided

		// Type inference gives precedence to type information obtained from typed operands before considering untyped constants

		// the bound type parameters in each type argument are substituted with the respective type arguments for those type parameters
		// until each type argument is free of bound type parameters

		var typeUnification = func() {
			show("Type unification, Type inference solves type equations through type unification")
			// Type unification recursively compares the LHS and RHS types of an equation
			// where ... types may ... contain bound type parameters,
			// and looks for type arguments for those type parameters
			// such that the LHS and RHS match

			// Unification uses a combination of `exact` and `loose` unification depending on whether two types have to be
			// `identical`, `assignment`-compatible, or only `structurally` equal.

			// Type inference repeats type unification as long as new type arguments are inferred
		}
		typeUnification()
	}
	typeInference()

	var operators = func() {
		show("Operators, Operators combine operands into expressions")
		// ... the operand types must be identical unless ...

		// if one operand is an untyped constant
		// and the other operand is not, the constant is implicitly converted to the type of the other operand

		// The right operand in a shift expression must have integer type
		// or be an untyped constant representable by a value of type uint.
		// If the left operand of a non-constant shift expression is an untyped constant, it is first implicitly converted
		// to the type it would assume if the shift expression were replaced by its left operand alone

		// var a [1024]byte
		var s uint = 33

		// The results of the following examples are given for 64-bit ints.
		var i = 1 << s         // 1 has type int
		var j int32 = 1 << s   // 1 has type int32; j == 0
		var k = uint64(1 << s) // 1 has type uint64; k == 1<<33
		var m int = 1.0 << s   // 1.0 has type int; m == 1<<33
		var n = 1.0<<s == j    // 1.0 has type int32; n == true
		var o = 1<<s == 2<<s   // 1 and 2 have type int; o == false
		var p = 1<<s == 1<<33  // 1 has type int; p == true

		// var u = 1.0<<s                 // illegal: 1.0 has type float64, cannot shift
		// var u1 = 1.0<<s != 0           // illegal: 1.0 has type float64, cannot shift
		// var u2 = 1<<s != 1.0           // illegal: 1 has type float64, cannot shift
		// var v1 float32 = 1<<s          // illegal: 1 has type float32, cannot shift
		// var v2 = string(1<<s)          // illegal: 1 is converted to a string, cannot shift

		var w int64 = 1.0 << 33 // 1.0<<33 is a constant shift expression; w == 1<<33
		// var x = a[1.0<<s]            // panics: 1.0 has type int, but 1<<33 overflows array bounds
		// var b = make([]byte, 1.0<<s) // 1.0 has type int; len(b) == 1<<33 // oom
		show("shift operators, x64: ", i, j, k, m, n, o, p, w)
		// shift operators, x64: int(8589934592); int32(0); uint64(8589934592); int(8589934592); bool(true); bool(false); bool(true); int64(8589934592);

		// The results of the following examples are given for 32-bit ints,
		// which means the shifts will overflow.
		var mm int = 1.0 << s // 1.0 has type int; mm == 0
		var oo = 1<<s == 2<<s // 1 and 2 have type int; oo == true
		// var pp = 1<<s == 1<<33         // illegal: 1 has type int, but 1<<33 overflows int
		// var xx = a[1.0<<s]            // 1.0 has type int; xx == a[0] // panic
		// var bb = make([]byte, 1.0<<s) // 1.0 has type int; len(bb) == 0 // oom
		show("shift operators, x32: ", mm, oo)
		// shift operators, x32: int(8589934592); bool(false);

		var operatorPrecedence = func() {
			show("Operator precedence, Unary operators have the highest precedence")
			// As the `++` and `--` operators form statements, not expressions, they fall outside the operator hierarchy
			// statement `*p++` is the same as `(*p)++`

			// There are five precedence levels for binary operators
			// Multiplication operators bind strongest,
			// followed by addition operators, comparison operators,
			// `&&` (logical AND), and finally `||` (logical OR):
			/*
				Precedence    Operator
					5             *  /  %  <<  >>  &  &^
					4             +  -  |  ^
					3             ==  !=  <  <=  >  >=
					2             &&
					1             ||
			*/
			// Binary operators of the same precedence associate from left to right
		}
		operatorPrecedence()
	}
	operators()

	var arithmeticOperators = func() {
		show("Arithmetic operators apply to numeric values and yield a result of the same type as the first operand")
		// `+` also applies to strings
		// The bitwise logical and shift operators apply to integers only
		/*
			+    sum                    integers, floats, complex values, strings
			-    difference             integers, floats, complex values
			*    product                integers, floats, complex values
			/    quotient               integers, floats, complex values
			%    remainder              integers

			&    bitwise AND            integers
			|    bitwise OR             integers
			^    bitwise XOR            integers
			&^   bit clear (AND NOT)    integers

			<<   left shift             integer << integer >= 0
			>>   right shift            integer >> integer >= 0
		*/

		// The operands are represented as values of the type argument that the type parameter is instantiated with,
		// and the operation is computed with the precision of that type argument
		var v1, v2 = []float64{1, 2}, []float64{3, 4}
		show("operands and result of generic expression: ", sfs.DotProduct(v1, v2))
		// operands and result of generic expression: float64(11);

		var IntegerOperators = func() {
			show("Integer operators, quotient, reminder, shift, unari ops")
			// For two `integer` values x and y,
			// the integer quotient `q = x / y` and remainder `r = x % y` satisfy the following relationships
			// with `x / y` truncated towards zero:
			// `x = q*y + r`  and  `|r| < |y|`
			// With exception:
			// the quotient q `x / -1 == x` and r == 0 if the dividend `x` is the most negative value for its type,
			// due to two's-complement integer overflow

			// The shift operators implement arithmetic shifts if the left operand is a signed integer
			// and logical shifts if it is an unsigned integer
			// There is no upper limit on the shift count.
			// Shifts behave as if the left operand is shifted `n` times by 1 for a shift count of `n`

			// For integer operands, the unary operators +, -, and ^ are defined as follows
			/*
				+x                          is 0 + x
				-x    negation              is 0 - x
				^x    bitwise complement    is m ^ x  with m = "all bits set to 1" for unsigned x
													and  m = -1 for signed x
			*/
		}
		IntegerOperators()

		var integerOverflow = func() {
			show("Integer overflow, For unsigned integer values, the operations +, -, *, and << are computed modulo 2n")
			// discard high bits upon overflow, and programs may rely on "wrap around"
			var a uint8 = 250
			show("overflow unsigned: ", a+a) // overflow unsigned: uint8(244);

			// For signed integers, the operations +, -, *, /, and << may legally overflow
			// and the resulting value exists and is deterministically defined by the signed integer representation, the operation, and its operands.
			// Overflow does not cause a run-time panic.
			var b int8 = 127
			show("overflow signed: ", b+b) // overflow signed: int8(-2);

		}
		integerOverflow()

		var floatingPointOperators = func() {
			show("Floating-point operators, For floating-point and complex numbers, +x is the same as x, while -x is the negation of x.")
			// The result of a floating-point or complex division by zero ... is implementation-specific

			// An implementation may combine multiple floating-point operations into a single fused operation ...
			// produce a result that differs from the value obtained by executing and rounding the instructions individually
			// For instance, some architectures provide a "fused multiply and add" (FMA) instruction
			// that computes `x*y + z` without rounding the intermediate result `x*y`

			// examples of FMA

			var x, y, z, r, t float32 = 1.1, 3.3, 5.5, 0, 0
			var p *float32 = &t
			// FMA allowed for computing r, because x*y is not explicitly rounded:
			r = x*y + z
			r = z
			r += x * y
			t = x * y
			r = t + z
			*p = x * y
			r = *p + z
			// r  = x*y + float64(z)

			// FMA disallowed for computing r, because it would omit rounding of x*y:
			// An explicit floating-point type conversion rounds to the precision of the target type, preventing fusion
			// r  = float64(x*y) + z
			// r  = z; r += float64(x*y)
			// t  = float64(x*y); r = t + z

		}
		floatingPointOperators()

		var stringConcatenation = func() {
			show("String concatenation, strings can be concatenated using the + operator or the += assignment operator")
			// String addition creates a new string
			var c = 65
			s := "hi " + string(c)
			s += " and good bye"
			show("string concatenation: ", s) // string concatenation: string(hi A and good bye);
		}
		stringConcatenation()
	}
	arithmeticOperators()

	var comparisonOperators = func() {
		show("Comparison operators, compare two operands and yield an untyped boolean value")
		// A type is `strictly comparable` if it is comparable and not an interface type nor composed of interface types
		// Boolean, numeric, string, pointer, and channel types are strictly comparable

		// In any comparison, the first operand must be assignable to the type of the second operand, or vice versa
		/*
			==    equal
			!=    not equal
			<     less
			<=    less or equal
			>     greater
			>=    greater or equal
		*/

		// The equality operators == and != apply to operands of `comparable` types.
		// The ordering operators <, <=, >, and >= apply to operands of `ordered` types.
		/*
			Boolean types are comparable
			Integer types are comparable and ordered
			Floating-point types are comparable and ordered ... as defined by the IEEE-754 standard
			Complex types are comparable
			String types are comparable and ordered. Two string values are compared lexically byte-wise
			Pointer types are comparable. Two pointer values are equal if they point to the same variable or if both have value nil
			Channel types are comparable. Two channel values are equal if they were created by the same call to `make` or if both have value nil
			Interface types ... are comparable. Two interface values are equal if they have identical dynamic types and equal dynamic values or if both have value nil
			Struct types are comparable if all their field types are comparable ... in source order
			Array types are comparable if their array element types are comparable
			Type parameters are comparable if they are strictly comparable
		*/
		// Slice, map, and function types are not comparable
		// as a special case, a slice, map, or function value may be compared to the predeclared identifier `nil`

		// run-time panic
		// A comparison of two interface values with identical dynamic types causes a run-time panic if that type is not comparable

	}
	comparisonOperators()

	// Logical operators
	// Logical operators apply to boolean values
	// The right operand is evaluated conditionally
	/*
		&&    conditional AND    p && q  is  "if p then q else false"
		||    conditional OR     p || q  is  "if p then true else q"
		!     NOT                !p      is  "not p"
	*/

	var addressOperators = func() {
		show("Address operators, For an operand x of type T, the address operation `&x` generates a pointer of type `*T` to x")
		// The operand must be addressable, that is, either a
		// variable, pointer indirection, or slice indexing operation; or a field selector of an addressable struct operand;
		// or an array indexing operation of an addressable array.
		// As an exception to the addressability requirement, x may also be a (possibly parenthesized) composite literal.

		// If x is nil, an attempt to evaluate `*x` will cause a run-time panic
		// var x *int = nil
		// *x   // causes a run-time panic
		// &*x  // causes a run-time panic
	}
	addressOperators()

	var receiveOperator = func() {
		show("Receive operator, For an operand ch whose core type is a channel, the value of the receive operation <-ch is the value received from the channel ch")
		//  the type of the receive operation is the element type of the channel

		// The expression blocks until a value is available
		// Receiving from a nil channel blocks forever
		// A receive operation on a closed channel ... immediately, yielding the element type's zero value

		// A `receive expression` used in an assignment statement or initialization of the special form
		// `ok` is `false` if it is a zero value generated because the channel is closed and empty
		var ch = make(chan byte)
		close(ch)
		var x, ok = <-ch
		show("receive expression, value, was-sent-to-opened-channel: ", x, ok)
		// receive expression, value, was-sent-to-opened-channel: uint8(0); bool(false);
	}
	receiveOperator()

	var conversions = func() {
		show("Conversions, A conversion changes the type of an expression")
		// A conversion may appear literally in the source, or it may be implied by the context in which an expression appears
		// An explicit conversion is an expression of the form `T(x)` where T is a type and x is an expression

		// If the type starts with the operator `*` or `<-`,
		// or if the type starts with the keyword `func` and has no result list, it must be parenthesized
		/*
			*Point(p)        // same as *(Point(p))
			(*Point)(p)      // p is converted to *Point

			<-chan int(c)    // same as <-(chan int(c))
			(<-chan int)(c)  // c is converted to <-chan int

			func()(x)        // function signature func() x
			(func())(x)      // x is converted to func()
			(func() int)(x)  // x is converted to func() int
			func() int(x)    // x is converted to func() int (unambiguous)
		*/

		// Converting a constant to a type that is not a type parameter yields a typed constant
		// nil is not a constant
		/*
			int(1.2)                 // illegal: 1.2 cannot be represented as an int
			string(65.0)             // illegal: 65.0 is not an integer constant
		*/

		// Converting a constant to a type parameter yields a non-constant value of that type
		/*
			func f[P ~float32|~float64]() {
				… P(1.1) …
				// results in a non-constant value of type P and the value 1.1 is represented as a float32 or a float64 depending on the type argument
			}
		*/

		// A non-constant value x can be converted to type T in any of these cases
		// ...
		// x is an integer or a slice of bytes or runes and T is a string type
		// x is a string and T is a slice of bytes or runes

		// Struct tags are ignored when comparing struct types for identity for the purpose of conversion

		// Specific rules apply to (non-constant) conversions between numeric types to/from a string type
		// These conversions may change the representation of x and incur a run-time cost.
		// All other conversions only change the type but not the representation of x

		// There is no linguistic mechanism to convert between pointers and integers.
		// The package unsafe implements this functionality under restricted circumstances

		// Conversions between numeric types
		// Conversions to and from a string type
		// Conversions from slice to array or array pointer
	}
	conversions()

	// Constant expressions
	// Order of evaluation
}

func show(msg string, xs ...any) {
	var line string = msg
	for _, x := range xs {
		line += fmt.Sprintf("%T(%v); ", x, x) // type(value)
		// line += fmt.Sprintf("%#v; ", x) // repr
	}
	fmt.Println(line)
}
