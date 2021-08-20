package main

import "fmt"

type Payer interface {
	Pay(int) error
}

type Wallet struct {
	Cash int
}

func (w *Wallet) Pay(amount int) error {
	// n.b. reference to Wallet object, type (*Wallet, w) is not the same as (Wallet, w)
	if w.Cash < amount {
		return fmt.Errorf("not enough money")
	}
	w.Cash -= amount // should be atomic
	return nil
}

// try to buy with anything

func Buy(in interface{}) {
	// empty interface, dynamic type check
	var p Payer
	var ok bool

	if p, ok = in.(Payer); !ok {
		fmt.Printf("%T is not a Payer\n", in)
		return
	}

	err := p.Pay(10)
	if err != nil {
		fmt.Printf("Paying error, %T %v\n", p, err)
		return
	}
	fmt.Printf("You paid with %T\n", p)
}

func main() {
	w := Wallet{Cash: 100}
	Buy(&w) // n.b. reference, it's vital for type cast
	Buy([]int{1, 2, 3})
	Buy(3.14)
}
