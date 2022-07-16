package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRacer(t *testing.T) {
	t.Run("Test fastest Server", func(t *testing.T) {
		slowServer := spinDelayedServer(20 * time.Millisecond)
		fastServer := spinDelayedServer(0 * time.Millisecond)

		defer slowServer.Close()
		defer fastServer.Close()

		slowUrl, fastUrl := slowServer.URL, fastServer.URL

		want := fastUrl
		got, err := Racer(slowUrl, fastUrl)

		if err != nil {
			t.Fatalf("did not expect an error but got one %v", err)
		}

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("Return error if no server respond within 10s", func(t *testing.T) {
		server := spinDelayedServer(11 * time.Millisecond)

		defer server.Close()

		_, err := ConfigurableRacer(server.URL, server.URL, 10*time.Millisecond)

		if err == nil {
			t.Errorf("expected an error but didn't get one")
		}
	})
}

func spinDelayedServer(sleepDuration time.Duration) (server *httptest.Server) {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(sleepDuration)
		w.WriteHeader(http.StatusOK)
	}))
	return
}
