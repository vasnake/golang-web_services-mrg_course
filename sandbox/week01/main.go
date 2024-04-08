package main // week 1

import (
	"bufio"
	"fmt"
	"go/types"
	"io"
	"os"
	"reflect"
	"time"
	"unicode/utf8"
	"week01/interfaces"
	m "week01/methods"
	"week01/person"
)

func vars_1() {
	var num0 int // platform dependent int, 32 or 64 or ? with default value
	fmt.Println("1. default int: ", num0)

	var num1 int = 1 // set init value
	fmt.Println("2. with init value: ", num1)

	var num2 = 2 // type inference
	fmt.Println("3. detected type: ", num2, types.BasicInfo(num2), reflect.TypeOf(num2).String())
	// 3. detected type:  2 2 int
	fmt.Printf("4. type: %T\n", num2)
	// 4. type: int

	num := 3 // skip `var` and `type`
	fmt.Println("5. new var: ", num, reflect.TypeOf(num).String())
	// 5. new var:  3 int

	num += 1 // increment
	num++    // postfix increment
	fmt.Println("6. incremented: ", num)
	// 6. incremented:  5

	var a, b = 11, 12 // define multiple vars
	a, c := 13, 14    // only if `c` is new
	a, b, c = 1, 2, 3 // multi assign
	fmt.Println("7. multi: ", a, b, c)
	// 7. multi:  1 2 3
}

func vars_2() {
	var platformInt int = 1          // platform dependent size
	var autoInt = 2                  // compiler decision, should be platform int
	var bigInt int64 = (1 << 63) - 1 // you can select from int8 .. int64
	var uInt uint = 1 << 63          // unsigned, same as int
	fmt.Println("ints: ", platformInt, autoInt, bigInt, uInt)
	// ints:  1 2 9223372036854775807 9223372036854775808

	var pi float32 = 3.14 // explicit 32 (or 64)
	var e float64 = 2.718
	var _pi = 3.14          // platform dependent float
	var flotDefault float32 // zero default
	fmt.Println("floats: ", pi, e, _pi, flotDefault)
	// floats:  3.14 2.718 3.14 0

	var defaultBool bool // default 0 == false
	var b bool = true
	condition := 0 == 0 // condition := (0 == 0)
	fmt.Println("bools: ", defaultBool, b, condition)
	// bools:  false true true

	// complex
	var c complex64 = -1.1 + 7.12i
	complex := 1.2 + 3.4i // complex128 by default
	fmt.Println("complex: ", c, complex)
	// complex:  (-1.1+7.12i) (1.2+3.4i)
}

func stringsDemo() {
	var defaultStr string               // empty by default
	var interprStr = "CR: \n, TAB: \t." // will be parsed (interpreted)
	var rawStr = `CR: \n, TAB: \t.`     // as is, no interpretaton
	var utfStr = "ЯЫЧфйДжЮËъЭщП"        // all strings are utf8
	var oneByte byte = '\x27'           // single quotes for symbols, uint8 = 39
	var oneRune rune = 'Ы'              // int32 = 1067
	fmt.Println("strings: ", defaultStr, interprStr, rawStr, utfStr, oneRune, oneByte)
	// strings:   CR:
	// , TAB:  . CR: \n, TAB: \t. ЯЫЧфйДжЮËъЭщП 1067 39

	// concat
	fmt.Println("concat: ", rawStr+" wha?")
	// concat:  CR: \n, TAB: \t. wha?

	// immutable
	//rawStr[0] = 39 // invalid

	// string length
	byteLen := len(utfStr)
	runeLen := utf8.RuneCountInString(utfStr)
	fmt.Println("string len: ", byteLen, runeLen)
	// string len:  26 13

	// slice, view
	fmt.Println("byte slice: ", utfStr[:3], utfStr[:4], rawStr[1])
	// byte slice:  Я� ЯЫ 82
	bytes := []byte(utfStr)
	symbols := string(bytes[24:])
	fmt.Println("bytes<->string: ", bytes, symbols)
	// bytes<->string:  [208 175 208 171 208 167 209 132 208 185 208 148 208 182 208 174 195 139 209 138 208 173 209 137 208 159] П
}

