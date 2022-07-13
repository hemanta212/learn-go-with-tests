package main

import (
	"testing"
)

func TestWallet(t *testing.T) {
	t.Run("Init wallet", func(t *testing.T) {
		var wallet Wallet
		assertBalance(t, wallet, Wallet(0))
	})
	t.Run("Deposit 10&10", func(t *testing.T) {
		var wallet Wallet
		wallet.Deposit(10)
		wallet.Deposit(10)
		assertBalance(t, wallet, Wallet(20))
	})
	t.Run("Withdraw 10", func(t *testing.T) {
		var wallet Wallet
		wallet.Deposit(20)
		err := wallet.Withdraw(10)
		assertNoError(t, err)
		assertBalance(t, wallet, Wallet(10))
	})
	t.Run("Withdraw Insufficient Funds", func(t *testing.T) {
		startingBalance := 20
		wallet := Wallet(startingBalance)
		err := wallet.Withdraw(30)
		assertError(t, err, ErrInsufficientFunds)
		assertBalance(t, wallet, Wallet(startingBalance))
	})
}

func assertBalance(t testing.TB, wallet Wallet, want Wallet) {
	t.Helper()
	got := wallet.Balance()

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
func assertError(t testing.TB, got, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("Wanted an error but didn't get one")
	}
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatal("got an error, but didn't want one")
	}
}
