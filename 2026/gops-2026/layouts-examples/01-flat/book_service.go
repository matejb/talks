package main

import (
	"fmt"
	"sync"
)

type bookService struct {
	mu     sync.RWMutex
	books  map[int]Book
	nextID int
	users  *userService // flat layout: direct reference to another service
}

func newBookService(us *userService) *bookService {
	return &bookService{
		books: map[int]Book{
			1: {ID: 1, Title: "The Go Programming Language", Author: "Donovan & Kernighan"},
			2: {ID: 2, Title: "Concurrency in Go", Author: "Katherine Cox-Buday"},
		},
		nextID: 3,
		users:  us,
	}
}

func (s *bookService) Book(id int) (Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.books[id]
	if !ok {
		return Book{}, fmt.Errorf("book %d not found", id)
	}
	return b, nil
}

func (s *bookService) Books() ([]Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Book, 0, len(s.books))
	for _, b := range s.books {
		result = append(result, b)
	}
	return result, nil
}

func (s *bookService) Borrow(bookID, userID int) error {
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
	b.Borrowed = &User{ID: u.ID, Name: u.Name, Email: u.Email}
	s.books[bookID] = b
	return nil
}

func (s *bookService) Return(bookID int) error {
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
