package sync

import (
	"sync"
	"testing"
)

func TestCount(t *testing.T) {

	t.Run("Incrementing counter 3 times leaves it at 3", func(t *testing.T) {
		counter := Counter{sync.Mutex{}, 0}

		counter.Inc()
		counter.Inc()
		counter.Inc()

		assertCounter(t, counter, 3)
	})
	t.Run("Incrementing counter concurrently", func(t *testing.T) {
		counter := Counter{sync.Mutex{}, 0}
		threads := 300

		var wg sync.WaitGroup
		wg.Add(threads)

		for i := 0; i < threads; i++ {
			go func(w *sync.WaitGroup) {
				counter.Inc()
				w.Done()
			}(&wg)
		}
		wg.Wait()
		assertCounter(t, counter, threads)
	})
}

func assertCounter(t *testing.T, counter Counter, expected int) {
	t.Helper()
	if counter.Value != expected {
		t.Fatalf("Wrong value for counter after incrementing, got %d, expected %d", counter.Value, expected)
	}
}
