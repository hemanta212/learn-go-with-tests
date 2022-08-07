package poker

import "time"

type Game interface {
	Start(numberOfPlayers int)
	Finish(winner string)
}

type TexasHoldemGame struct {
	Store   PlayerStore
	Alerter BlindAlerter
}

func NewTexasHoldem(store PlayerStore, alerter BlindAlerter) *TexasHoldemGame {
	return &TexasHoldemGame{
		Store:   store,
		Alerter: alerter,
	}
}

func (game *TexasHoldemGame) Start(numberOfPlayers int) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		game.Alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + blindIncrement
	}
}

func (game *TexasHoldemGame) Finish(winner string) {
	game.Store.IncrementPlayerScore(winner)
}
