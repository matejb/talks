package users

import "testing"

func TestService_CreateAndGet(t *testing.T) {
	svc := NewService()

	// Pre-seeded users exist.
	users, err := svc.Users()
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 pre-seeded users, got %d", len(users))
	}

	// Create a new user.
	err = svc.CreateUser(User{Name: "Charlie", Email: "charlie@example.com"})
	if err != nil {
		t.Fatal(err)
	}

	// Get the new user (ID 3).
	u, err := svc.User(3)
	if err != nil {
		t.Fatal(err)
	}
	if u.Name != "Charlie" {
		t.Fatalf("expected Charlie, got %s", u.Name)
	}

	// Non-existent user.
	_, err = svc.User(999)
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}
