package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPlayerStore struct {
	scores map[string]int
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	return s.scores[name]
}

func TestGETPlayerScore(t *testing.T) {
	scoresMap := map[string]int {
		"A": 20,
		"B": 10,
	}
	server := &PlayerServer{&StubPlayerStore{scoresMap}}
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
}

func assertRespondedScore(t *testing.T, response *httptest.ResponseRecorder, expected string) {
	t.Helper()
	got := response.Body.String()
	if got != expected {
		t.Fatalf("Wrong score for player returned, expected %s, got %s", expected, got)
	}
}

func getScoreRequest(player string) (*http.Request, error) {
	return http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", player), nil)
}
