package context

import (
	"context"
	"errors"
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

func (s *SpyStore) Fetch(ctx context.Context) (string, error) {
	data := make(chan string, 1)

	go func() {
		var result string
		for _, c := range s.response {
			select {
			case <- ctx.Done():
				s.t.Log("spy store got canceled")
				s.cancelled = true
				return
			default:
				time.Sleep(10 * time.Millisecond)
				result += string(c)
			}
		}
		data <- result
	}()
	select {
	case res := <-data:
		return res, nil
	case <-ctx.Done():
		s.cancelled = true
		return "", ctx.Err()
	}
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

type SpyResponseWriter struct {
	written bool
}

func (s *SpyResponseWriter) Header() http.Header {
	s.written = true
	return nil
}

func (s *SpyResponseWriter) Write([]byte) (int, error) {
	s.written = true
	return 0, errors.New("Not implemented")
}

func (s *SpyResponseWriter) WriteHeader(statusCode int) {
	s.written = true
}

func TestHandler(t *testing.T) {
	data := "Hello, world"

	t.Run("SpyStore should return correct string", func(t *testing.T) {
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
		response := &SpyResponseWriter{}

		cancelingCtx, cancel := context.WithCancel(request.Context())
		time.AfterFunc(5*time.Millisecond, cancel)

		request = request.WithContext(cancelingCtx)
		s.ServeHTTP(response, request)

		store.assertShouldCancel()
		if response.written {
			t.Fatalf("A response should not be written")
		}
	})
}
