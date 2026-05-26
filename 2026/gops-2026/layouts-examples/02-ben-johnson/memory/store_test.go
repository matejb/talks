package memory

import (
	"testing"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/02-ben-johnson"
)

func TestUserService(t *testing.T) {
	svc := NewUserService()

	// Pre-seeded data.
	users, err := svc.Users()
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 pre-seeded users, got %d", len(users))
	}

	// Get existing user.
	u, err := svc.User(1)
	if err != nil {
		t.Fatal(err)
	}
	if u.Name != "Alice" {
		t.Fatalf("expected Alice, got %s", u.Name)
	}

	// Get non-existent user.
	_, err = svc.User(999)
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}

	// Create a new user.
	err = svc.CreateUser(library.User{Name: "Charlie", Email: "charlie@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	u, err = svc.User(3)
	if err != nil {
		t.Fatal(err)
	}
	if u.Name != "Charlie" {
		t.Fatalf("expected Charlie, got %s", u.Name)
	}
}

func TestBookService_BorrowAndReturn(t *testing.T) {
	us := NewUserService()
	bs := NewBookService(us)

	// Pre-seeded data.
	books, err := bs.Books()
	if err != nil {
		t.Fatal(err)
	}
	if len(books) != 2 {
		t.Fatalf("expected 2 pre-seeded books, got %d", len(books))
	}

	// Get existing book.
	b, err := bs.Book(1)
	if err != nil {
		t.Fatal(err)
	}
	if b.Title != "The Go Programming Language" {
		t.Fatalf("unexpected title: %s", b.Title)
	}

	// Borrow a book.
	err = bs.Borrow(1, 1)
	if err != nil {
		t.Fatal(err)
	}

	// Verify it's borrowed.
	b, err = bs.Book(1)
	if err != nil {
		t.Fatal(err)
	}
	if b.Borrowed == nil || b.Borrowed.Name != "Alice" {
		t.Fatal("expected book to be borrowed by Alice")
	}

	// Cannot borrow again.
	err = bs.Borrow(1, 2)
	if err == nil {
		t.Fatal("expected error when borrowing already borrowed book")
	}

	// Return the book.
	err = bs.Return(1)
	if err != nil {
		t.Fatal(err)
	}

	// Verify it's free.
	b, err = bs.Book(1)
	if err != nil {
		t.Fatal(err)
	}
	if b.Borrowed != nil {
		t.Fatal("expected book to be returned (Borrowed == nil)")
	}

	// Borrow with non-existent user.
	err = bs.Borrow(1, 999)
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}

	// Borrow non-existent book.
	err = bs.Borrow(999, 1)
	if err == nil {
		t.Fatal("expected error for non-existent book")
	}
}

func TestReviewService(t *testing.T) {
	svc := NewReviewService()

	// Empty initially.
	reviews, err := svc.Reviews(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(reviews) != 0 {
		t.Fatalf("expected 0 reviews, got %d", len(reviews))
	}

	// Create a review.
	err = svc.CreateReview(library.Review{BookID: 1, UserID: 1, Text: "Great book!", Rating: 5})
	if err != nil {
		t.Fatal(err)
	}

	// Fetch reviews for book.
	reviews, err = svc.Reviews(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(reviews) != 1 {
		t.Fatalf("expected 1 review, got %d", len(reviews))
	}
	if reviews[0].Text != "Great book!" {
		t.Fatalf("unexpected review text: %s", reviews[0].Text)
	}
	if reviews[0].ID != 1 {
		t.Fatalf("expected review ID 1, got %d", reviews[0].ID)
	}
}
