package main

import (
	"log"
	"net/http"
	"sync"
)

func main() {
	server := NewPlayerServer(NewInMemoryPlayerStore())
	log.Fatal(http.ListenAndServe(":5000", server))
}

type inMemoryPlayerStore struct {
	mu     sync.Mutex
	scores map[string]int
}

func NewInMemoryPlayerStore() *inMemoryPlayerStore {
	return &inMemoryPlayerStore{scores: map[string]int{}}
}

func (i *inMemoryPlayerStore) GetPlayerScore(name string) (int, bool) {
	i.mu.Lock()
	defer i.mu.Unlock()
	value, found := i.scores[name]
	return value, found
}
func (i *inMemoryPlayerStore) IncrementPlayerScore(name string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.scores[name]++
}
func (i *inMemoryPlayerStore) GetPlayersList() []Player {
	players := []Player{}
	for player, score := range i.scores {
		players = append(players, Player{player, score})
	}
	return players
}