func constDemo() {
	const pi = 3.1415 // type inferred, but not here, in application place

	// multiple constants
	const (
		e  = 2.718
		hi = "Hi"
	)

	// enum by iota
	const (
		zero = iota
		_    // skip `1`
		two
		three
	)

	// complex enum by iota
	const (
		_         = iota             // skip zero
		KB uint64 = 1 << (10 * iota) // 1024 == 1 << (10 * 1)
		MB                           // 1048576 == 1 << (10 * 2)
	)

	// untyped const
	const (
		year = 2021
	)

	fmt.Println("const: ", pi, e, hi, zero, three, two, KB, MB, 1<<(10*2))
	// const:  3.1415 2.718 Hi 0 3 2 1024 1048576 1048576

	// const type inferred in-place
	var bigY int64 = 2021
	var smallY uint16 = 2021
	fmt.Println("untyped const: ", smallY+year, bigY+year)
	// untyped const:  4042 4042
}

func typesDemo() {
	type UserID int // not int type, UserID type based on int, incompatible

	idx := 1
	var uid UserID = 42

	//uid = idx // invalid
	uid = UserID(idx) // simple cast
	println("idx, uid: ", idx, uid)
	// idx, uid:  1 1
}

func pointersDemo() {
	// pointers: just references, not really pointers (can't do pointers aryphmetic, see `unsafe` package)

	a := 2
	b := &a // b ref to a

	*b = 3  // write `3` to a
	c := &a // c ref to a
	println("b, *b, a, c: ", b, *b, a, c)
	// b, *b, a, c:  0xc00007cf20 3 3 0xc00007cf20

	// ref to anon value
	d := new(int) // `new`always returns ref to created value
	println("new int: d, *d", d, *d)
	// new int: d, *d 0xc00007cf28 0

	*d = 12
	*c = *d // a = 12
	*d = 13 // a = 12, anon = 13
	println("a, c, *c, d, *d", a, c, *c, d, *d)
	// a, c, *c, d, *d 12 0xc00007cf20 12 0xc00007cf28 13

	c = d   // c ref to anon
	*c = 14 // anon = 14, a = 12
	println("a, c, *c, d, *d", a, c, *c, d, *d)
	// a, c, *c, d, *d 12 0xc00007cf28 14 0xc00007cf28 14
}

func arrayDemo() {
	// array size defines array type, you can't use [2]int where [3]int is required

	var a1 [3]int
	fmt.Printf("a1 short %v, a1 full %#v \n", a1, a1)
	// a1 short [0 0 0], a1 full [3]int{0, 0, 0}

	// const can be used as size
	const size = 2
	var a2 [2 * size]bool
	fmt.Printf("a2 %#v \n", a2)
	// a2 [4]bool{false, false, false, false}

	// infer size from values enum
	a3 := [...]int{1, 2, 3}
	fmt.Printf("a3 %#v \n", a3)
	// a3 [3]int{1, 2, 3}

	// compile time panic
	//a3[size + 1] = 42
}

func sliceDemo1() {
	// creation, make function, default values
	var buf0 []int    // like array but w/o size, len:0, cap:0; nil
	buf1 := []int{}   // empty enumeration using literal
	buf2 := []int{42} // one elem, len:1, cap:1

	buf3 := make([]int, 0)     // len:0, cap:0
	buf4 := make([]int, 5)     // len:5, cap:5
	buf5 := make([]int, 5, 10) // len:5, cap:10

	fmt.Printf("buf0 %#v \n", buf0) //buf0 []int(nil)
	fmt.Printf("buf1 %#v \n", buf1) //buf1 []int{}
	fmt.Printf("buf2 %#v \n", buf2) //buf2 []int{42}
	fmt.Printf("buf3 %#v \n", buf3) //buf3 []int{}
	fmt.Printf("buf4 %#v \n", buf4) //buf4 []int{0, 0, 0, 0, 0}
	fmt.Printf("buf5 %#v \n", buf5) //buf5 []int{0, 0, 0, 0, 0}

	someInt := buf2[0]
	buf2[0] = 37
	fmt.Printf("buf2 %#v, oldval %#v \n", buf2, someInt) // buf2 []int{37}, oldval 42

	// append elements
	var buf []int            // len:0, cap:0
	buf = append(buf, 9, 10) // len:2, cap:2
	buf = append(buf, 12)    // len:3, cap:4, capacity x2 every reallocation
	// old storage deleted if unused elsewhere, old reference points to previous buffer
	fmt.Printf("buf %#v \n", buf) // buf []int{9, 10, 12}

	// append elements from another slice
	otherBuf := make([]int, 3)     // [0, 0, 0]
	buf = append(buf, otherBuf...) // unpack `from` using `...`
	fmt.Printf("buf %#v \n", buf)  // buf []int{9, 10, 12, 0, 0, 0}

	// get len, cap from slice
	fmt.Printf("buf len: %#v, cap: %#v \n", len(buf), cap(buf)) // buf len: 6, cap: 8
}

