package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
)

type StubPlayerStore struct {
	scores map[string]int
	updateScoreCalls []string
	players []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) (int, error) {
	val, found := s.scores[name]
	if !found {
		return 0, ErrPlayerNotFound
	}
	return val, nil
}

func (s *StubPlayerStore) UpdatePlayerScore(name string) {
	s.updateScoreCalls = append(s.updateScoreCalls, name)
}

func (s *StubPlayerStore) GetPlayers() []Player {
	return s.players
}

func assertRespondedScore(t *testing.T, response *httptest.ResponseRecorder, expected string) {
	t.Helper()
	got := response.Body.String()
	if got != expected {
		t.Fatalf("Wrong score for player returned, expected %s, got %s", expected, got)
	}
}

func assertRespondedStatusCode(t *testing.T, response *httptest.ResponseRecorder, expected int) {
	t.Helper()
	got := response.Code
	if got != expected {
		t.Fatalf("Wrong status code returned, expected %d, got %d", expected, got)
	}
}

func getScoreRequest(player string) (*http.Request, error) {
	return http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", player), nil)
}

func postScoreRequest(player string) (*http.Request, error) {
	return http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", player), nil)
}


func TestGETPlayerScore(t *testing.T) {
	scoresMap := map[string]int {
		"A": 20,
		"B": 10,
	}
	server := NewPlayerServer(&StubPlayerStore{scoresMap, nil, nil})
	t.Run("Return correct score for A", func(t *testing.T) {
		request, _ := getScoreRequest("A")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertRespondedScore(t, response, "20")
	})

	t.Run("Return correct score for B", func(t *testing.T) {
		request, _ := getScoreRequest("B")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertRespondedScore(t, response, "10")
	})

	t.Run("Return 404 if player not found", func(t *testing.T) {
		request, _ := getScoreRequest("C")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertRespondedStatusCode(t, response, http.StatusNotFound)
	})
}

func TestUpdatePlayerStore(t *testing.T) {
	store := new(StubPlayerStore)
	server := NewPlayerServer(store)

	t.Run("POST score should return accepted", func(t *testing.T) {
		request, _ := postScoreRequest("C")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertRespondedStatusCode(t, response, http.StatusAccepted)

		expectedScoreCalls := []string{"C"}
		if !reflect.DeepEqual(store.updateScoreCalls, expectedScoreCalls) {
			t.Fatalf("wrong update score calls, got %v, expected %v", store.updateScoreCalls, expectedScoreCalls)
		}
	})
}

func TestUpdateAndShowPlayerScore(t *testing.T) {
	store := &InMemoryStore{sync.Mutex{}, make(map[string]int)}
	server := NewPlayerServer(store)
	player := "C"

	request, _ := postScoreRequest(player)
	// Post 3 times win for C
	server.ServeHTTP(httptest.NewRecorder(), request)
	server.ServeHTTP(httptest.NewRecorder(), request)
	server.ServeHTTP(httptest.NewRecorder(), request)

	t.Run("Update and show score of the same play should return consistent result", func(t *testing.T) {
		// Get C's score
		request, _ = getScoreRequest(player)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertRespondedStatusCode(t, response, http.StatusOK)
		assertRespondedScore(t, response, "3")
	})

	t.Run("Get league shoud return correct result", func(t *testing.T) {
		request, _ = http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := []Player{}

		json.NewDecoder(response.Body).Decode(&got)
		expected := []Player{
			{"C", 3},
		}

		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("Wrong league table returned, expected %v, got %v", expected, got)
		}
	})
}

func TestLeague(t *testing.T) {
	store := StubPlayerStore{}
	server := NewPlayerServer(&store)

	t.Run("/league should return list of players", func(t *testing.T) {
		request, _ := postScoreRequest("A")
		server.ServeHTTP(httptest.NewRecorder(), request)
		request, _ = postScoreRequest("B")
		server.ServeHTTP(httptest.NewRecorder(), request)
		request, _ = postScoreRequest("C")
		server.ServeHTTP(httptest.NewRecorder(), request)

		expected := []Player{
			{"A", 1},
			{"B", 1},
			{"C", 1},
		}
		store.players = expected

		request, _ = http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assertRespondedStatusCode(t, response, http.StatusOK)
		assertResponsedType(t, response, JsonContentType)
		var got []Player
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("Cannot parse the returned result from the endpoint")
		}

		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("Return player list is incorrect, expected %v, got %v", expected, got)
		}
	})
}

func assertResponsedType(t *testing.T, response *httptest.ResponseRecorder, expectedType string) {
	gotType := response.Header().Get("content-type")
	if gotType != expectedType {
		t.Fatalf("Wrong response header for content-type, expected %q, got %q", expectedType, gotType)
	}
}