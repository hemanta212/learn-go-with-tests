package main

import (
	"fmt"
	"net/http"
	"strings"
)

type PlayerStore interface {
	GetPlayerScore(string) (int, bool)
	IncrementPlayerScore(string)
}

type PlayerServer struct {
	Store PlayerStore
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.handleGet(w, r)
	} else if r.Method == http.MethodPost {
		p.handlePost(w, r)
	}
}

func (p *PlayerServer) handlePost(w http.ResponseWriter, r *http.Request) {
	playername := strings.TrimPrefix(r.URL.String(), "/players/")
	p.Store.IncrementPlayerScore(playername)
	w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) handleGet(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.String(), "/players/")
	score, found := p.Store.GetPlayerScore(player)
	if !found {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, score)
}
