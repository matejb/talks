package main

import (
	"testing"
)

func TestUserService_CreateAndGet(t *testing.T) {
	svc := newUserService()

	// Pre-seeded users exist
	users, err := svc.Users()
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}

	// Create a new user
	err = svc.CreateUser(User{Name: "Charlie", Email: "charlie@example.com"})
	if err != nil {
		t.Fatal(err)
	}

	// Get the new user
	u, err := svc.User(3)
	if err != nil {
		t.Fatal(err)
	}
	if u.Name != "Charlie" {
		t.Fatalf("expected Charlie, got %s", u.Name)
	}

	// Non-existent user returns error
	_, err = svc.User(999)
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestBookService_BorrowAndReturn(t *testing.T) {
	us := newUserService()
	bs := newBookService(us)

	// List books
	books, err := bs.Books()
	if err != nil {
		t.Fatal(err)
	}
	if len(books) != 2 {
		t.Fatalf("expected 2 books, got %d", len(books))
	}

	// Borrow a book
	err = bs.Borrow(1, 1)
	if err != nil {
		t.Fatal(err)
	}

	// Check it's borrowed
	b, err := bs.Book(1)
	if err != nil {
		t.Fatal(err)
	}
	if b.Borrowed == nil || b.Borrowed.Name != "Alice" {
		t.Fatal("expected book to be borrowed by Alice")
	}

	// Can't borrow again
	err = bs.Borrow(1, 2)
	if err == nil {
		t.Fatal("expected error when borrowing already borrowed book")
	}

	// Return the book
	err = bs.Return(1)
	if err != nil {
		t.Fatal(err)
	}

	// Check it's returned
	b, err = bs.Book(1)
	if err != nil {
		t.Fatal(err)
	}
	if b.Borrowed != nil {
		t.Fatal("expected book to be returned")
	}
}

func TestReviewService_Create(t *testing.T) {
	svc := newReviewService()

	err := svc.CreateReview(Review{BookID: 1, UserID: 1, Text: "Great book!", Rating: 5})
	if err != nil {
		t.Fatal(err)
	}

	reviews, err := svc.Reviews(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(reviews) != 1 {
		t.Fatalf("expected 1 review, got %d", len(reviews))
	}
	if reviews[0].Text != "Great book!" {
		t.Fatalf("expected 'Great book!', got '%s'", reviews[0].Text)
	}
}
