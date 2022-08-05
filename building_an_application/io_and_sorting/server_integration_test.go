package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

//Integration Test
func TestRecordingWinsAndRetrievingTheme(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[]`)
	defer cleanDatabase()
	store, err := NewFileSystemPlayerStore(database)
	assertNoError(t, err)

	server := NewPlayerServer(store)
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
		assertStatus(t, response.Code, http.StatusOK)
		assertPlayerScore(t, response.Body.String(), "3")

		response = httptest.NewRecorder()
		server.ServeHTTP(response, getRequestFor(player2))
		assertStatus(t, response.Code, http.StatusOK)
		assertPlayerScore(t, response.Body.String(), "1")

		response = httptest.NewRecorder()
		server.ServeHTTP(response, getRequestFor(player3))
		assertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("get league", func(t *testing.T) {
		wantedLeague := League{{"Pepper", 3}, {"Rick", 1}}
		response := httptest.NewRecorder()
		server.ServeHTTP(response, getLeagueRequest())
		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)
		got := parseLeagueResponse(t, response.Body)
		assertLeague(t, got, wantedLeague)
	})

}
