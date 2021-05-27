package main

func main() {
	myMap := map[string]string{"name": "Bender"}

	// simple if, only bool
	// if with init section
	// if-else

	// switch operator, no fallthrough by default, no need for `break`
	// `or` items in one `case` in switch
	// bool expressions in `case`
	// break from switch, break from outer loop

	boolVal := 1 > 0
	if boolVal {
		println("1 > 0")
	}

	if 0 < 1 {
		println("0 < 1")
	}

	if v, exists := myMap["name"]; exists {
		println("name:", v)
	}

	if len(myMap) == 0 {
		println("len 0")
	} else if len(myMap) == 2 {
		println("len 2")
	} else {
		println("len unknown:", len(myMap))
	}

	switch len(myMap) {
	case 0:
		println("len 0")
		fallthrough // go to next case
	case 2, 1: // multiple cases in one
		println("len 2 or 1")
	default:
		println("len unknown:", len(myMap))
	}

	switch {
	case 0 > 1 || 0 > 2:
		println("never")
	case 1 == 0 || 1 > 0:
		println("yep")
		fallthrough
	default:
		println("1 > 0 or 1 == 0? first")
	}

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
	}

Loop:
	for k, v := range myMap {
		println("switch in loop,", k, v)
		switch {
		case k == "name" && v == "Bender":
			println("breaking loop")
			break Loop
		default:
			println("k, v:", k, v)
		}
	}
}
