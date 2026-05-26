package http

import (
	"log"
	"net/http"
)

// Serve starts an HTTP server on the given address using the provided mux.
// This is a platform-level helper - it knows nothing about business logic.
func Serve(addr string, mux *http.ServeMux) error {
	log.Printf("Kennedy layout example listening on %s", addr)
	return http.ListenAndServe(addr, mux)
}
