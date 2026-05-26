package library

import "testing"

func TestUserService_Create(t *testing.T) {
	store := NewMemoryUserStore()
	svc := &UserService{Store: store}

	users, _ := svc.Users()
	if len(users) != 2 {
		t.Fatalf("expected 2 pre-seeded users, got %d", len(users))
	}

	err := svc.CreateUser(User{Name: "Charlie", Email: "charlie@example.com"})
	if err != nil {
		t.Fatal(err)
	}

	u, err := svc.User(3)
	if err != nil {
		t.Fatal(err)
	}
	if u.Name != "Charlie" {
		t.Fatalf("expected Charlie, got %s", u.Name)
	}

	_, err = svc.User(999)
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestUserService_Validation(t *testing.T) {
	store := NewMemoryUserStore()
	svc := &UserService{Store: store}

	err := svc.CreateUser(User{Name: "", Email: "test@example.com"})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestBookService_BorrowAndReturn(t *testing.T) {
	userStore := NewMemoryUserStore()
	bookStore := NewMemoryBookStore()
	svc := &BookService{Store: bookStore, Users: userStore}

	books, _ := svc.Books()
	if len(books) != 2 {
		t.Fatalf("expected 2 pre-seeded books, got %d", len(books))
	}

	// Borrow
	err := svc.Borrow(1, 1)
	if err != nil {
		t.Fatal(err)
	}

	b, _ := svc.Book(1)
	if b.Borrowed == nil || b.Borrowed.Name != "Alice" {
		t.Fatal("expected book to be borrowed by Alice")
	}

	// Can't borrow again
	err = svc.Borrow(1, 2)
	if err == nil {
		t.Fatal("expected error for already borrowed book")
	}

	// Return
	err = svc.Return(1)
	if err != nil {
		t.Fatal(err)
	}

	b, _ = svc.Book(1)
	if b.Borrowed != nil {
		t.Fatal("expected book to be returned")
	}

	// Borrow non-existent book
	err = svc.Borrow(999, 1)
	if err == nil {
		t.Fatal("expected error for non-existent book")
	}

	// Borrow with non-existent user
	err = svc.Borrow(1, 999)
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestReviewService_Create(t *testing.T) {
	store := NewMemoryReviewStore()
	svc := &ReviewService{Store: store}

	err := svc.CreateReview(Review{BookID: 1, UserID: 1, Text: "Great!", Rating: 5})
	if err != nil {
		t.Fatal(err)
	}

	reviews, _ := svc.Reviews(1)
	if len(reviews) != 1 {
		t.Fatalf("expected 1 review, got %d", len(reviews))
	}
	if reviews[0].Text != "Great!" {
		t.Fatalf("expected 'Great!', got '%s'", reviews[0].Text)
	}

	// Validation: missing book ID
	err = svc.CreateReview(Review{UserID: 1, Text: "Ok", Rating: 3})
	if err == nil {
		t.Fatal("expected error for missing book ID")
	}

	// Validation: invalid rating
	err = svc.CreateReview(Review{BookID: 1, UserID: 1, Text: "Bad", Rating: 0})
	if err == nil {
		t.Fatal("expected error for rating below 1")
	}

	err = svc.CreateReview(Review{BookID: 1, UserID: 1, Text: "Perfect", Rating: 6})
	if err == nil {
		t.Fatal("expected error for rating above 5")
	}
}
