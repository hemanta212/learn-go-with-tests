package poker

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
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
	template *template.Template
	game     Game
}

const htmlTemplatePath = "game.html"

func NewPlayerServer(store PlayerStore, game Game) (*PlayerServer, error) {
	p := &PlayerServer{
		store: store,
	}

	templ, err := template.ParseFiles(htmlTemplatePath)
	if err != nil {
		return nil, fmt.Errorf("problem opening %s %v", htmlTemplatePath, err)
	}

	p.template = templ
	p.store = store
	p.game = game

	router := http.NewServeMux()
	router.Handle("/league/", http.HandlerFunc(p.leagueRoute))
	router.Handle("/players/", http.HandlerFunc(p.playerRoute))
	router.Handle("/game", http.HandlerFunc(p.gameRoute))
	router.Handle("/ws", http.HandlerFunc(p.webSocket))

	p.Handler = router

	return p, nil
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

func (p *PlayerServer) gameRoute(w http.ResponseWriter, r *http.Request) {
	p.template.Execute(w, nil)
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type PlayerServerWS struct {
	*websocket.Conn
}

func newPlayerServerWS(w http.ResponseWriter, r *http.Request) *PlayerServerWS {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("problem upgrading connection to WebSockets %v\n", err)
	}
	return &PlayerServerWS{conn}
}

func (w *PlayerServerWS) WaitForMsg() string {
	_, msg, err := w.ReadMessage()
	if err != nil {
		log.Printf("error reading from websocket %v\n", err)
	}
	return string(msg)
}

func (w *PlayerServerWS) Write(p []byte) (int, error) {
	err := w.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	ws := newPlayerServerWS(w, r)

	numberOfPlayerMsg := ws.WaitForMsg()
	numberOfPlayers, _ := strconv.Atoi(string(numberOfPlayerMsg))
	p.game.Start(numberOfPlayers, ws) // todo: dont discard blinds alerts

	winner := ws.WaitForMsg()
	p.game.Finish(winner)
}
