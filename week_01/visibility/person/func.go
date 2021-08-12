package person

import (
	"fmt"
)

func NewPerson(id int, name, secret string) *Person {
	// public function
	return &Person{
		ID:     id,
		Name:   name,
		secret: secret,
	}
}

func GetSecret(p *Person) string {
	// public function
	return p.secret
}

func printSecret(p *Person) {
	// private function
	fmt.Println(p.secret)
}
