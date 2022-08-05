package main

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
		assertNoError(t, err)

		got := store.GetLeague()

		want := []Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		assertLeague(t, got, want)
		// Read again to see if reader is exhausted on one read
		got = store.GetLeague()
		assertLeague(t, got, want)
	})

	t.Run("Get player score ", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
                      {"Name":"Cleo", "Score":10},
                      {"Name":"Chris", "Score":33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		got, found := store.GetPlayerScore("Cleo")
		want := 10
		if !found {
			t.Fatal("Expected to find player Cleo but didn't find")
		}
		assertScoreEquals(t, got, want)
	})

	t.Run("Store wins for existing players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
                      {"Name":"Cleo", "Score":10},
                      {"Name":"Chris", "Score":33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		store.IncrementPlayerScore("Chris")

		got, found := store.GetPlayerScore("Chris")
		want := 34
		if !found {
			t.Fatal("Expected to find player Chris but didn't find")
		}
		assertScoreEquals(t, got, want)
	})

	t.Run("Record wins for new players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
                      {"Name":"Cleo", "Score":10},
                      {"Name":"Chris", "Score":33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		store.IncrementPlayerScore("Rick")

		got, found := store.GetPlayerScore("Rick")
		want := 1
		if !found {
			t.Fatal("Expected to find player Rick but didn't find")
		}
		assertScoreEquals(t, got, want)
	})

	t.Run("League sorted", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
                      {"Name":"Cleo", "Score":10},
                      {"Name":"Chris", "Score":33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		got := store.GetLeague()

		want := []Player{
			{Name: "Chris", Score: 33},
			{Name: "Cleo", Score: 10},
		}

		assertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		assertLeague(t, got, want)
	})

	t.Run("Works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, ``)
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)
	})

}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
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
