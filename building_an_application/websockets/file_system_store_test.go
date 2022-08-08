package poker_test

import (
	"os"
	"testing"

	poker "github.com/hemanta212/go-with-tdd-project"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("League from reader", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
                      {"Name":"Cleo", "Score":10},
                      {"Name":"Chris", "Score":33}]`)
		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)
		poker.AssertNoError(t, err)

		got := store.GetLeague()

		want := []poker.Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		poker.AssertLeague(t, got, want)
		// Read again to see if reader is exhausted on one read
		got = store.GetLeague()
		poker.AssertLeague(t, got, want)
	})

	t.Run("Get player score ", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
                      {"Name":"Cleo", "Score":10},
                      {"Name":"Chris", "Score":33}]`)
		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)
		poker.AssertNoError(t, err)
		poker.AssertPlayerWin(t, store, "Cleo", 10)
	})

	t.Run("Store wins for existing players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
                      {"Name":"Cleo", "Score":10},
                      {"Name":"Chris", "Score":33}]`)
		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)
		poker.AssertNoError(t, err)

		store.IncrementPlayerScore("Chris")
		poker.AssertPlayerWin(t, store, "Chris", 34)
	})

	t.Run("Record wins for new players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
                      {"Name":"Cleo", "Score":10},
                      {"Name":"Chris", "Score":33}]`)
		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)
		poker.AssertNoError(t, err)

		store.IncrementPlayerScore("Rick")

		poker.AssertPlayerWin(t, store, "Rick", 1)
	})

	t.Run("League sorted", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
                      {"Name":"Cleo", "Score":10},
                      {"Name":"Chris", "Score":33}]`)
		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)
		poker.AssertNoError(t, err)

		got := store.GetLeague()

		want := []poker.Player{
			{Name: "Chris", Score: 33},
			{Name: "Cleo", Score: 10},
		}

		poker.AssertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		poker.AssertLeague(t, got, want)
	})

	t.Run("Works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, ``)
		defer cleanDatabase()

		_, err := poker.NewFileSystemPlayerStore(database)
		poker.AssertNoError(t, err)
	})

}

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "db")
	if err != nil {
		t.Fatalf("could not create a temp file %v", err)
	}
	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}
