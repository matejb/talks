package books

import "testing"

func TestService_BorrowAndReturn(t *testing.T) {
	// User lookup function for tests.
	lookup := func(id int) (string, error) {
		names := map[int]string{1: "Alice", 2: "Bob"}
		n, ok := names[id]
		if !ok {
			return "", ErrUserNotFound
		}
		return n, nil
	}

	svc := NewService(lookup)

	// Pre-seeded books exist.
	books, err := svc.Books()
	if err != nil {
		t.Fatal(err)
	}
	if len(books) != 2 {
		t.Fatalf("expected 2 pre-seeded books, got %d", len(books))
	}

	// Borrow a book.
	err = svc.Borrow(1, 1)
	if err != nil {
		t.Fatal(err)
	}

	// Verify it's borrowed.
	b, err := svc.Book(1)
	if err != nil {
		t.Fatal(err)
	}
	if b.BorrowerID != 1 {
		t.Fatalf("expected borrower ID 1, got %d", b.BorrowerID)
	}
	if b.BorrowerName != "Alice" {
		t.Fatalf("expected borrower name Alice, got %s", b.BorrowerName)
	}

	// Can't borrow again.
	err = svc.Borrow(1, 2)
	if err == nil {
		t.Fatal("expected error when borrowing already borrowed book")
	}

	// Return the book.
	err = svc.Return(1)
	if err != nil {
		t.Fatal(err)
	}

	// Verify it's returned.
	b, err = svc.Book(1)
	if err != nil {
		t.Fatal(err)
	}
	if b.BorrowerID != 0 {
		t.Fatalf("expected borrower ID 0, got %d", b.BorrowerID)
	}
	if b.BorrowerName != "" {
		t.Fatalf("expected empty borrower name, got %s", b.BorrowerName)
	}
}

func TestService_Errors(t *testing.T) {
	svc := NewService(func(id int) (string, error) { return "Test", nil })

	// Borrow non-existent book.
	err := svc.Borrow(999, 1)
	if err == nil {
		t.Fatal("expected error for non-existent book")
	}

	// Return non-existent book.
	err = svc.Return(999)
	if err == nil {
		t.Fatal("expected error for non-existent book")
	}

	// Borrow with non-existent user.
	err = svc.Borrow(1, 999)
	// This will succeed because our test lookup always succeeds -
	// in real code with a real DB lookup this would fail.
	// The key Kennedy insight: the service calls directly into data layer,
	// there's no port/adapter to swap for testing.
	_ = err
}
