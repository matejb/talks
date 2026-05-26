package book

import (
	"context"
	"fmt"
)

// UserLookup is a function that verifies a user exists.
// This avoids importing the user package directly.
type UserLookup func(ctx context.Context, userID int) error

// Service provides book business logic.
type Service struct {
	repo   Repository
	lookup UserLookup
}

// NewService creates a new book Service.
func NewService(repo Repository, lookup UserLookup) *Service {
	return &Service{repo: repo, lookup: lookup}
}

// Book returns a single book by ID.
func (s *Service) Book(_ context.Context, id int) (Book, error) {
	return s.repo.ByID(id)
}

// Books returns all books.
func (s *Service) Books(_ context.Context) ([]Book, error) {
	return s.repo.All()
}

// Borrow marks a book as borrowed by the given user.
func (s *Service) Borrow(ctx context.Context, bookID, userID int) error {
	if err := s.lookup(ctx, userID); err != nil {
		return fmt.Errorf("user lookup: %w", err)
	}

	b, err := s.repo.ByID(bookID)
	if err != nil {
		return fmt.Errorf("find book: %w", err)
	}
	if b.BorrowerID != 0 {
		return fmt.Errorf("book %d already borrowed", bookID)
	}

	b.BorrowerID = userID
	return s.repo.Save(b)
}

// Return marks a book as returned.
func (s *Service) Return(_ context.Context, bookID int) error {
	b, err := s.repo.ByID(bookID)
	if err != nil {
		return fmt.Errorf("find book: %w", err)
	}
	if b.BorrowerID == 0 {
		return fmt.Errorf("book %d is not borrowed", bookID)
	}
	b.BorrowerID = 0
	return s.repo.Save(b)
}
