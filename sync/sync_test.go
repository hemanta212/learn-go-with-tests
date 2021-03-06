package main

import (
	"sync"
	"testing"
)

func TestSync(t *testing.T) {
	t.Run("Incrementing counter 3 times leaves it at 3", func(t *testing.T) {

		counter := NewCounter()
		counter.Inc()
		counter.Inc()
		counter.Inc()

		assertCounter(t, counter, 3)
	})

	t.Run("1k incs without goroutine", func(t *testing.T) {

		wantedCount := 1000
		counter := NewCounter()

		for i := 0; i < wantedCount; i++ {
			counter.Inc()
		}
		assertCounter(t, counter, wantedCount)
	})

	t.Run("It runs safe concurrently", func(t *testing.T) {
		wantedCount := 1000
		counter := NewCounter()

		var wg sync.WaitGroup
		wg.Add(wantedCount)

		for i := 0; i < wantedCount; i++ {
			go func() {
				counter.Inc()
				wg.Done()
			}()
		}
		wg.Wait()

		assertCounter(t, counter, wantedCount)
	})

}

func assertCounter(t testing.TB, got *Counter, want int) {
	t.Helper()
	if got.Value() != want {
		t.Errorf("got %d want %d", got.Value(), want)
	}
}
