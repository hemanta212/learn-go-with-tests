package poker

import (
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
)

type SpyPlayerStore struct {
	scores map[string]int
}

func NewSpyPlayerStore() *SpyPlayerStore {
	return &SpyPlayerStore{map[string]int{}}
}

func (s SpyPlayerStore) GetPlayerScore(player string) (int, bool) {
	score, found := s.scores[player]
	return score, found
}

func (s SpyPlayerStore) IncrementPlayerScore(player string) {
	s.scores[player] = s.scores[player] + 1
}
func (s SpyPlayerStore) GetLeague() League {
	players := make(League, 0, len(s.scores))
	for name, score := range s.scores {
		players = append(players, Player{name, score})
	}
	sort.Slice(players, func(i, j int) bool {
		return players[i].Score > players[j].Score
	})
	return players
}

func AssertPlayerWin(t testing.TB, store PlayerStore, winner string, want int) {
	t.Helper()

	got, found := store.GetPlayerScore(winner)

	if !found {
		t.Fatalf("Expected to find player %s, but didn't find!", winner)
	}

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func AssertPlayerScore(t testing.TB, got, want string) {
	t.Helper()
	if want != got {
		t.Errorf("got %q, want %q", got, want)
	}
}

func AssertLeague(t testing.TB, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func AssertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got status %d, want %d", got, want)
	}
}

func AssertContentType(t testing.TB, res *httptest.ResponseRecorder, want string) {
	t.Helper()
	if res.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of application/json, got %v",
			res.Result().Header)
	}
}

func AssertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}
