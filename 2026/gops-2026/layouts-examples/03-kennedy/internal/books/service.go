package books

import (
	"errors"
	"fmt"
	"sync"
)

// ErrUserNotFound is returned when a user lookup fails.
var ErrUserNotFound = errors.New("user not found")

// UserLookup is a function the books service uses to resolve user names.
// In the Kennedy layout, cross-domain references are wired at the
// composition root rather than through shared interfaces.
type UserLookup func(id int) (string, error)

// Service provides book business logic.
type Service struct {
	mu     sync.RWMutex
	books  map[int]Book
	nextID int
	lookup UserLookup
}

// NewService creates a book service with pre-seeded data.
func NewService(lookup UserLookup) *Service {
	return &Service{
		books: map[int]Book{
			1: {ID: 1, Title: "The Go Programming Language", Author: "Donovan & Kernighan"},
			2: {ID: 2, Title: "Concurrency in Go", Author: "Katherine Cox-Buday"},
		},
		nextID: 3,
		lookup: lookup,
	}
}

// Book returns a book by ID.
func (s *Service) Book(id int) (Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.books[id]
	if !ok {
		return Book{}, fmt.Errorf("book %d not found", id)
	}
	return b, nil
}

// Books returns all books.
func (s *Service) Books() ([]Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Book, 0, len(s.books))
	for _, b := range s.books {
		result = append(result, b)
	}
	return result, nil
}

// Borrow marks a book as borrowed by the given user.
func (s *Service) Borrow(bookID, userID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.books[bookID]
	if !ok {
		return fmt.Errorf("book %d not found", bookID)
	}
	if b.BorrowerID != 0 {
		return fmt.Errorf("book %d already borrowed", bookID)
	}
	name, err := s.lookup(userID)
	if err != nil {
		return fmt.Errorf("user %d not found", userID)
	}
	b.BorrowerID = userID
	b.BorrowerName = name
	s.books[bookID] = b
	return nil
}

// Return marks a book as returned.
func (s *Service) Return(bookID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.books[bookID]
	if !ok {
		return fmt.Errorf("book %d not found", bookID)
	}
	b.BorrowerID = 0
	b.BorrowerName = ""
	s.books[bookID] = b
	return nil
}
