package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type League []Player

func NewLeague(rdr io.Reader) ([]Player, error) {
	players := []Player{}
	err := json.NewDecoder(rdr).Decode(&players)

	if err != nil {
		return nil, fmt.Errorf("Unable to parse league, '%v'", err)
	}

	return players, nil
}

func (l League) Find(playerName string) *Player {
	for i, player := range l {
		if player.Name == playerName {
			return &l[i]
		}
	}
	return nil
}
