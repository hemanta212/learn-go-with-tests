package main

import (
	"context"
	"fmt"
	"net/http"
)

func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := store.Fetch(r.Context())
		if err != nil {
			return // TODO: log error etc
		}
		fmt.Fprintf(w, data)
	}
}

type Store interface {
	Fetch(ctx context.Context) (string, error)
}
