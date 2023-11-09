package interfaces

import (
	"fmt"
	"strconv"
)

type Payer interface {
	Pay(int) error
}

////////////////////////////////////////////////////////////////////////

type Wallet struct {
	Cash int
}

func (w *Wallet) Pay(amount int) error {
	// implement Payer interface

	// N.B. don't do that in prod
	if w.Cash < amount {
		return fmt.Errorf("you're too poor")
	}
	w.Cash -= amount
	return nil
}

////////////////////////////////////////////////////////////////////////

type Card struct {
	Balance    int
	ValidUntil string
	CardHolder string
	CVV        string
	Number     string
}

func (c *Card) Pay(amount int) error {
	// implement Payer interface

	// N.B. don't do that in prod
	if c.Balance < amount {
		return fmt.Errorf("not enough money")
	}
	c.Balance -= amount
	return nil
}

////////////////////////////////////////////////////////////////////////

// wallet with String method
type WalletNice struct {
	Cash int
}

func (w *WalletNice) String() string {
	// Stringer interface
	return "Wallet with " + strconv.Itoa(w.Cash) + " currency points"
}

////////////////////////////////////////////////////////////////////////

type Phone struct {
	// should implement Payer and Ringer
	Money int
}

func (p *Phone) Pay(amount int) error {
	fmt.Println("pay", amount, "coins")
	return nil
}

func (p *Phone) Ring(what string) error {
	fmt.Println("ring", what)
	return nil
}
