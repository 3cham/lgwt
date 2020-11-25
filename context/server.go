package context

/**
Incoming requests to a server should create a Context, and outgoing calls to
servers should accept a Context. The chain of function calls between them must
propagate the Context, optionally replacing it with a derived Context created
using WithCancel, WithDeadline, WithTimeout, or WithValue. When a Context is
canceled, all Contexts derived from it are also canceled.
 */

import (
	"context"
	"fmt"
	"net/http"
)

type Store interface {
	Fetch(ctx context.Context) (string, error)
	Cancel()
}

func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := store.Fetch(r.Context())
		if err == nil {
			fmt.Fprint(w, response)
		}
	}
}
