package main

import (
	"testing"
)

func TestWallet(t *testing.T) {
	t.Run("Init wallet", func(t *testing.T) {
		var wallet Wallet
		assertBalance(t, wallet, Bitcoin(0))
	})
	t.Run("Deposit 10&10", func(t *testing.T) {
		var wallet Wallet
		wallet.Deposit(Bitcoin(10))
		wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(20))
	})
	t.Run("Withdraw 10", func(t *testing.T) {
		var wallet Wallet
		wallet.Deposit(Bitcoin(20))
		err := wallet.Withdraw(Bitcoin(10))
		assertNoError(t, err)
		assertBalance(t, wallet, Bitcoin(10))
	})
	t.Run("Withdraw Insufficient Funds", func(t *testing.T) {
		startingBalance := Bitcoin(20)
		wallet := Wallet{balance: startingBalance}
		err := wallet.Withdraw(Bitcoin(30))
		assertError(t, err, ErrInsufficientFunds)
		assertBalance(t, wallet, startingBalance)
	})
}

func assertBalance(t testing.TB, wallet Wallet, want Bitcoin) {
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