func sliceDemo2() {
	// slice operation, buf[start, stop]
	// slice op gives reference to the buffer

	buf := []int{1, 2, 3, 4, 5}
	fmt.Printf("buf %#v, len: %#v, cap: %#v \n", buf, len(buf), cap(buf))
	//buf []int{1, 2, 3, 4, 5}, len: 5, cap: 5

	sl1 := buf[1:4] // idx from 0, inclusive
	sl2 := buf[:2]
	sl3 := buf[2:]

	fmt.Printf("sl1 %#v, len: %#v, cap: %#v \n", sl1, len(sl1), cap(sl1))
	fmt.Printf("sl2 %#v, len: %#v, cap: %#v \n", sl2, len(sl2), cap(sl2))
	fmt.Printf("sl3 %#v, len: %#v, cap: %#v \n", sl3, len(sl3), cap(sl3))
	//sl1 []int{2, 3, 4}, len: 3, cap: 4
	//sl2 []int{1, 2}, len: 2, cap: 5
	//sl3 []int{3, 4, 5}, len: 3, cap: 3

	newBuf := buf[:] // same storage
	newBuf[0] = 9    // put value to storage
	fmt.Printf("buf %#v, len: %#v, cap: %#v \n", buf, len(buf), cap(buf))
	// buf []int{9, 2, 3, 4, 5}, len: 5, cap: 5

	// reference detached if buf was reallocated (e.g. using append)
	newBuf = append(newBuf, 6) // got a new storage
	newBuf[0] = 1
	fmt.Printf("newBuf %#v, len: %#v, cap: %#v \n", newBuf, len(newBuf), cap(newBuf))
	//newBuf []int{1, 2, 3, 4, 5, 6}, len: 6, cap: 10
	fmt.Printf("buf %#v, len: %#v, cap: %#v \n", buf, len(buf), cap(buf))
	//buf []int{9, 2, 3, 4, 5}, len: 5, cap: 5

	// copy(newBuf, existingBuf) checks lengths and copy only min(len1, len2) elements

	var emptyBuf []int
	numCopied := copy(emptyBuf, buf) // wrong
	println(numCopied)               // 0, because emptyBuf len=0

	emptyBuf = make([]int, len(buf), len(buf))
	numCopied = copy(emptyBuf, buf) // ok
	println(numCopied)              // 5

	// copy slice to slice, replacing existing values in buf
	buf = []int{1, 2, 3, 4}
	copy(buf[1:3], []int{5, 6})
	fmt.Printf("buf %#v, len: %#v, cap: %#v \n", buf, len(buf), cap(buf))
	// buf []int{1, 5, 6, 4}, len: 4, cap: 4
}

func mapDemo() {
	// hash table, associative array, keys unordered

	// creation, literal
	var user map[string]string = map[string]string{
		"name":     "Bart",
		"lastName": "Simpson",
	}
	fmt.Printf("map %#v, len: %#v\n", user, len(user))
	// map map[string]string{"lastName":"Simpson", "name":"Bart"}, len: 2

	// creation, `make` function
	profile := make(map[string]string, 10) // cap = 10
	fmt.Printf("map %#v, len: %#v \n", profile, len(profile))
	// map map[string]string{}, len: 0

	// element absence = element type default value
	// solution: val, exists = map[key]
	// `_` as blank var

	name := user["name"]        // Bart
	mName := user["middleName"] // nothing, default value ""
	println(name, mName)        // Bart

	mName, mNameExists := user["middleName"]
	_, mNameExists = user["middleName"] // only existence flag
	println(name, mNameExists)          // Bart false

	// function delete(map, key)
	delete(user, "lastName")
	fmt.Printf("map %#v, len: %#v\n", user, len(user))
	// map map[string]string{"name":"Bart"}, len: 1
}

