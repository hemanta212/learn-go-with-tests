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

	game := poker.NewTexasHoldem(store, poker.BlindAlerterFunc(poker.Alerter))

	server, err := poker.NewPlayerServer(store, game)
	if err != nil {
		log.Fatalf("Failed to init player server, %v", err)
	}
	log.Fatal(http.ListenAndServe(":5000", server))
}
