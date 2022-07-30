package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetPlayers(t *testing.T) {
	store := SpyPlayerStore{
		map[string]int{
			"Amily": 20,
			"Rick":  0,
		},
	}
	server := NewPlayerServer(store)

	t.Run("check scores for Amily", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := getRequestFor("Amily")

		server.ServeHTTP(response, request)

		assertPlayerScore(t, response.Body.String(), "20")
	})
	t.Run("check scores for Rick", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := getRequestFor("Rick")

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertPlayerScore(t, response.Body.String(), "0")
	})
	t.Run("Return 404 on missing players", func(t *testing.T) {
		request := getRequestFor("nonplayer")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestPostPlayerWins(t *testing.T) {
	store := SpyPlayerStore{
		scores: map[string]int{},
	}
	server := NewPlayerServer(store)

	t.Run("It returns accepted in POST", func(t *testing.T) {
		request := postRequestFor("Amily")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, 200)

		want := 1
		got, found := server.store.GetPlayerScore("Amily")
		if !found {
			t.Fatal("Expected to find player Amily but didn't find")
		}
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
	t.Run("Increments already existing players", func(t *testing.T) {
		request := postRequestFor("Robin")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, 200)

		want := 2
		got, found := server.store.GetPlayerScore("Robin")
		if !found {
			t.Fatal("Expected to find player Amily but didn't find")
		}
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

}

//Integration Test
func TestRecordingWinsAndRetrievingTheme(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
	player := "Pepper"
	player2 := "Rick"

	server.ServeHTTP(httptest.NewRecorder(), postRequestFor(player))
	server.ServeHTTP(httptest.NewRecorder(), postRequestFor(player2))
	server.ServeHTTP(httptest.NewRecorder(), postRequestFor(player))
	server.ServeHTTP(httptest.NewRecorder(), postRequestFor(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, getRequestFor(player))
		assertStatus(t, response.Code, http.StatusOK)
		assertPlayerScore(t, response.Body.String(), "3")

		response = httptest.NewRecorder()
		server.ServeHTTP(response, getRequestFor(player2))
		assertStatus(t, response.Code, http.StatusOK)
		assertPlayerScore(t, response.Body.String(), "1")
	})

	t.Run("get league", func(t *testing.T) {
		wantedLeague := []Player{{"Pepper", 3}, {"Rick", 1}}
		response := httptest.NewRecorder()
		server.ServeHTTP(response, getLeagueRequest())
		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)
		got := parseLeagueResponse(t, response.Body)
		assertLeague(t, got, wantedLeague)
	})

}

func TestLeagueEndpoint(t *testing.T) {
	store := SpyPlayerStore{
		scores: map[string]int{
			"Rick":  10,
			"Morty": 12,
		},
	}
	leaguePlayers := []Player{{"Rick", 10}, {"Morty", 12}}
	server := NewPlayerServer(store)

	req := getLeagueRequest()
	res := httptest.NewRecorder()

	server.ServeHTTP(res, req)

	assertStatus(t, res.Code, http.StatusOK)
	assertContentType(t, res, jsonContentType)

	got := parseLeagueResponse(t, res.Body)
	assertLeague(t, got, leaguePlayers)
}

type SpyPlayerStore struct {
	scores map[string]int
}

func (s SpyPlayerStore) GetPlayerScore(player string) (int, bool) {
	score, found := s.scores[player]
	return score, found
}

func (s SpyPlayerStore) IncrementPlayerScore(player string) {
	s.scores[player] = s.scores[player] + 1
}
func (s SpyPlayerStore) GetPlayersList() []Player {
	players := make([]Player, 0, len(s.scores))
	for name, score := range s.scores {
		players = append(players, Player{name, score})
	}
	return players
}

func getRequestFor(player string) *http.Request {
	req_url := fmt.Sprintf("/players/%s", player)
	req := httptest.NewRequest(http.MethodGet, req_url, nil)
	return req
}

func getLeagueRequest() *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/league/", nil)
	return req
}

func postRequestFor(player string) *http.Request {
	req_url := fmt.Sprintf("/players/%s", player)
	return httptest.NewRequest(http.MethodPost, req_url, nil)
}

func assertPlayerScore(t testing.TB, got, want string) {
	t.Helper()
	if want != got {
		t.Errorf("got %q, want %q", got, want)
	}
}

func parseLeagueResponse(t testing.TB, body io.Reader) []Player {
	t.Helper()
	players := []Player{}
	err := json.NewDecoder(body).Decode(&players)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}

	return players
}

func assertLeague(t testing.TB, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got status %d, want %d", got, want)
	}
}

func assertContentType(t testing.TB, res *httptest.ResponseRecorder, want string) {
	t.Helper()
	if res.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of application/json, got %v",
			res.Result().Header)
	}
}
