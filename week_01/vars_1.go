package main

import (
	"fmt"
	"go/types"
	"reflect"
)

func main() {
	/*
		// значение по умолчанию
		var num0 int

		// значение при инициализации
		var num1 int = 1

		// пропуск типа
		var num2 = 20
		fmt.Println(num0, num1, num2)

		// короткое объявление переменной
		num := 30
		// только для новых переменных
		// no new variables on left side of :=
		// num := 31

		num += 1
		fmt.Println("+=", num)

		// ++num нету
		num++
		fmt.Println("++", num)
	*/

	var num0 int // platform dependent int, 32 or 64 or ? with default value
	fmt.Println("1. default int: ", num0)

	var num1 int = 1 // set init value
	fmt.Println("2. with init value: ", num1)

	var num2 = 2 // type inference
	fmt.Println("3. detected type: ", num2, types.BasicInfo(num2), reflect.TypeOf(num2).String())
	fmt.Printf("4. type: %T\n", num2)

	num := 3 // skip `var` and `type`
	fmt.Println("5. new var: ", num, reflect.TypeOf(num).String())

	num += 1 // increment
	num++    // postfix increment
	fmt.Println("6. incremented: ", num)

	var a, b = 11, 12 // define multiple vars
	a, c := 13, 14    // only if `c` is new
	a, b, c = 1, 2, 3 // multi assign
	fmt.Println("7. multi: ", a, b, c)

	// camelCase - принятый стиль
	userIndex := 10
	// under_score - не принято
	user_index := 10
	fmt.Println(userIndex, user_index)

	// объявление нескольких переменных
	var weight, height int = 10, 20

	// присваивание в существующие переменные
	weight, height = 11, 21

	// короткое присваивание
	// хотя-бы одна переменная должна быть новой!
	weight, age := 12, 22

	fmt.Println(weight, height, age)
}