func controlDemo() {
	myMap := map[string]string{"name": "Bender"}

	// simple if, only bool type
	// if with init section
	// if-else

	// switch operator, no fallthrough by default, no need for `break`
	// `or` items in one `case` in switch
	// bool expressions in `case`
	// break from switch, break from outer loop, labels

	boolVal := 1 > 0 // true
	if boolVal {
		println("1 > 0")
	} // 1 > 0

	if 0 < 1 {
		println("0 < 1")
	} // 0 < 1

	if v, exists := myMap["name"]; exists {
		println("name:", v)
	} // name: Bender

	if len(myMap) == 0 {
		println("len 0")
	} else if len(myMap) == 2 {
		println("len 2")
	} else {
		println("len unknown:", len(myMap))
	} // len unknown: 1

	switch len(myMap) {
	case 0:
		println("len 0")
		fallthrough // go to next case
	case 2, 1: // multiple cases in one
		println("len 2 or 1")
	default:
		println("len unknown:", len(myMap))
	} // len 2 or 1

	switch {
	case 0 > 1 || 0 > 2:
		println("never")
	case 1 == 0 || 1 > 0:
		println("yep")
		fallthrough
	default:
		println("1 > 0 or 1 == 0? first")
	}
	// yep
	// 1 > 0 or 1 == 0? first

	switch len(myMap) {
	case 0:
		println("len 0")
	case 1, 2:
		if len(myMap) == 1 {
			println("don't tell it's 1")
			break
		}
		println("len 1")
	default:
		println("len unknown:", len(myMap))
	} // don't tell it's 1

Loop: // mark block of code
	for k, v := range myMap {
		println("switch in loop,", k, v)
		switch {
		case k == "name" && v == "Bender":
			println("breaking loop")
			break Loop // n.b. it's not goto
		default:
			println("k, v:", k, v)
		}
	}
	// switch in loop, name Bender
	// breaking the loop
}

func loopDemo() {
	// only one keyword: `for`, in many forms
	// while(true): for { ... break }
	// while(condition): for cond { ... cond = ??? }
	// classic: for i := 0; i < cnt; i++ { ??? }
	// slice iteration: `for idx, val := range mySlice { ??? }`
	// map iteration, the same as slice
	// string iteration `for pos, rune := range myStr { ??? }` // not byte but rune!

	for {
		println("while(true)")
		break
	} // while(true)

	cond := true
	for cond {
		println("while condition is true")
		cond = false
	} // while condition is true

	for i := 0; i < 3; i++ {
		if i == 1 {
			continue
		}
		print("i:", i, ",")
	} // i:0,i:2,
	println()

	mySlice := []int{1, 2, 3}

	for i := range mySlice {
		println("i:", i, "v:", mySlice[i]) // index iteration
	}
	// i: 0 v: 1
	// i: 1 v: 2
	// i: 2 v: 3

	for i, v := range mySlice {
		println("i:", i, "v:", v) // index and value iteration
	}
	// i: 0 v: 1
	// i: 1 v: 2
	// i: 2 v: 3

	myMap := map[string]string{"name": "Bender"}

	for k, v := range myMap {
		println("k:", k, "v:", v)
	} // k: name v: Bender

	for k := range myMap {
		println("k:", k, "v:", myMap[k]) // w/o creating value var on each iteration
	} // k: name v: Bender

	for _, v := range myMap {
		println("v:", v)
	} // v: Bender

	myStr := "ЧЯДНТ"
	println("str len, bytes: ", len(myStr)) // str len, bytes: 10
	for pos, symb := range myStr {
		println("byte pos:", pos, ", rune as int32:", symb, ", rune as str:", string(symb))
	}
	// byte pos: 0 , rune as int32: 1063 , rune as str: Ч
	// byte pos: 2 , rune as int32: 1071 , rune as str: Я
	// byte pos: 4 , rune as int32: 1044 , rune as str: Д
	// byte pos: 6 , rune as int32: 1053 , rune as str: Н
	// byte pos: 8 , rune as int32: 1058 , rune as str: Т
}

