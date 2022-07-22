package main

import (
	"reflect"
	"testing"
)

func TestSum(t *testing.T) {
	t.Run("Test array of any size 4-10:49", func(t *testing.T) {
		numbers := []int{4, 5, 6, 7, 8, 9, 10}
		got := Sum(numbers)
		want := 49
		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

}

func TestSumAll(t *testing.T) {
	t.Run("Test 2 slices	", func(t *testing.T) {
		got := SumAll([]int{1, 2}, []int{0, 9})
		want := []int{3, 9}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestSumAllTails(t *testing.T) {
	checkSums := func(t testing.TB, got, want []int) {
		t.Helper()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	}

	t.Run("Test 2 slices", func(t *testing.T) {
		got := SumAllTails([]int{1, 2}, []int{0, 9})
		want := []int{2, 9}
		checkSums(t, got, want)
	})
	t.Run("Test with lenght 3+", func(t *testing.T) {
		got := SumAllTails([]int{1, 2, 3}, []int{0, 5, 9})
		want := []int{5, 14}
		checkSums(t, got, want)
	})
	t.Run("Test with empty array", func(t *testing.T) {
		got := SumAllTails([]int{}, []int{3, 4, 5})
		want := []int{0, 9}
		checkSums(t, got, want)
	})
}

func TestReduce(t *testing.T) {
	t.Run("Multiply all elements", func(t *testing.T) {
		got := Reduce([]int{2, 3, 5, 6, 3}, func(x, y int) int {
			return x * y
		}, 1)
		AssertEqual(t, got, 540)
	})
	t.Run("Concatenate Strings", func(t *testing.T) {
		concat := func(x, y string) string {
			return x + y
		}
		got := Reduce([]string{"a", "b", "c"}, concat, "")
		AssertEqual(t, got, "abc")
	})
}

func TestBadBank(t *testing.T) {
	transactions := []Transaction{
		{
			From: "Chris",
			To:   "Riya",
			Sum:  100,
		},
		{
			From: "Adil",
			To:   "Chris",
			Sum:  25,
		},
	}
	AssertEqual(t, BalanceFor(transactions, "Riya"), 100)
	AssertEqual(t, BalanceFor(transactions, "Chris"), -75)
	AssertEqual(t, BalanceFor(transactions, "Adil"), -25)
}

func AssertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func AssertNotEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got == want {
		t.Errorf("got %v want %v", got, want)
	}
}
func AssertTrue(t *testing.T, got bool) {
	t.Helper()
	if !got {
		t.Errorf("got %v, want true", got)
	}
}
func AssertFalse(t *testing.T, got bool) {
	t.Helper()
	if got {
		t.Errorf("got %v, want false", got)
	}
}
