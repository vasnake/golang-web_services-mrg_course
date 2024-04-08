package main

import "fmt"

func main() {
	// hash table, associative array, keys unordered

	// creation, make function
	var user map[string]string = map[string]string{
		"name":     "Bart",
		"lastName": "Simpson",
	}
	fmt.Printf("map %#v, len: %#v\n", user, len(user))
	// map map[string]string{"lastName":"Simpson", "name":"Bart"}, len: 2

	profile := make(map[string]string, 10) // cap = 10
	fmt.Printf("map %#v, len: %#v \n", profile, len(profile))
	// map map[string]string{}, len: 0

	// element absence = element type default value
	// solution: val, exists = map[key]
	// `_` as empty var

	name := user["name"]        // Bart
	mName := user["middleName"] // wrong, default value ""
	println(name, mName)

	mName, mNameExists := user["middleName"]
	_, mNameExists = user["middleName"] // only existence flag
	println(name, mNameExists)

	// function delete(map, key)
	delete(user, "lastName")
	fmt.Printf("map %#v, len: %#v\n", user, len(user))
	// map map[string]string{"name":"Bart"}, len: 1

}

/*
	// инициализация при создании
	var user map[string]string = map[string]string{
		"name":     "Vasily",
		"lastName": "Romanov",
	}

	// сразу с нужной ёмкостью
	profile := make(map[string]string, 10)

	// количество элементов
	mapLength := len(user)

	fmt.Printf("%d %+v\n", mapLength, profile)

	// если ключа нет - вернёт значение по умолчанию для типа
	mName := user["middleName"]
	fmt.Println("mName:", mName)

	// проверка на существование ключа
	mName, mNameExist := user["middleName"]
	fmt.Println("mName:", mName, "mNameExist:", mNameExist)

	// пустая переменная - только проверяем что ключ есть
	_, mNameExist2 := user["middleName"]
	fmt.Println("mNameExist2", mNameExist2)

	// удаление ключа
	delete(user, "lastName")
	fmt.Printf("%#v\n", user)
}

*/
