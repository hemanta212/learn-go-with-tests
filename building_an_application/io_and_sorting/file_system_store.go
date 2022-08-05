package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type FileSystemPlayerStore struct {
	database *json.Encoder
	league   League
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
	err := initializePlayerDBFile(file)

	if err != nil {
		return nil, fmt.Errorf("Problem initializing db file, %v", err)
	}

	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v",
			file.Name(), err)
	}

	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}, nil
}

func (f *FileSystemPlayerStore) GetLeague() League {
	sort.Slice(f.league, func(i, j int) bool {
		return f.league[i].Score > f.league[j].Score
	})
	return f.league
}

func (f *FileSystemPlayerStore) GetPlayerScore(playerName string) (int, bool) {
	player := f.league.Find(playerName)
	if player == nil {
		return 0, false
	}
	return player.Score, player != nil
}

func (f *FileSystemPlayerStore) IncrementPlayerScore(playerName string) {
	player := f.league.Find(playerName)

	if player != nil {
		player.Score++
	} else {
		f.league = append(f.league, Player{playerName, 1})
	}

	f.database.Encode(f.league)
}

func initializePlayerDBFile(file *os.File) error {
	file.Seek(0, 0)

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("problem getting stat info from file %s, %v",
			file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, 0)
	}
	return nil
}
