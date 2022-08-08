package poker_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	poker "github.com/hemanta212/go-with-tdd-project"
)

//Integration Test
func TestRecordingWinsAndRetrievingTheme(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[]`)
	defer cleanDatabase()
	store, err := poker.NewFileSystemPlayerStore(database)
	poker.AssertNoError(t, err)

	server := mustMakePlayerServer(t, store, dummyGame)

	player := "Pepper"
	player2 := "Rick"
	player3 := "Monty"

	server.ServeHTTP(httptest.NewRecorder(), postRequestFor(player))
	server.ServeHTTP(httptest.NewRecorder(), postRequestFor(player2))
	server.ServeHTTP(httptest.NewRecorder(), postRequestFor(player))
	server.ServeHTTP(httptest.NewRecorder(), postRequestFor(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, getRequestFor(player))
		poker.AssertStatus(t, response.Code, http.StatusOK)
		poker.AssertPlayerScore(t, response.Body.String(), "3")

		response = httptest.NewRecorder()
		server.ServeHTTP(response, getRequestFor(player2))
		poker.AssertStatus(t, response.Code, http.StatusOK)
		poker.AssertPlayerScore(t, response.Body.String(), "1")

		response = httptest.NewRecorder()
		server.ServeHTTP(response, getRequestFor(player3))
		poker.AssertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("get league", func(t *testing.T) {
		wantedLeague := poker.League{{"Pepper", 3}, {"Rick", 1}}
		response := httptest.NewRecorder()
		server.ServeHTTP(response, getLeagueRequest())
		poker.AssertStatus(t, response.Code, http.StatusOK)
		poker.AssertContentType(t, response, "application/json")
		got := parseLeagueResponse(t, response.Body)
		poker.AssertLeague(t, got, wantedLeague)
	})

}
