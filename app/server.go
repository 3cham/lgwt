package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var JsonContentType = "application/json"


type PlayerServer struct {
	store PlayerStore
	http.Handler
}

func (p *PlayerServer) playerHandler(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/players/"):]
	switch r.Method {
	case http.MethodGet:
		p.showScore(w, player)
	case http.MethodPost:
		p.updateScore(w, player)
	}
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(p.getPlayers())
	w.Header().Set("content-type", JsonContentType)
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score, err := p.store.GetPlayerScore(player)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		fmt.Fprint(w, score)
	}
}

func (p *PlayerServer) updateScore(w http.ResponseWriter, player string) {
	p.store.UpdatePlayerScore(player)
	w.WriteHeader(http.StatusAccepted)
}

func (p *PlayerServer) getPlayers() League {
	return p.store.GetPlayers()
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)
	p.store = store

	router := http.ServeMux{}
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playerHandler))

	p.Handler = &router
	return p
}