package main

import (
	"errors"
	"log"
	"net/http"
	"sync"
)

var (
	ErrPlayerNotFound = errors.New("Player not found")
)

func main() {
	server := &PlayerServer{&InMemoryStore{sync.Mutex{}, make(map[string]int)}}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}