package main

import "testing"

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
