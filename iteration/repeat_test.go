package main

import (
	"fmt"
	"testing"
)

func TestRepeat(t *testing.T) {
	t.Run("Test simple", func(t *testing.T) {
		repeated := Repeat("a", 5)
		expected := "aaaaa"
		if repeated != expected {
			t.Errorf("expected %q but got %q", expected, repeated)
		}
	})
	t.Run("Test 10 reps", func(t *testing.T) {
		repeated := Repeat("z", 10)
		expected := "zzzzzzzzzz"
		if repeated != expected {
			t.Errorf("expected %q but got %q", expected, repeated)
		}
	})
}

func BenchmarkRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Repeat("a", 5)
	}
}

func ExampleRepeat() {
	result := Repeat("%", 3)
	fmt.Println(result)
	//output: %%%
}
