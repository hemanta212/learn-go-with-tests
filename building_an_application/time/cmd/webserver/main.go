package main

import (
	"log"
	"net/http"

	poker "github.com/hemanta212/go-with-tdd-project"
)

const dbFileName = "game.db.json"

func main() {
	store, closeDb, err := poker.FileSystemPlayerStoreFromFilePath(dbFileName)

	if err != nil {
		log.Fatalf("Failed to init player store, %v", err)
	}
	defer closeDb()

	server := poker.NewPlayerServer(store)
	log.Fatal(http.ListenAndServe(":5000", server))
}
