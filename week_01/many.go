package main

import (
	"fmt"
)

// payable interface

type Payer interface {
	Pay(int) error
}

// payable wallet

type Wallet struct {
	Cash int
}

func (w *Wallet) Pay(amount int) error {
	if w.Cash < amount {
		return fmt.Errorf("not enough money")
	}
	w.Cash -= amount // should be atomic
	return nil
}

// payable card

type Card struct {
	Balance    int
	ValidUntil string
	CardHolder string
	CVV        string
	Number     string
}

func (c *Card) Pay(amount int) error {
	if c.Balance < amount {
		return fmt.Errorf("not enough money")
	}
	c.Balance -= amount // see Wallet.Pay
	return nil
}

// uses interface, don't give a fuck about data type

func Buy(p Payer) {
	// any object with Payer implementation will do, e.g. Wallet
	err := p.Pay(10)
	if err != nil {
		panic(err)
	}
	fmt.Printf("You paid with %T\n", p)
}

func main() {
	// buy with any kind of payable object
	w := &Wallet{Cash: 100}
	Buy(w)

	var p Payer
	p = &Card{Balance: 100}
	Buy(p)

	p = w
	Buy(p)
}
