package main

import (
	"fmt"
	"testing"
)

func TestAddr(t *testing.T) {
	t.Run("Add 2 plus 2", func(t *testing.T) {
		sum := Add(2, 2)
		expected := 4
		if sum != expected {
			t.Errorf("Expected '%d' but got '%d'", expected, sum)
		}
	})
	t.Run("Add 200 plus -200", func(t *testing.T) {
		sum := Add(200, -200)
		expected := 0
		if sum != expected {
			t.Errorf("Expected '%d' but got '%d'", expected, sum)
		}
	})
}

func ExampleAdd() {
	sum := Add(1, 5)
	fmt.Println(sum)
	//output: 6
}
