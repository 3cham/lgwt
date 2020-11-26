package main

import (
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
	server := &PlayerServer{&StubPlayerStore{scoresMap, nil}}
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
	store := StubPlayerStore{make(map[string]int), []string{}}
	server := &PlayerServer{&store}

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
	store := InMemoryStore{sync.Mutex{}, make(map[string]int)}
	server := &PlayerServer{&store}
	player := "C"

	t.Run("Update and show score of the same play should return consistent result", func(t *testing.T) {
		request, _ := postScoreRequest(player)
		// Post 3 times win for C
		server.ServeHTTP(httptest.NewRecorder(), request)
		server.ServeHTTP(httptest.NewRecorder(), request)
		server.ServeHTTP(httptest.NewRecorder(), request)

		// Get C's score
		request, _ = getScoreRequest(player)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertRespondedStatusCode(t, response, http.StatusOK)
		assertRespondedScore(t, response, "3")
	})
}

func TestLeague(t *testing.T) {
	store := StubPlayerStore{}
	server := PlayerServer{&store}

	t.Run("/league should return list of players", func(t *testing.T) {
		request, _ := postScoreRequest("A")
		server.ServeHTTP(httptest.NewRecorder(), request)
		request, _ = postScoreRequest("B")
		server.ServeHTTP(httptest.NewRecorder(), request)
		request, _ = postScoreRequest("C")
		server.ServeHTTP(httptest.NewRecorder(), request)

		request, _ = http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertRespondedStatusCode(t, response, http.StatusOK)
	})
}