package poker

import (
	"bufio"
	"io"
	"strings"
)

type CLI struct {
	store PlayerStore
	in    *bufio.Scanner
}

func NewCLI(store PlayerStore, input io.Reader) *CLI {
	return &CLI{store, bufio.NewScanner(input)}
}

func (cli CLI) PlayPoker() {
	userInput := cli.readLine()
	cli.store.IncrementPlayerScore(extractWinner(userInput))
}

func (cli CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

func extractWinner(text string) string {
	return strings.Split(text, " ")[0]
}
