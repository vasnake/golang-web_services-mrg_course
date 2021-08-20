package main

import (
	"fmt"
	"strconv"
)

type Wallet struct {
	Cash int
}

// wallet with String method

type WalletNice struct {
	Cash int
}

func (w *WalletNice) String() string {
	// Stringer interface
	return "Wallet with " + strconv.Itoa(w.Cash) + " currency points"
}

func main() {
	w := &Wallet{Cash: 100}
	// printf accepts any empty interface
	// empty interface don't impose any limitations, it's just tuple (type, value)
	fmt.Printf("repr: %#v\n", w)                               // &main.Wallet{Cash:100}
	fmt.Printf("Stringer String: %s\n", w)                     // &{%!s(int=100)}
	fmt.Printf("Stringer String: %s\n", &WalletNice{Cash: 99}) // Wallet with 99 currency points

}
