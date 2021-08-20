package main

import "fmt"

type Payer interface {
	Pay(int) error
}

type Wallet struct {
	Cash int
}

func (w *Wallet) Pay(amount int) error {
	// implement Payer interface

	// TODO: should be atomic
	if w.Cash < amount {
		return fmt.Errorf("you're too poor")
	}
	w.Cash -= amount
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
	w := &Wallet{Cash: 100} // n.b. reference to wallet, so it can mutate Cash
	Buy(w)
}