func functionsDemo() {

	// simple definition
	// multi-parameter definition
	// named return var
	// return tuple
	// return named vars as tuple // beware default values of uninitialized vars
	// unpacked list of parameters, triple-dot notation

	// func sqrt(in int) int {
	// 	return in * in
	// }
	// this syntax is valid only on the package level

	var sqrt = func(in int) int {
		return in * in
	}

	var sum3 = func(a, b int, c int) int { // n.b. a and b are the same type
		return a + b + c
	}

	var namedReturn = func() (out int) {
		// out = zero value here

		out = 42
		return
		// return 37 // also will work
	}

	var withError = func(in int) (int, error) {
		if in > 1 {
			return 1, fmt.Errorf("shit happens")
		}
		return 0, nil
	}

	var withErrorNamed = func(cond bool) (res int, err error) {
		// res = 0, err = nil

		if cond {
			err = fmt.Errorf("shit happens")
			return // n.b. res will be 0 by default
			// return 0, fmt.Errorf("???")
		}
		return 42, nil
	}

	println("shut up compiler, I'm using these functions: ", sqrt, sum3, namedReturn, withError, withErrorNamed)
	// shut up compiler, I'm using these functions:  0x4e1f40 0x4e1f48 0x4e1f50 0x4e1f58 0x4e1f60

	var sum = func(in ...int) (res int) {
		fmt.Printf("in: %#v \n", in) // slice in: []int{1, 2, 3, 4}
		for _, v := range in {
			res += v
		}
		return
	}

	/*
		in: []int{1, 2, 3, 4}
		10:  10

		in: []int{1, 2, 3, 4}
		10:  10
	*/
	println("10: ", sum(1, 2, 3, 4))

	xs := []int{1, 2, 3, 4}
	println("10: ", sum(xs...)) // unpack slice
}

func firstclassFunctions() {
	// like any other var.
	// named function vs anon. function.
	// assign anon.func to var and call.
	// define var type as some function(in)out.
	// using closure as currying substitute.

	// named function
	var doNothing = func() {
		fmt.Println("I'm a regular function")
	}
	doNothing() // I'm a regular function

	// anonymous function, creation and call
	func(in string) {
		fmt.Println("anon func, in: ", in)
	}("no name function")
	// anon func, in:  no name function

	// assign func to var
	printer := func(in string) {
		fmt.Println("printer input:", in)
	}
	printer("ref to function saved in var")
	// printer input: ref to function saved in var

	// define func type
	type stringPrinterType func(msg string)

	// define function which take a function and call it
	worker := func(callback stringPrinterType) {
		callback("callback called by worker")
	}
	worker(printer)
	// printer input: callback called by worker

	// closure, access to a var outside of a function body
	prefixer := func(prefix string) stringPrinterType { // prefixer returns a function constructed using passed prefix
		return func(in string) {
			fmt.Printf("[%s] %s\n", prefix, in) // n.b. `prefix` linked to outer var
		}
	}
	successLogger := prefixer("SUCCESS")         // created a function
	successLogger("should be marked as success") // call a function
	// [SUCCESS] should be marked as success
}

func deferDemo() {
	var doStuff = func(in string) string {
		fmt.Printf("doStuff: `%s`\n", in)
		return fmt.Sprintf("stuff `%s` len: %d bytes", in, len(in))
	}

	// defer put call on stack, actual call will be made before exit from block
	defer fmt.Println("Work is done.")
	defer fmt.Println("closing ...") // n.b. closing called before `Work is done.`

	defer fmt.Println(doStuff("before work ...")) // n.b. parameters are evaluated right here, doStuff called right now!
	// doStuff: `before work ...` // but print "stuff `before work ...` len: 15 bytes" is deferred

	defer doStuff("after work ...")

	// anon functions let you to defer params evaluation
	defer func() {
		fmt.Println(doStuff("first deferred call")) // doStuff NOT called right here, it will be deferred
	}()

	fmt.Println("Working ...")
	// Working ...
	time.Sleep(3 * time.Millisecond)

	fmt.Println("Work is done, now deferred steps ...")
	// Work is done, now deferred steps ...

	/*
	   doStuff: `before work ...`

	   Working ...
	   Work is done, now deferred steps ...

	   doStuff: `first deferred call`
	   stuff `first deferred call` len: 19 bytes

	   doStuff: `after work ...`

	   stuff `before work ...` len: 15 bytes

	   closing ...
	   Work is done.
	*/
}

