package poker_test

import (
	poker "github.com/hemanta212/go-with-tdd-project"
	"io"
	"testing"
)

func TestTape_Write(t *testing.T) {
	file, clean := createTempFile(t, "12345")
	defer clean()

	tape := &poker.Tape{file}
	tape.Write([]byte("abc"))

	file.Seek(0, 0)
	newFileContents, _ := io.ReadAll(file)

	got := string(newFileContents)
	want := "abc"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
