package context

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type SpyStore struct {
	response  string
	cancelled bool
	t         *testing.T
}

func (s *SpyStore) Fetch() string {
	time.Sleep(100 * time.Millisecond)
	return s.response
}

func (s *SpyStore) Cancel() {
	s.cancelled = true
}

func (s *SpyStore) assertShouldNotCancel() {
	s.t.Helper()
	if s.cancelled {
		s.t.Fatalf("Request should not be cancelled")
	}
}

func (s *SpyStore) assertShouldCancel() {
	s.t.Helper()
	if !s.cancelled {
		s.t.Fatalf("Request should be cancelled")
	}
}

func TestHandler(t *testing.T) {
	data := "Hello, world"

	t.Run("StubStore should return correct string", func(t *testing.T) {
		store := &SpyStore{data, false, t}
		s := Server(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		s.ServeHTTP(response, request)

		store.assertShouldNotCancel()
		if response.Body.String() != data {
			t.Fatalf("Wrong result, expected %s, got %s", data, response.Body.String())
		}
	})

	t.Run("Cancel request should return empty string", func(t *testing.T) {
		store := &SpyStore{data, false, t}
		s := Server(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		cancelingCtx, cancel := context.WithCancel(request.Context())
		time.AfterFunc(5*time.Millisecond, cancel)

		request = request.WithContext(cancelingCtx)
		s.ServeHTTP(response, request)

		store.assertShouldCancel()
	})
}