func recoverDemo() {
	// recover from panic using defer

	// panic: program stop with stacktrace.
	// panic .. recover: is not a try .. catch, don't use it this way.
	// you should catch panic if your program is a daemon of some sort, and it should continue processing another requests.

	var panicButton = func() {
		// catch panic event on exit and suppress it
		defer func() {
			// process panic event if there was some shit
			if err := recover(); err != nil {
				fmt.Println("panic detected:", err)
			}
		}()

		_, err := fmt.Println("working ...")
		if err == nil {
			panic("some shit happened! (not)")
		}
	}

	panicButton()
	println("Panic? I can't feel any panic, working as usual ...")
	/*
	   working ...
	   panic detected: some shit happened! (not)
	   Panic? I can't feel any panic, working as usual ...
	*/
}

func structsDemo() {
	// struct is a type

	type Person struct {
		Id      int
		Name    string
		Address string
	}

	type Account struct {
		Id      int
		Name    string
		Cleaner func(string) string // func as struct fiekd
		Owner   Person              // struct as field, not a composition
	}

	// PersonalAccount shows structs composition: add Person fields to Account
	type PersonalAccount struct {
		Id      int
		Name    string
		Cleaner func(string) string
		Owner   Person
		Person  // composition: Account id, name and address
	}

	// struct initialization, named fields init
	var acc Account = Account{
		Id:   42,
		Name: "foo",
	} // remember: you'll have default values for unset fields

	fmt.Printf("Account w/o owner: %#v\n", acc) // repr
	// Account w/o owner: main.Account{Id:42, Name:"foo", Cleaner:(func(string) string)(nil), Owner:main.Person{Id:0, Name:"", Address:""}}

	// short init, all fields present and accounted for
	acc.Owner = Person{33, "Foo Bar", "Under The Bridge"}
	fmt.Printf("Account with owner: %#v\n", acc)
	// Account with owner: main.Account{Id:42, Name:"foo", Cleaner:(func(string) string)(nil), Owner:main.Person{Id:33, Name:"Foo Bar", Address:"Under The Bridge"}}

	// composed struct
	pa := PersonalAccount{Id: 77, Name: "baz", Person: Person{Address: "Far Far Away", Name: "zab"}}
	fmt.Printf("Account with owner: %#v\n", pa)
	// Account with owner: main.PersonalAccount{Id:77, Name:"baz", Cleaner:(func(string) string)(nil), Owner:main.Person{Id:0, Name:"", Address:""}, Person:main.Person{Id:0, Name:"zab", Address:"Far Far Away"}}

	// PA.Person fields, n.b. w/o `Person` prefix, using Account namespace
	fmt.Printf("address: %#v\n", pa.Address)
	// address: "Far Far Away"

	// what about `Name`? outer field selected
	fmt.Printf("PA.Name: %#v\n", pa.Name)
	// PA.Name: "baz"
	fmt.Printf("PA.Person.Name: %#v\n", pa.Person.Name)
	// PA.Person.Name: "zab"
}

