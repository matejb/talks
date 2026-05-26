package main

import (
	"log"
	"net/http"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/03-kennedy/internal/books"
	platformHTTP "github.com/matejb/talks/2026/gops-2026/layouts-examples/03-kennedy/internal/platform/http"
	"github.com/matejb/talks/2026/gops-2026/layouts-examples/03-kennedy/internal/reviews"
	"github.com/matejb/talks/2026/gops-2026/layouts-examples/03-kennedy/internal/users"
)

func main() {
	// Domain services - each with their own storage (no shared DB interface in this example).
	// In a real Kennedy layout, these would take *database.DB or *sql.DB directly.
	userSvc := users.NewService()

	// The books service needs a way to look up user names.
	// In the Kennedy pattern, domain packages can reference each other
	// at the composition root via a closure/function.
	bookSvc := books.NewService(func(id int) (string, error) {
		u, err := userSvc.User(id)
		if err != nil {
			return "", err
		}
		return u.Name, nil
	})
	reviewSvc := reviews.NewService()

	// HTTP handlers - each domain owns its own handlers.
	userHandler := &users.Handler{Service: userSvc}
	bookHandler := &books.Handler{Service: bookSvc}
	reviewHandler := &reviews.Handler{Service: reviewSvc}

	// Wire routes.
	mux := http.NewServeMux()
	userHandler.RegisterRoutes(mux)
	bookHandler.RegisterRoutes(mux)
	reviewHandler.RegisterRoutes(mux)

	// Start server via platform helper.
	log.Fatal(platformHTTP.Serve(":8080", mux))
}
