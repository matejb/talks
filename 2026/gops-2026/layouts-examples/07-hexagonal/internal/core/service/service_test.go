package service_test

import (
	"context"
	"testing"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/adapter/out/memory"
	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/core/domain"
	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/core/service"
)

// TestBorrowBook tests the Borrow use case using in-memory adapters.
// No database needed - this is the killer feature of hexagonal architecture.
func TestBorrowBook(t *testing.T) {
	bookRepo := memory.NewBookRepository()
	userRepo := memory.NewUserRepository()
	svc := service.NewBookService(bookRepo, userRepo)

	ctx := context.Background()

	// Book 1 and User 1 are pre-seeded in the memory adapters.
	err := svc.Borrow(ctx, 1, 1)
	if err != nil {
		t.Fatalf("Borrow: %v", err)
	}

	book, err := svc.Get(ctx, 1)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if book.Borrower == nil {
		t.Fatal("expected book to be borrowed")
	}
	if book.Borrower.Name != "Alice" {
		t.Fatalf("expected borrower Alice, got %s", book.Borrower.Name)
	}
}

func TestBorrowAlreadyBorrowedBook(t *testing.T) {
	bookRepo := memory.NewBookRepository()
	userRepo := memory.NewUserRepository()
	svc := service.NewBookService(bookRepo, userRepo)

	ctx := context.Background()

	// Borrow the book first
	if err := svc.Borrow(ctx, 1, 1); err != nil {
		t.Fatalf("first Borrow: %v", err)
	}

	// Second borrow should fail
	err := svc.Borrow(ctx, 1, 2)
	if err == nil {
		t.Fatal("expected error when borrowing already borrowed book")
	}
}

func TestBorrowNonExistentBook(t *testing.T) {
	bookRepo := memory.NewBookRepository()
	userRepo := memory.NewUserRepository()
	svc := service.NewBookService(bookRepo, userRepo)

	ctx := context.Background()

	err := svc.Borrow(ctx, 999, 1)
	if err == nil {
		t.Fatal("expected error for non-existent book")
	}
}

func TestBorrowNonExistentUser(t *testing.T) {
	bookRepo := memory.NewBookRepository()
	userRepo := memory.NewUserRepository()
	svc := service.NewBookService(bookRepo, userRepo)

	ctx := context.Background()

	err := svc.Borrow(ctx, 1, 999)
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestReturnBook(t *testing.T) {
	bookRepo := memory.NewBookRepository()
	userRepo := memory.NewUserRepository()
	svc := service.NewBookService(bookRepo, userRepo)

	ctx := context.Background()

	// Borrow first
	if err := svc.Borrow(ctx, 1, 1); err != nil {
		t.Fatalf("Borrow: %v", err)
	}

	// Return
	if err := svc.Return(ctx, 1); err != nil {
		t.Fatalf("Return: %v", err)
	}

	book, err := svc.Get(ctx, 1)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if book.Borrower != nil {
		t.Fatal("expected book to be returned (Borrower should be nil)")
	}
}

func TestReturnNotBorrowedBook(t *testing.T) {
	bookRepo := memory.NewBookRepository()
	userRepo := memory.NewUserRepository()
	svc := service.NewBookService(bookRepo, userRepo)

	ctx := context.Background()

	err := svc.Return(ctx, 1)
	if err == nil {
		t.Fatal("expected error when returning a book that is not borrowed")
	}
}

func TestGetBook(t *testing.T) {
	bookRepo := memory.NewBookRepository()
	userRepo := memory.NewUserRepository()
	svc := service.NewBookService(bookRepo, userRepo)

	ctx := context.Background()

	book, err := svc.Get(ctx, 1)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if book.Title != "The Go Programming Language" {
		t.Fatalf("unexpected title: %s", book.Title)
	}
}

func TestAllBooks(t *testing.T) {
	bookRepo := memory.NewBookRepository()
	userRepo := memory.NewUserRepository()
	svc := service.NewBookService(bookRepo, userRepo)

	ctx := context.Background()

	books, err := svc.All(ctx)
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(books) != 2 {
		t.Fatalf("expected 2 books, got %d", len(books))
	}
}

func TestUserService(t *testing.T) {
	userRepo := memory.NewUserRepository()
	svc := service.NewUserService(userRepo)

	ctx := context.Background()

	// Get pre-seeded user
	u, err := svc.Get(ctx, 1)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if u.Name != "Alice" {
		t.Fatalf("expected Alice, got %s", u.Name)
	}

	// List users
	users, err := svc.All(ctx)
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}

	// Create user
	err = svc.Create(ctx, domain.User{Name: "Charlie", Email: "charlie@example.com"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	users, err = svc.All(ctx)
	if err != nil {
		t.Fatalf("All after create: %v", err)
	}
	if len(users) != 3 {
		t.Fatalf("expected 3 users, got %d", len(users))
	}
}

func TestReviewService(t *testing.T) {
	reviewRepo := memory.NewReviewRepository()
	svc := service.NewReviewService(reviewRepo)

	ctx := context.Background()

	// No reviews initially
	reviews, err := svc.ByBook(ctx, 1)
	if err != nil {
		t.Fatalf("ByBook: %v", err)
	}
	if len(reviews) != 0 {
		t.Fatalf("expected 0 reviews, got %d", len(reviews))
	}

	// Create a review
	err = svc.Create(ctx, domain.Review{BookID: 1, UserID: 1, Text: "Excellent!", Rating: 5})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	// Should now have 1 review for book 1
	reviews, err = svc.ByBook(ctx, 1)
	if err != nil {
		t.Fatalf("ByBook after create: %v", err)
	}
	if len(reviews) != 1 {
		t.Fatalf("expected 1 review, got %d", len(reviews))
	}
	if reviews[0].Text != "Excellent!" {
		t.Fatalf("unexpected review text: %s", reviews[0].Text)
	}
	if reviews[0].Rating != 5 {
		t.Fatalf("unexpected rating: %d", reviews[0].Rating)
	}
}
