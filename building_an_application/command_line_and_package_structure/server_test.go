package poker

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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

		AssertPlayerScore(t, response.Body.String(), "20")
	})
	t.Run("check scores for Rick", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := getRequestFor("Rick")

		server.ServeHTTP(response, request)

		AssertStatus(t, response.Code, http.StatusOK)
		AssertPlayerScore(t, response.Body.String(), "0")
	})
	t.Run("Return 404 on missing players", func(t *testing.T) {
		request := getRequestFor("nonplayer")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response.Code, http.StatusNotFound)
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
		AssertStatus(t, response.Code, 200)
		AssertPlayerWin(t, server.store, "Amily", 1)
	})
	t.Run("Increments already existing players", func(t *testing.T) {
		request := postRequestFor("Robin")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		server.ServeHTTP(response, request)
		AssertStatus(t, response.Code, 200)
		AssertPlayerWin(t, server.store, "Robin", 2)
	})

}

func TestLeagueEndpoint(t *testing.T) {
	store := SpyPlayerStore{
		scores: map[string]int{
			"Rick":  10,
			"Morty": 12,
		},
	}
	// should return in sorted form
	leaguePlayers := []Player{{"Morty", 12}, {"Rick", 10}}
	server := NewPlayerServer(store)

	req := getLeagueRequest()
	res := httptest.NewRecorder()

	server.ServeHTTP(res, req)

	AssertStatus(t, res.Code, http.StatusOK)
	AssertContentType(t, res, jsonContentType)

	got := parseLeagueResponse(t, res.Body)
	AssertLeague(t, got, leaguePlayers)
}

func parseLeagueResponse(t testing.TB, body io.Reader) []Player {
	t.Helper()

	players, err := NewLeague(body)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
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
