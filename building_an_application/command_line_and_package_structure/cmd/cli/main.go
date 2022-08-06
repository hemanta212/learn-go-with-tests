package main

import (
	"fmt"
	"log"
	"os"

	poker "github.com/hemanta212/go-with-tdd-project"
)

const dbFileName = "poker.db.json"

func main() {
	store, closeDb, err := poker.FileSystemPlayerStoreFromFilePath(dbFileName)

	if err != nil {
		log.Fatalf("Failed to init player store, %v", err)
	}
	defer closeDb()

	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	cli := poker.NewCLI(store, os.Stdin)
	cli.PlayPoker()
}
