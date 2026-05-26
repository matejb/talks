package main

import (
	"log"
	"net/http"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/08-monorepo/libs/library"
	"github.com/matejb/talks/2026/gops-2026/layouts-examples/08-monorepo/services/api/internal/handler"
)

func main() {
	store := library.NewStore()

	h := &handler.Handler{Store: store}
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	log.Println("Monorepo API service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
