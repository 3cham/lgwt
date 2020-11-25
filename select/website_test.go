package website_racer

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRacer(t *testing.T) {
	t.Run("Race should return fastHttpServer", func(t *testing.T) {
		slowHttpServer := makeHTTPServer(20 * time.Millisecond)
		fastHttpServer := makeHTTPServer(0 * time.Millisecond)

		defer fastHttpServer.Close()
		defer slowHttpServer.Close()

		fastUrl := fastHttpServer.URL
		slowURL := slowHttpServer.URL

		got, _ := Racer(fastUrl, slowURL)
		expected := fastHttpServer.URL

		if got != expected {
			t.Fatalf("Wrong URL returned, got %q, expected %q", got, expected)
		}
	})

	t.Run("Racer should throw an error if request takes longer than 10 milliseconds", func(t *testing.T) {
		slowHttpServer := makeHTTPServer(12 * time.Millisecond)
		fastHttpServer := makeHTTPServer(11 * time.Millisecond)

		defer fastHttpServer.Close()
		defer slowHttpServer.Close()
		timeout := 10 * time.Millisecond

		_, err := ConfigurableRacer(fastHttpServer.URL, slowHttpServer.URL, timeout)

		if err == nil {
			t.Fatalf("Error expected")
		}
	})
}

func makeHTTPServer(delayedTime time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delayedTime)
		w.WriteHeader(200)
	}))
}
