package review

import (
	"context"
	"fmt"
	"testing"
)

// stubBookLookup succeeds for book IDs 1 and 2.
func stubBookLookup() BookLookup {
	return func(_ context.Context, bookID int) error {
		if bookID != 1 && bookID != 2 {
			return errBookNotFound
		}
		return nil
	}
}

var errBookNotFound = fmt.Errorf("book not found")

func TestService_Reviews(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo, stubBookLookup())

	// No reviews initially
	reviews, err := svc.Reviews(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(reviews) != 0 {
		t.Fatalf("expected 0 reviews, got %d", len(reviews))
	}

	// Create a review
	err = svc.CreateReview(context.Background(), Review{
		BookID: 1,
		UserID: 1,
		Text:   "Great book!",
		Rating: 5,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Now there is one
	reviews, _ = svc.Reviews(context.Background(), 1)
	if len(reviews) != 1 {
		t.Fatalf("expected 1 review, got %d", len(reviews))
	}
	if reviews[0].Text != "Great book!" {
		t.Fatalf("expected 'Great book!', got '%s'", reviews[0].Text)
	}
}

func TestService_CreateReview_InvalidBook(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo, stubBookLookup())

	err := svc.CreateReview(context.Background(), Review{
		BookID: 999,
		UserID: 1,
		Text:   "Nope",
		Rating: 1,
	})
	if err == nil {
		t.Fatal("expected error for non-existent book")
	}
}
