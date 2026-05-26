package library

import (
	"fmt"
	"sync"
)

// MemoryUserStore is an in-memory implementation of UserStore.
type MemoryUserStore struct {
	mu     sync.RWMutex
	users  map[int]User
	nextID int
}

func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{
		users: map[int]User{
			1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
			2: {ID: 2, Name: "Bob", Email: "bob@example.com"},
		},
		nextID: 3,
	}
}

func (s *MemoryUserStore) GetUser(id int) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return User{}, fmt.Errorf("user %d not found", id)
	}
	return u, nil
}

func (s *MemoryUserStore) ListUsers() ([]User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]User, 0, len(s.users))
	for _, u := range s.users {
		result = append(result, u)
	}
	return result, nil
}

func (s *MemoryUserStore) CreateUser(user User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	user.ID = s.nextID
	s.nextID++
	s.users[user.ID] = user
	return nil
}

// MemoryBookStore is an in-memory implementation of BookStore.
type MemoryBookStore struct {
	mu     sync.RWMutex
	books  map[int]Book
	nextID int
}

func NewMemoryBookStore() *MemoryBookStore {
	return &MemoryBookStore{
		books: map[int]Book{
			1: {ID: 1, Title: "The Go Programming Language", Author: "Donovan & Kernighan"},
			2: {ID: 2, Title: "Concurrency in Go", Author: "Katherine Cox-Buday"},
		},
		nextID: 3,
	}
}

func (s *MemoryBookStore) GetBook(id int) (Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.books[id]
	if !ok {
		return Book{}, fmt.Errorf("book %d not found", id)
	}
	return b, nil
}

func (s *MemoryBookStore) ListBooks() ([]Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Book, 0, len(s.books))
	for _, b := range s.books {
		result = append(result, b)
	}
	return result, nil
}

func (s *MemoryBookStore) SaveBook(book Book) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.books[book.ID] = book
	return nil
}

// MemoryReviewStore is an in-memory implementation of ReviewStore.
type MemoryReviewStore struct {
	mu      sync.RWMutex
	reviews map[int]Review
	nextID  int
}

func NewMemoryReviewStore() *MemoryReviewStore {
	return &MemoryReviewStore{
		reviews: make(map[int]Review),
		nextID:  1,
	}
}

func (s *MemoryReviewStore) ListReviews(bookID int) ([]Review, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []Review
	for _, r := range s.reviews {
		if r.BookID == bookID {
			result = append(result, r)
		}
	}
	return result, nil
}

func (s *MemoryReviewStore) CreateReview(review Review) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	review.ID = s.nextID
	s.nextID++
	s.reviews[review.ID] = review
	return nil
}
