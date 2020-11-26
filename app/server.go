package main

import (
	"fmt"
	"net/http"
)


type PlayerServer struct {
	store PlayerStore
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	router := http.ServeMux{}

	router.Handle("/league", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	router.Handle("/players/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		player := r.URL.Path[len("/players/"):]
		switch r.Method {
		case http.MethodGet:
			p.showScore(w, player)
		case http.MethodPost:
			p.updateScore(w, player)
		}
	}))

	router.ServeHTTP(w, r)
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
