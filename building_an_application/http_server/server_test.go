package main

import (
	"fmt"
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
	server := &PlayerServer{Store: store}

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
	server := &PlayerServer{store}

	t.Run("It returns accepted in POST", func(t *testing.T) {
		request := postRequestFor("Amily")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, 200)

		want := 1
		got, found := server.Store.GetPlayerScore("Amily")
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
		got, found := server.Store.GetPlayerScore("Robin")
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
	server := PlayerServer{store}
	player := "pepper"

	server.ServeHTTP(httptest.NewRecorder(), postRequestFor(player))
	server.ServeHTTP(httptest.NewRecorder(), postRequestFor(player))
	server.ServeHTTP(httptest.NewRecorder(), postRequestFor(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, getRequestFor(player))
	assertStatus(t, response.Code, http.StatusOK)

	assertPlayerScore(t, response.Body.String(), "3")

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

func getRequestFor(player string) *http.Request {
	req_url := fmt.Sprintf("/players/%s", player)
	req := httptest.NewRequest(http.MethodGet, req_url, nil)
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

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got status %d, want %d", got, want)
	}
}
