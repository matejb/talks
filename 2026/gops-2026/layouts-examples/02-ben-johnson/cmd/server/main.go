package main

import (
	"log"
	"net/http"

	httpHandler "github.com/matejb/talks/2026/gops-2026/layouts-examples/02-ben-johnson/http"
	"github.com/matejb/talks/2026/gops-2026/layouts-examples/02-ben-johnson/memory"
)

func main() {
	// Create in-memory implementations of domain services.
	us := memory.NewUserService()
	bs := memory.NewBookService(us)
	rs := memory.NewReviewService()

	// Wire HTTP handlers to the domain interfaces.
	h := &httpHandler.Handler{
		Users:   us,
		Books:   bs,
		Reviews: rs,
	}

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	log.Println("Ben Johnson Standard Package Layout example running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
