package poker_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

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
		game := &SpyGame{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertMessageToUser(t, stdout, poker.PlayerPrompt)
		assertPlayerNo(t, game.PlayersNo, 3)
		assertWinner(t, game.Winner, "Chris")
	})

	t.Run("Start with 8 players and 'Cleo' as winner", func(t *testing.T) {
		in := userSends("8", "Cleo wins")
		game := &SpyGame{}

		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		assertPlayerNo(t, game.PlayersNo, 8)
		assertWinner(t, game.Winner, "Cleo")
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := userSends("Pies")
		game := &SpyGame{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertGameNotStarted(t, game)
		assertMessageToUser(t, stdout, poker.PlayerPrompt, poker.BadPlayerInputErrMsg)

	})

	t.Run("it prints an error when a non player is sent as winner", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := userSends("7", "Llyod")
		game := &SpyGame{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertGameNotFinished(t, game)
		assertMessageToUser(t, stdout, poker.PlayerPrompt, poker.BadWinnerInputMsg)
	})

}

type SpyGame struct {
	StartCalled  bool
	PlayersNo    int
	Winner       string
	FinishCalled bool
}

func (game *SpyGame) Start(playersNo int) {
	game.StartCalled = true
	game.PlayersNo = playersNo
}
func (game *SpyGame) Finish(winner string) {
	game.Winner = winner
	game.FinishCalled = true
}

func assertPlayerNo(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("game called with player number %d, instead of %d", got, want)
	}
}

func assertWinner(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, expected %q", got, want)
	}
}

func assertGameNotStarted(t *testing.T, game *SpyGame) {
	if game.StartCalled {
		t.Errorf("game should not have started")
	}
}
func assertGameNotFinished(t *testing.T, game *SpyGame) {
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
