package methods

// structs and methods
// method is a function associated with type; struct is a type

type Person struct {
	Id      int
	Name    string
	Address string
}

func (p Person) CantSetName(newName string) {
	// N.B. read only method, pass-by-value, you don't need this behaviour in setter
	p.Name = newName
}

func (p *Person) SetName(newName string) {
	// access to original person, pass-by-reference
	p.Name = newName
}

/////////////////////////////////////////////////////////////////////////////////////////////////

type Account struct {
	Id     int
	Name   string
	No     uint64
	Person // id, name, address: composition
}

// same method name for different structures
func (a *Account) SetNameDup(newName string) {
	a.Name = newName
}
func (p *Person) SetNameDup(newName string) {
	p.Name = newName
}

/////////////////////////////////////////////////////////////////////////////////////////////////

type MySlice []int

func (s *MySlice) Size() int {
	return len(*s)
}

func (s *MySlice) Append(v int) *MySlice {
	*s = append(*s, v)
	return s
}
