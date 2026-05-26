package library

import "testing"

func TestStore_CreateAndGetUser(t *testing.T) {
	s := NewStore()

	users, err := s.Users()
	if err != nil {
		t.Fatal(err)
	}
	if len(users) < 2 {
		t.Fatalf("expected at least 2 seeded users, got %d", len(users))
	}

	u, err := s.CreateUser(User{Name: "Charlie", Email: "charlie@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if u.ID != 3 {
		t.Fatalf("expected ID 3, got %d", u.ID)
	}

	got, err := s.User(u.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "Charlie" {
		t.Fatalf("expected Charlie, got %s", got.Name)
	}

	_, err = s.User(999)
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestStore_BorrowAndReturn(t *testing.T) {
	s := NewStore()

	err := s.BorrowBook(1, 1)
	if err != nil {
		t.Fatal(err)
	}

	b, err := s.Book(1)
	if err != nil {
		t.Fatal(err)
	}
	if b.Borrowed == nil || b.Borrowed.Name != "Alice" {
		t.Fatal("expected book to be borrowed by Alice")
	}

	// Can't borrow again
	err = s.BorrowBook(1, 2)
	if err == nil {
		t.Fatal("expected error when borrowing already borrowed book")
	}

	// Borrowed books list
	borrowed, err := s.BorrowedBooks()
	if err != nil {
		t.Fatal(err)
	}
	if len(borrowed) != 1 {
		t.Fatalf("expected 1 borrowed book, got %d", len(borrowed))
	}

	// Return
	err = s.ReturnBook(1)
	if err != nil {
		t.Fatal(err)
	}

	b, err = s.Book(1)
	if err != nil {
		t.Fatal(err)
	}
	if b.Borrowed != nil {
		t.Fatal("expected book to be returned")
	}

	// Return non-existent
	err = s.ReturnBook(999)
	if err == nil {
		t.Fatal("expected error returning non-existent book")
	}
}

func TestStore_Reviews(t *testing.T) {
	s := NewStore()

	// No reviews initially
	reviews, err := s.Reviews(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(reviews) != 0 {
		t.Fatalf("expected 0 reviews, got %d", len(reviews))
	}

	r, err := s.CreateReview(Review{BookID: 1, UserID: 1, Text: "Great book!", Rating: 5})
	if err != nil {
		t.Fatal(err)
	}
	if r.ID != 1 {
		t.Fatalf("expected review ID 1, got %d", r.ID)
	}

	reviews, err = s.Reviews(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(reviews) != 1 || reviews[0].Text != "Great book!" {
		t.Fatalf("unexpected reviews: %+v", reviews)
	}
}

func TestStore_Books(t *testing.T) {
	s := NewStore()

	books, err := s.Books()
	if err != nil {
		t.Fatal(err)
	}
	if len(books) != 2 {
		t.Fatalf("expected 2 seeded books, got %d", len(books))
	}

	b, err := s.CreateBook(Book{Title: "Go in Action", Author: "Kennedy, Ketelsen & Martin"})
	if err != nil {
		t.Fatal(err)
	}
	if b.ID != 3 {
		t.Fatalf("expected ID 3, got %d", b.ID)
	}

	_, err = s.Book(999)
	if err == nil {
		t.Fatal("expected error for non-existent book")
	}
}
