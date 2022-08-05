package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const jsonContentType = "application/json"

type PlayerStore interface {
	GetPlayerScore(string) (int, bool)
	IncrementPlayerScore(string)
	GetLeague() League
}

type Player struct {
	Name  string
	Score int
}

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := &PlayerServer{
		store: store,
	}
	router := http.NewServeMux()
	router.Handle("/league/", http.HandlerFunc(p.leagueRoute))
	router.Handle("/players/", http.HandlerFunc(p.playerRoute))
	p.Handler = router
	return p
}

func (p *PlayerServer) playerRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.handleGet(w, r)
	} else if r.Method == http.MethodPost {
		p.handlePost(w, r)
	}
}

func (p *PlayerServer) handlePost(w http.ResponseWriter, r *http.Request) {
	playername := strings.TrimPrefix(r.URL.String(), "/players/")
	p.store.IncrementPlayerScore(playername)
	w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) handleGet(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.String(), "/players/")
	score, found := p.store.GetPlayerScore(player)
	if !found {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, score)
}

func (p *PlayerServer) leagueRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(p.store.GetLeague())
}
