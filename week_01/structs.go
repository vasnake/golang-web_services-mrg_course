package main

import "fmt"

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
	Owner   Person              // struct as field
}

// PersonalAccount shows structs composition: add Person fields to Account
type PersonalAccount struct {
	Id      int
	Name    string
	Cleaner func(string) string
	Owner   Person
	Person  // composition: Account id, name and address
}

func main() {

	// struct initialization, named fields init
	var acc Account = Account{
		Id:   42,
		Name: "foo",
	} // remember: you'll have default values for unset fields

	fmt.Printf("Account w/o owner: %#v\n", acc) // repr
	// main.Account{Id:42, Name:"foo", Cleaner:(func(string) string)(nil), Owner:main.Person{Id:0, Name:"", Address:""}}

	// short init, fixed sequence only
	acc.Owner = Person{33, "Foo Bar", "Under The Bridge"}
	fmt.Printf("Account with owner: %#v\n", acc)
	// main.Account{Id:42, Name:"foo", Cleaner:(func(string) string)(nil), Owner:main.Person{Id:33, Name:"Foo Bar", Address:"Under The Bridge"}}

	// composed struct
	pa := PersonalAccount{Id: 77, Name: "baz", Person: Person{Address: "Far Far Away", Name: "zab"}}
	fmt.Printf("Account with owner: %#v\n", pa)
	// Account with owner: main.PersonalAccount{
	// Id:77, Name:"baz", Cleaner:(func(string) string)(nil),
	// Owner:main.Person{Id:0, Name:"", Address:""},
	// Person:main.Person{Id:0, Name:"", Address:"Far Far Away"}}

	// PA.Person fields, n.b. w/o `Person` prefix, using Account namespace
	fmt.Printf("address: %#v\n", pa.Address) // address: "Far Far Away"

	// what about `Name`? outer field selected
	fmt.Printf("PA.Name: %#v\n", pa.Name)               // PA.Name: "baz"
	fmt.Printf("PA.Person.Name: %#v\n", pa.Person.Name) // PA.Person.Name: "zab"

}
