package context

/**
Incoming requests to a server should create a Context, and outgoing calls to
servers should accept a Context. The chain of function calls between them must
propagate the Context, optionally replacing it with a derived Context created
using WithCancel, WithDeadline, WithTimeout, or WithValue. When a Context is
canceled, all Contexts derived from it are also canceled.
 */

import (
	"fmt"
	"net/http"
)

type Store interface {
	Fetch() string
	Cancel()
}

func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		data := make(chan string, 1)

		go func() {
			data <- store.Fetch()
		}()

		select {
		case d := <-data:
			fmt.Fprint(w, d)
		case <-ctx.Done():
			store.Cancel()
		}
	}
}
