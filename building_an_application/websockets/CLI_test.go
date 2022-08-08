package poker_test

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

	poker "github.com/hemanta212/go-with-tdd-project"
)

var (
	dummyStdIn  = &bytes.Buffer{}
	dummyStdOut = &bytes.Buffer{}
)

func TestCLI(t *testing.T) {

	t.Run("Start with 3 players and 'Chris' as winner", func(t *testing.T) {
		in := userSends("3", "Chris wins the game")
		stdout := &bytes.Buffer{}
		game := &poker.SpyGame{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertMessageToUser(t, stdout, poker.PlayerPrompt)
		assertGamePlayerNo(t, game, 3)
		assertGameWinner(t, game, "Chris")
	})

	t.Run("Start with 8 players and 'Cleo' as winner", func(t *testing.T) {
		in := userSends("8", "Cleo wins")
		game := &poker.SpyGame{}

		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		assertGamePlayerNo(t, game, 8)
		assertGameWinner(t, game, "Cleo")
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := userSends("Pies")
		game := &poker.SpyGame{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertGameNotStarted(t, game)
		assertMessageToUser(t, stdout, poker.PlayerPrompt, poker.BadPlayerInputErrMsg)

	})

	t.Run("it prints an error when a non player is sent as winner", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := userSends("7", "Llyod")
		game := &poker.SpyGame{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertGameNotFinished(t, game)
		assertMessageToUser(t, stdout, poker.PlayerPrompt, poker.BadWinnerInputMsg)
	})

}

func assertGamePlayerNo(t *testing.T, game *poker.SpyGame, want int) {
	t.Helper()

	passed := retryUntil(500*time.Millisecond, func() bool {
		game.Mu.Lock()
		defer game.Mu.Unlock()
		return game.PlayersNo == want
	})
	if !passed {
		t.Errorf("Expected numbers of players to be %d but got %d", want, game.PlayersNo)
	}
}

func assertGameWinner(t *testing.T, game *poker.SpyGame, want string) {
	t.Helper()
	passed := retryUntil(500*time.Millisecond, func() bool {
		game.Mu.Lock()
		defer game.Mu.Unlock()
		return game.Winner == want
	})
	if !passed {
		t.Errorf("expected winner called with %q but got %q", want, game.Winner)
	}
}

func assertGameNotStarted(t *testing.T, game *poker.SpyGame) {
	if game.StartCalled {
		t.Errorf("game should not have started")
	}
}
func assertGameNotFinished(t *testing.T, game *poker.SpyGame) {
	if game.FinishCalled {
		t.Error("Game should not have ended, but it did")
	}
}

func userSends(messages ...string) io.Reader {
	return strings.NewReader(strings.Join(messages, "\n"))
}

func assertMessageToUser(t *testing.T, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if got != want {
		t.Errorf("got %q sent to stdOut, expected %q", got, want)
	}
}

func retryUntil(duration time.Duration, f func() bool) bool {
	deadline := time.Now().Add(duration)
	for time.Now().Before(deadline) {
		if f() {
			return true
		}
	}
	return false
}
