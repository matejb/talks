package memory

import (
	"fmt"
	"sync"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/04-bourgon-pkg/pkg/library"
)

// BookService is an in-memory implementation of library.BookService.
type BookService struct {
	mu     sync.RWMutex
	books  map[int]library.Book
	nextID int
	users  library.UserService
}

// NewBookService returns a BookService pre-seeded with sample data.
func NewBookService(users library.UserService) *BookService {
	return &BookService{
		books: map[int]library.Book{
			1: {ID: 1, Title: "The Go Programming Language", Author: "Donovan & Kernighan"},
			2: {ID: 2, Title: "Concurrency in Go", Author: "Katherine Cox-Buday"},
		},
		nextID: 3,
		users:  users,
	}
}

func (s *BookService) Book(id int) (library.Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.books[id]
	if !ok {
		return library.Book{}, fmt.Errorf("book %d not found", id)
	}
	return b, nil
}

func (s *BookService) Books() ([]library.Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]library.Book, 0, len(s.books))
	for _, b := range s.books {
		result = append(result, b)
	}
	return result, nil
}

func (s *BookService) Borrow(bookID, userID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.books[bookID]
	if !ok {
		return fmt.Errorf("book %d not found", bookID)
	}
	if b.Borrowed != nil {
		return fmt.Errorf("book %d already borrowed by %s", bookID, b.Borrowed.Name)
	}
	u, err := s.users.User(userID)
	if err != nil {
		return fmt.Errorf("user %d not found", userID)
	}
	b.Borrowed = &library.User{ID: u.ID, Name: u.Name, Email: u.Email}
	s.books[bookID] = b
	return nil
}

func (s *BookService) Return(bookID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.books[bookID]
	if !ok {
		return fmt.Errorf("book %d not found", bookID)
	}
	b.Borrowed = nil
	s.books[bookID] = b
	return nil
}
