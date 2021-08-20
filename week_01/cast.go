package main

import "fmt"

func Buy(p Payer) {
	// if you want to check type of the object or use custom logic, type switch
	switch p.(type) {
	case *Wallet:
		fmt.Println("cash")
	case *Card:
		card, err := p.(*Card)
		if err { // not much sense in doing it here
			fmt.Println("not a card, really")
		}
		fmt.Println("prepare to authorize payment,", card.CardHolder)
	default:
		fmt.Println("unknown payable object")
	}

	err := p.Pay(10)
	if err != nil {
		panic(err)
	}
	fmt.Printf("You paid with %T\n", p)
}
