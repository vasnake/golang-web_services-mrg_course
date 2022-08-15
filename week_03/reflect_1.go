package main

import (
	"fmt"
	"reflect"
)

type User struct {
	ID       int
	RealName string `unpack:"-"`
	Login    string
	Flags    int
}

func PrintReflect(u interface{}) error {

	val := reflect.ValueOf(u).Elem() // it's a key to all
	fmt.Printf("%T have %d fields:\n", u, val.NumField())

	for i := 0; i < val.NumField(); i++ {

		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		fmt.Printf(
			"\tname=%v, type=%v, value=%v, tag=`%v`\n",
			typeField.Name,
			typeField.Type.Kind(),
			valueField,
			typeField.Tag,
		)
	}

	return nil
}

func main() {

	u := &User{
		ID:       42,
		RealName: "rvasily",
		Flags:    32,
	}

	err := PrintReflect(u) // examine the output, please

	if err != nil {
		panic(err)
	}
}
