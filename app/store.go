package main

import (
	"encoding/json"
	"io"
)

type PlayerStore interface {
	GetPlayerScore(name string) (int, error)
	UpdatePlayerScore(name string)
	GetPlayers() League
}

type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
	league League
}

func (f *FileSystemPlayerStore) GetPlayers() League {
	return f.league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) (int, error) {
	players := f.GetPlayers()

	player := players.Find(name)
	if player == nil {
		return 0, ErrPlayerNotFound
	} else {
		return player.Wins, nil
	}
}

func (f *FileSystemPlayerStore) UpdatePlayerScore(name string) {
	players := f.GetPlayers()
	player := players.Find(name)

	if player == nil {
		players = append(players, Player{name, 1})
	} else {
		player.Wins++
	}

	f.league = players
	f.database.Seek(0, 0)
	json.NewEncoder(f.database).Encode(players)
}

func NewLeague(r io.ReadWriteSeeker) *League {
	result := League{}
	r.Seek(0, 0)
	json.NewDecoder(r).Decode(&result)
	return &result
}

func NewFileSystemPlayerStore(r io.ReadWriteSeeker) FileSystemPlayerStore{
	store := FileSystemPlayerStore{r, nil}
	store.league = *NewLeague(store.database)
	return store
}
