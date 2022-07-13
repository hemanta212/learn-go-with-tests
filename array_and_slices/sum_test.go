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
