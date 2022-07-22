package main

import "testing"

func TestStack(t *testing.T) {
	t.Run("Integer stack", func(t *testing.T) {
		stack := new(Stack[int])

		// check if new stack is empty
		AssertTrue(t, stack.IsEmpty())

		// add a thing, then check its not empty
		stack.Push(123)
		AssertFalse(t, stack.IsEmpty())

		// Add another thing and pop back again
		stack.Push(22)
		value, _ := stack.Pop()
		AssertEqual(t, value, 22)
		value, _ = stack.Pop()
		AssertEqual(t, value, 123)
		AssertTrue(t, stack.IsEmpty())

		// can get num as we put in num, not untyped interface{}
		stack.Push(1)
		stack.Push(2)
		firstNum, _ := stack.Pop()
		secondNum, _ := stack.Pop()
		AssertEqual(t, firstNum+secondNum, 3)
	})
}

func TestAssertFunctions(t *testing.T) {
	t.Run("asserting on intergers", func(t *testing.T) {
		AssertEqual(t, 1, 1)
		AssertNotEqual(t, 1, 2)
	})
	t.Run("asserting on strings", func(t *testing.T) {
		AssertEqual(t, "hello", "hello")
		AssertNotEqual(t, "hello", "hi")
	})
	// should fail due to error
	// AssertEqual(t, 1, "1")
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
