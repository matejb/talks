package main

import (
	"log"
	"net/http"

	adapterhttp "github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/adapter/in/http"
	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/adapter/out/memory"
	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/core/service"
)

// THE ONLY PLACE that knows about concrete implementations.
// The core (domain, port, service) has zero knowledge of adapters.
func main() {
	// Outbound adapters (driven) - in-memory for this example
	bookRepo := memory.NewBookRepository()
	userRepo := memory.NewUserRepository()
	reviewRepo := memory.NewReviewRepository()

	// Core services (pure business logic)
	bookSvc := service.NewBookService(bookRepo, userRepo)
	userSvc := service.NewUserService(userRepo)
	reviewSvc := service.NewReviewService(reviewRepo)

	// Inbound adapters (driving)
	bookHandler := adapterhttp.NewBookHandler(bookSvc)
	userHandler := adapterhttp.NewUserHandler(userSvc)
	reviewHandler := adapterhttp.NewReviewHandler(reviewSvc)

	// Wire routes (Go 1.22+ pattern-based routing)
	mux := http.NewServeMux()
	for pattern, handler := range bookHandler.Routes() {
		mux.HandleFunc(pattern, handler)
	}
	for pattern, handler := range userHandler.Routes() {
		mux.HandleFunc(pattern, handler)
	}
	for pattern, handler := range reviewHandler.Routes() {
		mux.HandleFunc(pattern, handler)
	}

	log.Println("Hexagonal layout example running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