func methodsDemo() {
	p := m.Person{Id: 1, Name: "Foo"} // Address = ""

	p.CantSetName("Bar") // name is set in a copy of person
	fmt.Printf("person: %#v\n", p)
	// person: methods.Person{Id:1, Name:"Foo", Address:""}

	p.SetName("Bar")
	fmt.Printf("person: %#v\n", p)
	// person: methods.Person{Id:1, Name:"Bar", Address:""}
	// N.B. compiler automagically treated `p` as reference, converting call to smth like:
	(&p).SetName("Baz")
	fmt.Printf("person: %#v\n", p)
	// person: methods.Person{Id:1, Name:"Baz", Address:""}

	// alternatively, w/o magic
	pRef := new(m.Person) // return a reference
	pRef.SetName("Quix")
	fmt.Printf("person: %#v\n", pRef)
	// person: &methods.Person{Id:0, Name:"Quix", Address:""}
	// N.B. `&` symbol in `&methods.Person`

	// nested structures (composed), outer structure have all methods defined for inner structure
	acc := m.Account{}
	acc.SetName("Qux") // account.person.name
	fmt.Printf("account: %#v\n", acc)
	// account: methods.Account{Id:0, Name:"", No:0x0, Person:methods.Person{Id:0, Name:"Qux", Address:""}}

	// if method name is the same for outer and inner structs, outer struct have priority
	acc.SetNameDup("Quuz")
	fmt.Printf("account: %#v\n", acc)
	// account: methods.Account{Id:0, Name:"Quuz", No:0x0, Person:methods.Person{Id:0, Name:"Qux", Address:""}}

	acc.Person.SetNameDup("Corge")
	fmt.Printf("account: %#v\n", acc)
	// account: methods.Account{Id:0, Name:"Quuz", No:0x0, Person:methods.Person{Id:0, Name:"Corge", Address:""}}

	// another type methods
	s := m.MySlice([]int{1, 2})
	_s := s.Append(3)
	fmt.Printf("slice %#v, size %#v, another ref %#v\n", s, s.Size(), _s)
	// slice methods.MySlice{1, 2, 3}, size 3, another ref &methods.MySlice{1, 2, 3}
}

func packageDemo() {
	// package `person` code contains in a few files

	p := person.NewPerson(1, "Foo", "bar")
	fmt.Printf("main, person %+v\n", p)
	//main, person &{ID:1 Name:Foo secret:bar}

	// p.secret is private, only Capital letters are exported
	fmt.Println("main, secret", person.GetSecret(p))
	// main, secret bar
}

func interfaceBasic() {
	// uses interface, don't give a fuck about data type
	var Buy = func(p interfaces.Payer) {
		// any object with Payer implementation will do, e.g. Wallet
		err := p.Pay(10)
		if err != nil {
			panic(err)
		}
		fmt.Printf("You paid with %T\n", p)
	}

	w := &interfaces.Wallet{Cash: 100} // n.b. reference to wallet, so it can mutate Cash
	Buy(w)
	// You paid with *interfaces.Wallet
}

func interfaceMany() {
	// uses interface, don't give a fuck about data type
	var Buy = func(p interfaces.Payer) {
		// any object with Payer implementation will do, e.g. Wallet
		err := p.Pay(10)
		if err != nil {
			panic(err)
		}
		fmt.Printf("You paid with %T\n", p)
	}

	// buy with any kind of payable object
	w := &interfaces.Wallet{Cash: 100}
	Buy(w)

	var p interfaces.Payer
	p = &interfaces.Card{Balance: 100}
	Buy(p)

	p = w
	Buy(p)

	/*
	   You paid with *interfaces.Wallet
	   You paid with *interfaces.Card
	   You paid with *interfaces.Wallet
	*/
}

func interfaceCast() {
	var Buy = func(p interfaces.Payer) {
		// if you want to check type of the object or use custom logic, match object type
		switch p.(type) {
		case *interfaces.Wallet:
			fmt.Println("cash")
		case *interfaces.Card:
			card, confirmed := p.(*interfaces.Card) // type assertion
			if confirmed {                          // not much sense in doing it here, but for demo it's good enough
				fmt.Println("it is a card, no BS")
			}
			fmt.Println("prepare to authorize payment,", card.CardHolder)
		default:
			fmt.Println("unknown payable object", p)
		}

		err := p.Pay(10)
		if err != nil {
			panic(err)
		}
		fmt.Printf("You paid with %T\n", p)
	}

	Buy(&interfaces.Card{Balance: 11, CardHolder: "Foo"})
	/*
		it is a card, no BS
		prepare to authorize payment, Foo
		You paid with *interfaces.Card
	*/
}

