package poker

import (
	"os"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("League from reader", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
                      {"Name":"Cleo", "Score":10},
                      {"Name":"Chris", "Score":33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)

		got := store.GetLeague()

		want := []Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		AssertLeague(t, got, want)
		// Read again to see if reader is exhausted on one read
		got = store.GetLeague()
		AssertLeague(t, got, want)
	})

	t.Run("Get player score ", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
                      {"Name":"Cleo", "Score":10},
                      {"Name":"Chris", "Score":33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)
		AssertPlayerWin(t, store, "Cleo", 10)
	})

	t.Run("Store wins for existing players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
                      {"Name":"Cleo", "Score":10},
                      {"Name":"Chris", "Score":33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)

		store.IncrementPlayerScore("Chris")
		AssertPlayerWin(t, store, "Chris", 34)
	})

	t.Run("Record wins for new players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
                      {"Name":"Cleo", "Score":10},
                      {"Name":"Chris", "Score":33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)

		store.IncrementPlayerScore("Rick")

		AssertPlayerWin(t, store, "Rick", 1)
	})

	t.Run("League sorted", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
                      {"Name":"Cleo", "Score":10},
                      {"Name":"Chris", "Score":33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)

		got := store.GetLeague()

		want := []Player{
			{Name: "Chris", Score: 33},
			{Name: "Cleo", Score: 10},
		}

		AssertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		AssertLeague(t, got, want)
	})

	t.Run("Works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, ``)
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)
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
