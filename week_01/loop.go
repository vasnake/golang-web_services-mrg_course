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

/*
	// цикл без условия, while(true) OR for(;;;)
	for {
		fmt.Println("loop iteration")
		break
	}

	// цикл без условия, while(isRun)
	isRun := true
	for isRun {
		fmt.Println("loop iteration with condition")
		isRun = false
	}

	// цикл с условие и блоком инициализации
	for i := 0; i < 2; i++ {
		fmt.Println("loop iteration", i)
		if i == 1 {
			continue
		}
	}

	// операции по slice
	sl := []int{1, 2, 3}
	idx := 0

	for idx < len(sl) {
		fmt.Println("while-stype loop, idx:", idx, "value:", sl[idx])
		idx++
	}

	for i := 0; i < len(sl); i++ {
		fmt.Println("c-style loop", i, sl[i])
	}
	for idx := range sl {
		fmt.Println("range slice by index", sl[idx])
	}
	for idx, val := range sl {
		fmt.Println("range slice by idx-value", idx, val)
	}

	// операции по map
	profile := map[int]string{1: "Vasily", 2: "Romanov"}

	for key := range profile {
		fmt.Println("range map by key", key)
	}

	for key, val := range profile {
		fmt.Println("range map by key-val", key, val)
	}

	for _, val := range profile {
		fmt.Println("range map by val", val)
	}

	str := "Привет, Мир!"
	for pos, char := range str {
		fmt.Printf("%#U at pos %d\n", char, pos)
	}

}

*/
