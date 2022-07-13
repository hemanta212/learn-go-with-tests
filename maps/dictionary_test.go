package main

import (
	"testing"
)

func TestSearch(t *testing.T) {
	dictionary := Dictionary{"test": "a test"}

	t.Run("Known key", func(t *testing.T) {
		got, _ := dictionary.Search("test")
		want := "a test"
		assertStrings(t, got, want)
	})
	t.Run("Unknown key", func(t *testing.T) {
		_, err := dictionary.Search("unknown")
		assertError(t, err, ErrorNotFound)
	})

}

func TestAdd(t *testing.T) {
	t.Run("Add new key", func(t *testing.T) {
		dictionary := Dictionary{}
		key, value := "test", "a test"
		err := dictionary.Add(key, value)
		assertError(t, err, nil)
		assertDefinition(t, dictionary, key, value)
	})
	t.Run("Add existing key", func(t *testing.T) {
		key, value := "test", "a test"
		dictionary := Dictionary{key: value}
		err := dictionary.Add(key, value)
		if err == nil {
			t.Fatal("Key got overwritten")
		}
		assertError(t, err, ErrorKeyExists)
		assertDefinition(t, dictionary, key, value)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Update existing key", func(t *testing.T) {
		key, value := "test", "a test"
		dictionary := Dictionary{key: value}
		newValue := "another test"
		err := dictionary.Update(key, newValue)
		assertError(t, err, nil)
		assertDefinition(t, dictionary, key, newValue)
	})
	t.Run("Update unknown/new key", func(t *testing.T) {
		dictionary := Dictionary{}
		key, value := "test", "a test"
		err := dictionary.Update(key, value)
		assertError(t, err, ErrorNotFound)
		_, err = dictionary.Search(key)
		assertError(t, err, ErrorNotFound)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Delete existing key", func(t *testing.T) {
		key, value := "test", "a test"
		dictionary := Dictionary{key: value}
		dictionary.Delete(key)
		_, err := dictionary.Search(key)
		assertError(t, err, ErrorNotFound)
	})
}

func assertDefinition(t *testing.T, d Dictionary, key, want string) {
	t.Helper()
	got, err := d.Search(key)
	if err != nil {
		t.Fatal("Added word not found:", err)
	}
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertStrings(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
func assertError(t *testing.T, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("got error %q want %q", got, want)
	}
}
