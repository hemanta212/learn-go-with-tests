package poker_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	poker "github.com/hemanta212/go-with-tdd-project"
)

var (
	dummyGame = &poker.SpyGame{}
)

func TestGetPlayers(t *testing.T) {
	store := poker.SpyPlayerStore{
		Scores: map[string]int{
			"Amily": 20,
			"Rick":  0,
		},
	}
	server := mustMakePlayerServer(t, store, dummyGame)

	t.Run("check scores for Amily", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := getRequestFor("Amily")

		server.ServeHTTP(response, request)

		poker.AssertPlayerScore(t, response.Body.String(), "20")
	})
	t.Run("check scores for Rick", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := getRequestFor("Rick")

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response.Code, http.StatusOK)
		poker.AssertPlayerScore(t, response.Body.String(), "0")
	})
	t.Run("Return 404 on missing players", func(t *testing.T) {
		request := getRequestFor("nonplayer")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestPostPlayerWins(t *testing.T) {
	store := poker.SpyPlayerStore{
		Scores: map[string]int{},
	}
	server := mustMakePlayerServer(t, store, dummyGame)

	t.Run("It returns accepted in POST", func(t *testing.T) {
		request := postRequestFor("Amily")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		poker.AssertStatus(t, response.Code, 200)
		poker.AssertPlayerWin(t, &store, "Amily", 1)
	})
	t.Run("Increments already existing players", func(t *testing.T) {
		request := postRequestFor("Robin")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		server.ServeHTTP(response, request)
		poker.AssertStatus(t, response.Code, 200)
		poker.AssertPlayerWin(t, &store, "Robin", 2)
	})

}

func TestLeagueEndpoint(t *testing.T) {
	store := poker.SpyPlayerStore{
		Scores: map[string]int{
			"Rick":  10,
			"Morty": 12,
		},
	}
	// should return in sorted form
	leaguePlayers := []poker.Player{{"Morty", 12}, {"Rick", 10}}
	server := mustMakePlayerServer(t, store, dummyGame)

	req := getLeagueRequest()
	res := httptest.NewRecorder()

	server.ServeHTTP(res, req)

	poker.AssertStatus(t, res.Code, http.StatusOK)
	poker.AssertContentType(t, res, "application/json")

	got := parseLeagueResponse(t, res.Body)
	poker.AssertLeague(t, got, leaguePlayers)
}

func TestGame(t *testing.T) {
	t.Run("Get /game returns 200", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/game", nil)
		response := httptest.NewRecorder()

		server := mustMakePlayerServer(t, poker.NewSpyPlayerStore(), dummyGame)
		server.ServeHTTP(response, request)

		poker.AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("start game with 3 players, send some blind alerts down the WS and declare ruth the winner", func(t *testing.T) {
		wantedBlindAlert := "Blind is 100"
		winner := "Ruth"

		game := &poker.SpyGame{BlindAlert: []byte(wantedBlindAlert)}
		server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
		ws := mustDialWS(t, wsURL)

		defer server.Close()
		defer ws.Close()

		writeWSMessage(t, ws, "3")
		writeWSMessage(t, ws, winner)

		assertGamePlayerNo(t, game, 3)
		assertGameWinner(t, game, winner)

		within(t, 10*time.Millisecond, func() {
			assertWebsocketGotMsg(t, ws, wantedBlindAlert)
		})
	})

}

func parseLeagueResponse(t testing.TB, body io.Reader) []poker.Player {
	t.Helper()

	players, err := poker.NewLeague(body)

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

func mustMakePlayerServer(t *testing.T, store poker.PlayerStore, game poker.Game) *poker.PlayerServer {
	t.Helper()
	server, err := poker.NewPlayerServer(store, game)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}
	return server
}

func mustDialWS(t *testing.T, url string) *websocket.Conn {
	t.Helper()
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", url, err)
	}

	return ws
}

func writeWSMessage(t testing.TB, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("couldn't send message over ws connection %v", err)
	}
}

func within(t testing.TB, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("timed Out")
	case <-done:

	}
}

func assertWebsocketGotMsg(t *testing.T, ws *websocket.Conn, want string) {
	_, msg, _ := ws.ReadMessage()
	if string(msg) != want {
		t.Errorf("got %q, want %q", string(msg), want)
	}
}
