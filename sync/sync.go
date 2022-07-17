package main

import "sync"

type Counter struct {
	mu    sync.Mutex
	count int
}

// Initialize and use reference of counter
// to avoid copying mutex
func NewCounter() *Counter {
	return &Counter{}
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *Counter) Value() int {
	return c.count
}
