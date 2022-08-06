package poker_test

import (
	"strings"
	"testing"

	poker "github.com/hemanta212/go-with-tdd-project"
)

func TestCLI(t *testing.T) {

	t.Run("Record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris\n")
		playerStore := poker.NewSpyPlayerStore()

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Chris", 1)
	})

	t.Run("Record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo Wins\n")
		playerStore := poker.NewSpyPlayerStore()

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Cleo", 1)
	})

}
