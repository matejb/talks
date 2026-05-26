package memory

import (
	"testing"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/04-bourgon-pkg/pkg/library"
)

func TestUserService_CreateAndGet(t *testing.T) {
	svc := NewUserService()

	users, err := svc.Users()
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 pre-seeded users, got %d", len(users))
	}

	err = svc.CreateUser(library.User{Name: "Charlie", Email: "charlie@example.com"})
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

func TestBookService_BorrowAndReturn(t *testing.T) {
	us := NewUserService()
	bs := NewBookService(us)

	books, err := bs.Books()
	if err != nil {
		t.Fatal(err)
	}
	if len(books) != 2 {
		t.Fatalf("expected 2 pre-seeded books, got %d", len(books))
	}

	err = bs.Borrow(1, 1)
	if err != nil {
		t.Fatal(err)
	}

	b, err := bs.Book(1)
	if err != nil {
		t.Fatal(err)
	}
	if b.Borrowed == nil || b.Borrowed.Name != "Alice" {
		t.Fatal("expected book to be borrowed by Alice")
	}

	err = bs.Borrow(1, 2)
	if err == nil {
		t.Fatal("expected error when borrowing already borrowed book")
	}

	err = bs.Return(1)
	if err != nil {
		t.Fatal(err)
	}

	b, err = bs.Book(1)
	if err != nil {
		t.Fatal(err)
	}
	if b.Borrowed != nil {
		t.Fatal("expected book to be returned")
	}
}

func TestBookService_BorrowInvalidUser(t *testing.T) {
	us := NewUserService()
	bs := NewBookService(us)

	err := bs.Borrow(1, 999)
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestBookService_BorrowInvalidBook(t *testing.T) {
	us := NewUserService()
	bs := NewBookService(us)

	err := bs.Borrow(999, 1)
	if err == nil {
		t.Fatal("expected error for non-existent book")
	}
}

func TestReviewService_Create(t *testing.T) {
	svc := NewReviewService()

	err := svc.CreateReview(library.Review{BookID: 1, UserID: 1, Text: "Great book!", Rating: 5})
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
	if reviews[0].ID != 1 {
		t.Fatalf("expected review ID 1, got %d", reviews[0].ID)
	}
}

func TestReviewService_MultipleReviews(t *testing.T) {
	svc := NewReviewService()

	svc.CreateReview(library.Review{BookID: 1, UserID: 1, Text: "Great!", Rating: 5})
	svc.CreateReview(library.Review{BookID: 1, UserID: 2, Text: "Good", Rating: 4})
	svc.CreateReview(library.Review{BookID: 2, UserID: 1, Text: "Okay", Rating: 3})

	reviews, _ := svc.Reviews(1)
	if len(reviews) != 2 {
		t.Fatalf("expected 2 reviews for book 1, got %d", len(reviews))
	}

	reviews, _ = svc.Reviews(2)
	if len(reviews) != 1 {
		t.Fatalf("expected 1 review for book 2, got %d", len(reviews))
	}
}
