package main

import "sync"

type PlayerStore interface {
	GetPlayerScore(name string) (int, error)
	UpdatePlayerScore(name string)
	GetPlayers() []Player
}

type InMemoryStore struct {
	lock sync.Mutex
	scores map[string]int
}
func (s *InMemoryStore) GetPlayerScore(name string) (int, error) {
	score, found := s.scores[name]
	if !found {
		return 0, ErrPlayerNotFound
	}

	return score, nil
}

func (s *InMemoryStore) UpdatePlayerScore(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	currentScore, found := s.scores[name]

	if !found {
		s.scores[name] = 1
	} else {
		s.scores[name] = currentScore + 1
	}
}

func (s *InMemoryStore) GetPlayers() []Player {
	players := []Player{}
	for player, score := range s.scores {
		players = append(players, Player{player, score})
	}
	return players
}
