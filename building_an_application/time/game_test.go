package poker_test

import (
	"fmt"
	"testing"
	"time"

	poker "github.com/hemanta212/go-with-tdd-project"
)

var (
	dummyBlindAlerter = &SpyBlindAlerter{}
	dummyPlayerStore  = poker.NewSpyPlayerStore()
)

func TestGame_Start(t *testing.T) {
	t.Run("it schedules printing of blind values", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}

		game := poker.TexasHoldemGame{dummyPlayerStore, blindAlerter}
		game.Start(5)

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {

				if len(blindAlerter.alerts) <= i {
					t.Fatalf("Alert %d was not scheduled %v",
						i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]
				assertScheduledAlert(t, got, want)
			})

		}
	})

	t.Run("Schedule alerts on game start for 7 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}

		game := poker.NewTexasHoldem(dummyPlayerStore, blindAlerter)
		game.Start(7)

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}
		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}
				got := blindAlerter.alerts[i]
				assertScheduledAlert(t, got, want)
			})
		}
	})

}

func TestGame_Finish(t *testing.T) {
	store := poker.NewSpyPlayerStore()
	game := poker.NewTexasHoldem(store, dummyBlindAlerter)
	winner := "Ruth"

	game.Finish(winner)
	poker.AssertPlayerWin(t, store, winner, 1)
}

type scheduledAlert struct {
	at     time.Duration
	amount int
}

func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduledAlert{duration, amount})
}

func assertScheduledAlert(t testing.TB, got, want scheduledAlert) {
	t.Helper()
	if got.amount != want.amount {
		t.Errorf("got amount %d want %d", got.amount, want.amount)
	}

	if got.at != want.at {
		t.Errorf("got scheduled time of %v, want %v", got.at, want.at)
	}
}
