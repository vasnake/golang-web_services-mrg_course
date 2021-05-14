package main

// not addresses, just references

func main() {
	a := 2
	b := &a // b ref to a

	*b = 3
	c := &a
	println("b, *b, a, c: ", b, *b, a, c) //  0xc000046770 3 3 0xc000046770

	// ref to anon value
	d := new(int)
	println("new int: d, *d", d, *d) // 0xc000046770 0

	*d = 12
	*c = *d                                     // a = 12
	*d = 13                                     // a = 12, anon = 13
	println("a, c, *c, d, *d", a, c, *c, d, *d) // 12 0xc000046768 12 0xc000046770 13

	c = d                                       // c ref to anon
	*c = 14                                     // anon = 14, a = 12
	println("a, c, *c, d, *d", a, c, *c, d, *d) // 12 0xc000046770 14 0xc000046770 14
}
