package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	ErrPlayerNotFound = errors.New("Player not found")
)

type PlayerStore interface {
	GetPlayerScore(name string) (int, error)
	UpdatePlayerScore(name string)
}

type PlayerServer struct {
	store PlayerStore
}

type InMemoryStore struct {}
func (s *InMemoryStore) GetPlayerScore(name string) (int, error) {
	return 123, nil
}

func (s *InMemoryStore) UpdatePlayerScore(name string) {
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	switch r.Method {
	case http.MethodGet:
		p.showScore(w, player)
	case http.MethodPost:
		p.updateScore(w, player)
	}
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

func main() {
	server := &PlayerServer{&InMemoryStore{}}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}