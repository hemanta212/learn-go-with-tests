package main

import (
	"log"
	"net/http"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Problem opening %s %v", dbFileName, err)
	}
	store, err := NewFileSystemPlayerStore(db)
	if err != nil {
		log.Fatalf("Problem creating file system player store, %v ", err)
	}
	server := NewPlayerServer(store)
	log.Fatal(http.ListenAndServe(":5000", server))
}
