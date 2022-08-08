package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	PlayerPrompt         = "Please enter the number of players: "
	BadPlayerInputErrMsg = "Bad value received for number of players, please try again with a number"
	BadWinnerInputMsg    = "Invalid winner input, expect format of 'PlayerName wins'"
)

type CLI struct {
	in   *bufio.Scanner
	out  io.Writer
	game Game
}

func NewCLI(input io.Reader, output io.Writer, game Game) *CLI {
	return &CLI{
		in:   bufio.NewScanner(input),
		out:  output,
		game: game,
	}
}

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)

	playersNo, err := strconv.Atoi(cli.readLine())

	if err != nil {
		fmt.Fprint(cli.out, BadPlayerInputErrMsg)
		return
	}

	cli.game.Start(playersNo, cli.out)

	userInput := cli.readLine()

	if !strings.Contains(userInput, " wins") {
		fmt.Fprint(cli.out, BadWinnerInputMsg)
		return
	}

	winner := extractWinner(userInput)

	cli.game.Finish(winner)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

func extractWinner(text string) string {
	return strings.Split(text, " ")[0]
}
