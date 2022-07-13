package main

import (
	"errors"
	"fmt"
)

type Wallet int

var ErrInsufficientFunds = errors.New("cannot withdraw, insufficient funds")

func (w *Wallet) Deposit(amount Wallet) {
	*w += amount
}
func (w *Wallet) Withdraw(amount Wallet) error {
	if *w < amount {
		return ErrInsufficientFunds
	}
	*w -= amount
	return nil
}

func (w *Wallet) Balance() Wallet {
	return *w
}

func (w Wallet) String() string {
	return fmt.Sprintf("%d BTC", w)
}
