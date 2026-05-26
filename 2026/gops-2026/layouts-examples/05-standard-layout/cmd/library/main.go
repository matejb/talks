package main

import (
	"log"
	"net/http"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/05-standard-layout/internal/app/library"
)

func main() {
	userStore := library.NewMemoryUserStore()
	bookStore := library.NewMemoryBookStore()
	reviewStore := library.NewMemoryReviewStore()

	h := &library.Handler{
		Users:   &library.UserService{Store: userStore},
		Books:   &library.BookService{Store: bookStore, Users: userStore},
		Reviews: &library.ReviewService{Store: reviewStore},
	}

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	log.Println("Standard layout example running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
