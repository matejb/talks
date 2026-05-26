package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/06-feature-based/internal/book"
	"github.com/matejb/talks/2026/gops-2026/layouts-examples/06-feature-based/internal/review"
	"github.com/matejb/talks/2026/gops-2026/layouts-examples/06-feature-based/internal/user"
)

func main() {
	// --- Wire repositories (in-memory for this example) ---
	userRepo := user.NewMemoryRepository()
	bookRepo := book.NewMemoryRepository()
	reviewRepo := review.NewMemoryRepository()

	// --- Wire services ---
	userSvc := user.NewService(userRepo)

	// book.Service needs a user lookup - bridge between features
	userLookup := func(ctx context.Context, userID int) error {
		u, err := userSvc.User(ctx, userID)
		if err != nil {
			return fmt.Errorf("user %d not found: %w", userID, err)
		}
		_ = u // user exists
		return nil
	}
	bookSvc := book.NewService(bookRepo, userLookup)

	// review.Service needs a book lookup
	bookLookup := func(ctx context.Context, bookID int) error {
		b, err := bookSvc.Book(ctx, bookID)
		if err != nil {
			return fmt.Errorf("book %d not found: %w", bookID, err)
		}
		_ = b // book exists
		return nil
	}
	reviewSvc := review.NewService(reviewRepo, bookLookup)

	// --- Wire handlers ---
	userHandler := user.NewHandler(userSvc)
	bookHandler := book.NewHandler(bookSvc)
	reviewHandler := review.NewHandler(reviewSvc)

	// --- Register routes (Go 1.22+ pattern-based routing) ---
	mux := http.NewServeMux()
	for pattern, handler := range userHandler.Routes() {
		mux.HandleFunc(pattern, handler)
	}
	for pattern, handler := range bookHandler.Routes() {
		mux.HandleFunc(pattern, handler)
	}
	for pattern, handler := range reviewHandler.Routes() {
		mux.HandleFunc(pattern, handler)
	}

	log.Println("Feature-based layout example running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
