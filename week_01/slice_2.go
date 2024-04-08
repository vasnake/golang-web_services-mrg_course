package main

import "fmt"

func main() {
	// slice operation, buf[start, stop]
	// slice op gives reference to buffer memory

	buf := []int{1, 2, 3, 4, 5}
	fmt.Printf("buf %#v, len: %#v, cap: %#v \n", buf, len(buf), cap(buf))
	//buf []int{1, 2, 3, 4, 5}, len: 5, cap: 5

	sl1 := buf[1:4]
	sl2 := buf[:2]
	sl3 := buf[2:]

	fmt.Printf("sl1 %#v, len: %#v, cap: %#v \n", sl1, len(sl1), cap(sl1))
	fmt.Printf("sl2 %#v, len: %#v, cap: %#v \n", sl2, len(sl2), cap(sl2))
	fmt.Printf("sl3 %#v, len: %#v, cap: %#v \n", sl3, len(sl3), cap(sl3))
	//sl1 []int{2, 3, 4}, len: 3, cap: 4
	//sl2 []int{1, 2}, len: 2, cap: 5
	//sl3 []int{3, 4, 5}, len: 3, cap: 3

	newBuf := buf[:] // same storage
	newBuf[0] = 9
	fmt.Printf("buf %#v, len: %#v, cap: %#v \n", buf, len(buf), cap(buf))
	// buf []int{9, 2, 3, 4, 5}, len: 5, cap: 5

	// reference detached if buf was reallocated (e.g. using append)
	newBuf = append(newBuf, 6)
	newBuf[0] = 1
	fmt.Printf("newBuf %#v, len: %#v, cap: %#v \n", newBuf, len(newBuf), cap(newBuf))
	fmt.Printf("buf %#v, len: %#v, cap: %#v \n", buf, len(buf), cap(buf))
	//newBuf []int{1, 2, 3, 4, 5, 6}, len: 6, cap: 10
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

/*
	buf := []int{1, 2, 3, 4, 5}
	fmt.Println(buf)

	// получение среза, указывающего на ту же память
	sl1 := buf[1:4] // [2, 3, 4]
	sl2 := buf[:2]  // [1, 2]
	sl3 := buf[2:]  // [3, 4, 5]
	fmt.Println(sl1, sl2, sl3)

	newBuf := buf[:] // [1, 2, 3, 4, 5]
	// buf = [9, 2, 3, 4, 5], т.к. та же память
	newBuf[0] = 9

	// newBuf теперь указывает на другие данные
	newBuf = append(newBuf, 6)

	// buf    = [9, 2, 3, 4, 5], не изменился
	// newBuf = [1, 2, 3, 4, 5, 6], изменился
	newBuf[0] = 1
	fmt.Println("buf", buf)
	fmt.Println("newBuf", newBuf)

	// копирование одного слайса в другой
	var emptyBuf []int // len=0, cap=0
	// неправильно - скопирует меньшее (по len) из 2-х слайсов
	copied := copy(emptyBuf, buf) // copied = 0
	fmt.Println(copied, emptyBuf)

	// правильно
	newBuf = make([]int, len(buf), len(buf))
	copy(newBuf, buf)
	fmt.Println(newBuf)

	// можно копировать в часть существующего слайса
	ints := []int{1, 2, 3, 4}
	copy(ints[1:3], []int{5, 6}) // ints = [1, 5, 6, 4]
	fmt.Println(ints)
}

*/
