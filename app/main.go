package main

import (
	"errors"
	"log"
	"net/http"
	"os"
)

var (
	ErrPlayerNotFound = errors.New("Player not found")
	dbFileName = "game.db.json"
)

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 666)

	if err != nil {
		log.Fatalf("Cannot open file %s for writing database", dbFileName)
	}
	defer db.Close()
	store := NewFileSystemPlayerStore(db)
	server := NewPlayerServer(&store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}