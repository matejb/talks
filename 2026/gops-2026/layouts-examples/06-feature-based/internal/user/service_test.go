package user

import (
	"context"
	"testing"
)

func TestService_User(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)

	// Get pre-seeded user
	u, err := svc.User(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if u.Name != "Alice" {
		t.Fatalf("expected Alice, got %s", u.Name)
	}

	// Non-existent user
	_, err = svc.User(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestService_Users(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)

	users, err := svc.Users(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}

func TestService_CreateUser(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)

	err := svc.CreateUser(context.Background(), User{Name: "Charlie", Email: "charlie@example.com"})
	if err != nil {
		t.Fatal(err)
	}

	users, _ := svc.Users(context.Background())
	if len(users) != 3 {
		t.Fatalf("expected 3 users after create, got %d", len(users))
	}

	// Verify the new user has ID 3
	u, err := svc.User(context.Background(), 3)
	if err != nil {
		t.Fatal(err)
	}
	if u.Name != "Charlie" {
		t.Fatalf("expected Charlie, got %s", u.Name)
	}
}
