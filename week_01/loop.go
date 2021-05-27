package main

func main() {
	// only one keyword: `for`, in many forms
	// while(true): for { ... break }
	// while(condition): for cond { ... cond = ??? }
	// classic: for i := 0; i < cnt; i++ { ??? }
	// slice iteration: `for idx, val := range mySlice { ??? }`
	// map iteration, the same as slice
	// string iteration `for pos, rune := range myStr { ??? }` // not byte!

	for {
		println("while(true)")
		break
	}

	cond := true
	for cond {
		println("while condition is true")
		cond = false
	}

	for i := 0; i < 3; i++ {
		if i == 1 {
			continue
		}
		println("i:", i)
	}

	mySlice := []int{1, 2, 3}

	for i := range mySlice {
		println("i:", i, "v:", mySlice[i]) // only index on stack
	}
	for i, v := range mySlice {
		println("i:", i, "v:", v) // index and value on stack
	}

	myMap := map[string]string{"name": "Bender"}

	for k, v := range myMap {
		println("k:", k, "v:", v)
	}
	for k := range myMap {
		println("k:", k, "v:", myMap[k]) // w/o creating value var on each iteration
	}
	for _, v := range myMap {
		println("v:", v)
	}

	myStr := "ЧЯДНТ"
	println("str len", len(myStr))
	for pos, symb := range myStr {
		println("byte pos:", pos, ", rune as int32:", symb, ", rune as str:", string(symb))
	}
}
