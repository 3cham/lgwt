package context

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type StubStore struct {
	response string
}

func (s *StubStore) Fetch() string {
	return s.response
}

func (s *StubStore) Cancel() {
}

type SpyStore struct {
	response  string
	cancelled bool
}

func (s *SpyStore) Fetch() string {
	time.Sleep(100 * time.Millisecond)
	return s.response
}

func (s *SpyStore) Cancel() {
	s.cancelled = true
}

func TestHandler(t *testing.T) {
	t.Run("StubStore should return correct string", func(t *testing.T) {
		data := "Hello, world"
		store := &SpyStore{data, false}
		s := Server(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		s.ServeHTTP(response, request)

		if store.cancelled {
			t.Fatalf("Request should not be cancelled")
		}
		if response.Body.String() != data {
			t.Fatalf("Wrong result, expected %s, got %s", data, response.Body.String())
		}
	})

	t.Run("Cancel request should return empty string", func(t *testing.T) {
		data := "Hello, world"
		store := &SpyStore{data, false}
		s := Server(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		cancelingCtx, cancel := context.WithCancel(request.Context())
		time.AfterFunc(5*time.Millisecond, cancel)

		request = request.WithContext(cancelingCtx)
		s.ServeHTTP(response, request)

		if !store.cancelled {
			t.Fatalf("Store is not cancelled")
		}
	})

}
