package main

import (
	"fmt"
)

// structs and methods
// method is a function associated with type; struct is a type

type Person struct {
	Id      int
	Name    string
	Address string
}

type Account struct {
	Id     int
	Name   string
	No     uint64
	Person // id, name, Address
}

func (p Person) CantSetName(newName string) {
	// N.B. read only method, pass-by-value, you don't need this behaviour in setter
	p.Name = newName
}

func (p *Person) SetName(newName string) {
	// access to original person, pass-by-reference
	p.Name = newName
}

// same method name for different structures

func (a *Account) SetNameDup(newName string) {
	a.Name = newName
}
func (p *Person) SetNameDup(newName string) {
	p.Name = newName
}

type MySlice []int

func (s *MySlice) Size() int {
	return len(*s)
}
func (s *MySlice) Append(v int) *MySlice {
	*s = append(*s, v)
	return s
}

func main() {

	p := Person{Id: 1, Name: "Foo"} // Address = ""

	p.CantSetName("Bar")
	fmt.Printf("person: %#v\n", p) // person: main.Person{Id:1, Name:"Foo"}

	p.SetName("Bar")
	fmt.Printf("person: %#v\n", p) // person: main.Person{Id:1, Name:"Bar"}
	// N.B. compiler automagically treated `p` as reference, converting call to smth like:
	(&p).SetName("Baz")
	fmt.Printf("person: %#v\n", p) // person: main.Person{Id:1, Name:"Baz"}

	// alternatively, w/o magic
	pRef := new(Person) // return a reference
	pRef.SetName("Quix")
	fmt.Printf("person: %#v\n", pRef) // person: &main.Person{Id:0, Name:"Quix"} // N.B. `&` symbol

	// nested structures (composed), outer structure have all methods defined for inner structure
	acc := Account{}
	acc.SetName("Qux")
	fmt.Printf("account: %#v\n", acc)
	// account: main.Account{Id:0, Name:"", No:0x0, Person:main.Person{Id:0, Name:"Qux", Address:""}}

	// if method name is the same for outer and inner structs, outer struct have priority
	acc.SetNameDup("Quuz")
	fmt.Printf("account: %#v\n", acc)
	// account: main.Account{Id:0, Name:"Quuz", No:0x0, Person:main.Person{Id:0, Name:"Qux", Address:""}}

	acc.Person.SetNameDup("Corge")
	fmt.Printf("account: %#v\n", acc)
	// account: main.Account{Id:0, Name:"Quuz", No:0x0, Person:main.Person{Id:0, Name:"Corge", Address:""}}

	// another type methods
	s := MySlice([]int{1, 2})
	_s := s.Append(3)
	fmt.Printf("slice %#v, size %#v, another ref %#v\n", s, s.Size(), _s)
	// slice main.MySlice{1, 2, 3}, size 3, another ref &main.MySlice{1, 2, 3}

}
