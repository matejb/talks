package book

import (
	"context"
	"fmt"
	"testing"
)

// stubUserLookup returns a lookup that succeeds for user IDs 1 and 2.
func stubUserLookup() UserLookup {
	return func(_ context.Context, userID int) error {
		if userID != 1 && userID != 2 {
			return errNotFound
		}
		return nil
	}
}

var errNotFound = fmt.Errorf("not found")

func TestService_Book(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo, stubUserLookup())

	// Get pre-seeded book
	b, err := svc.Book(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if b.Title != "The Go Programming Language" {
		t.Fatalf("unexpected title: %s", b.Title)
	}

	// Non-existent book
	_, err = svc.Book(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error for non-existent book")
	}
}

func TestService_Books(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo, stubUserLookup())

	books, err := svc.Books(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(books) != 2 {
		t.Fatalf("expected 2 books, got %d", len(books))
	}
}

func TestService_Borrow(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo, stubUserLookup())

	// Borrow successfully
	err := svc.Borrow(context.Background(), 1, 1)
	if err != nil {
		t.Fatal(err)
	}

	// Verify book is borrowed
	b, _ := svc.Book(context.Background(), 1)
	if b.BorrowerID != 1 {
		t.Fatalf("expected borrower ID 1, got %d", b.BorrowerID)
	}

	// Cannot borrow again
	err = svc.Borrow(context.Background(), 1, 2)
	if err == nil {
		t.Fatal("expected error when borrowing already borrowed book")
	}

	// Non-existent user
	err = svc.Borrow(context.Background(), 2, 999)
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestService_Return(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo, stubUserLookup())

	// Borrow first
	svc.Borrow(context.Background(), 1, 1)

	// Return successfully
	err := svc.Return(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}

	// Verify returned
	b, _ := svc.Book(context.Background(), 1)
	if b.BorrowerID != 0 {
		t.Fatalf("expected borrower ID 0, got %d", b.BorrowerID)
	}

	// Cannot return a not-borrowed book
	err = svc.Return(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error when returning non-borrowed book")
	}
}