func interfaceEmpty() {
	w := &interfaces.Wallet{Cash: 100}

	// printf accepts any empty interface
	// empty interface don't impose any limitations, it's just a tuple (type, value)

	fmt.Printf("repr: %#v\n", w)
	// repr: &interfaces.Wallet{Cash:100}

	// fmt.Printf("Stringer String: %s\n", w)
	// Stringer String: &{%!s(int=100)}

	fmt.Printf("Stringer String: %s\n", &interfaces.WalletNice{Cash: 99})
	// Stringer String: Wallet with 99 currency points

	var Buy = func(in interface{}) {
		// empty interface, dynamic type check

		var p interfaces.Payer
		var ok bool

		// type assertion
		if p, ok = in.(interfaces.Payer); !ok {
			fmt.Printf("%T is not a Payer\n", in)
			return
		}

		err := p.Pay(10)
		if err != nil {
			fmt.Printf("Paying error, %T %v\n", p, err)
			return
		}

		fmt.Printf("You paid with %T\n", p)
	}

	Buy(&interfaces.Wallet{Cash: 100}) // n.b. reference, it's vital for type cast
	// You paid with *interfaces.Wallet

	Buy([]int{1, 2, 3})
	// []int is not a Payer

	Buy(3.14)
	// float64 is not a Payer
}

func interfaceEmbed() {

	type Payer interface {
		Pay(int) error
	}

	type Ringer interface {
		Ring(string) error
	}

	type NFCPhone interface {
		// combined interface (embedded)
		Payer
		Ringer
	}

	var PayAndRing = func(phone NFCPhone) {
		// func wants phone with methods: Pay, Ring

		err := phone.Pay(1)
		if err != nil {
			fmt.Printf("paying error %v\n", err)
			return
		}

		fmt.Printf("payed via %T\n", phone)

		err = phone.Ring("3 times")
		if err != nil {
			fmt.Printf("ringing error %v\n", err)
			return
		}

	}

	p := &interfaces.Phone{Money: 9} // implements Pay, Ring methods
	PayAndRing(p)
	/*
	   pay 1 coins
	   payed via *interfaces.Phone
	   ring 3 times
	*/
}

func uniqueProgram() {
	// cat data_map.txt | go run uniq || exit

	var naiveUnique = func() {
		in := bufio.NewScanner(os.Stdin)
		seenAlready := make(map[string]bool)

		for in.Scan() {
			txt := in.Text()
			fmt.Printf("input line: `%s`", txt)

			if _, found := seenAlready[txt]; found {
				fmt.Printf("\n")
				continue
			}

			seenAlready[txt] = true
			fmt.Printf(" is unique \n")
		}
	}
	println("I'm using it:", naiveUnique)

	var sortedInputUnique = func() {
		in := bufio.NewScanner(os.Stdin)
		var prev string // empty string

		for in.Scan() {
			txt := in.Text()
			fmt.Printf("input line: `%s`", txt)

			if txt == prev {
				fmt.Printf("\n")
				continue
			}

			if txt < prev {
				panic("input is not sorted")
			}

			prev = txt
			fmt.Printf(" is unique \n")
		}
	}
	println("I'm using it:", sortedInputUnique)

	println("Start typing lines of text ...")
	// naiveUnique()
	sortedInputUnique()
}

// Logic can be described as: drop duplicates from an already sorted input.
func sortedInputUnique(input io.Reader, output io.Writer) error {
	in := bufio.NewScanner(input)
	var prev string

	for in.Scan() {
		txt := in.Text()
		fmt.Printf("input line: `%s`", txt)

		if txt == prev {
			fmt.Printf(" was seen already\n")
			continue
		}

		if txt < prev {
			return fmt.Errorf("input is not sorted")
		}

		prev = txt
		fmt.Printf(" is unique\n")

		fmt.Fprintln(output, txt)
	}

	return nil
}

func uniqueFuncTestable() {
	// cat data_map.txt | sort | go run uniq || exit
	err := sortedInputUnique(os.Stdin, os.Stdout)
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	vars_1()
	// vars_2()
	// stringsDemo()
	// constDemo()
	// typesDemo()
	// pointersDemo()
	// arrayDemo()
	// sliceDemo1()
	// sliceDemo2()
	// mapDemo()
	// controlDemo()
	// loopDemo()
	// functionsDemo()
	// firstclassFunctions()
	// deferDemo()
	// recoverDemo()
	// structsDemo()
	// methodsDemo()
	// packageDemo()
	// interfaceBasic()
	// interfaceMany()
	// interfaceCast()
	// interfaceEmpty()
	// interfaceEmbed()
	// uniqueProgram()
	// uniqueFuncTestable()
}
