package main

import "fmt"

type Payer interface {
	Pay(int) error
}

type Ringer interface {
	Ring(string) error
}

type NFCPhone interface {
	// combined interface
	Payer
	Ringer
}

func PayAndRing(phone NFCPhone) {
	// it want Payer and Ringer

	err := phone.Pay(1)
	if err != nil {
		fmt.Printf("paying error %v\n", err)
		return
	}

	fmt.Printf("payed via %T\n", phone)

	err = phone.Ring("me")
	if err != nil {
		fmt.Printf("ringing error %v\n", err)
		return
	}

}

func main() {
	p := &Phone{Money: 9}
	PayAndRing(p) // using both interfaces
}

type Phone struct {
	// should implement Payer and Ringer
	Money int
}

func (p *Phone) Pay(amount int) error {
	fmt.Println("pay", amount, "coins")
	return nil
}

func (p *Phone) Ring(number string) error {
	fmt.Println("ring to", number)
	return nil
}
