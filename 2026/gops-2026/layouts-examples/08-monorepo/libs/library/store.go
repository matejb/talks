package library

import (
	"fmt"
	"sync"
)

// Store provides an in-memory store for library entities.
// In a real application this would be replaced by a database.
type Store struct {
	mu         sync.RWMutex
	users      map[int]User
	books      map[int]Book
	reviews    map[int]Review
	nextUser   int
	nextBook   int
	nextReview int
}

// NewStore creates a Store pre-seeded with sample data.
func NewStore() *Store {
	s := &Store{
		users:      make(map[int]User),
		books:      make(map[int]Book),
		reviews:    make(map[int]Review),
		nextUser:   1,
		nextBook:   1,
		nextReview: 1,
	}

	// Seed users
	s.CreateUser(User{Name: "Alice", Email: "alice@example.com"})
	s.CreateUser(User{Name: "Bob", Email: "bob@example.com"})

	// Seed books
	s.CreateBook(Book{Title: "The Go Programming Language", Author: "Donovan & Kernighan"})
	s.CreateBook(Book{Title: "Concurrency in Go", Author: "Katherine Cox-Buday"})

	return s
}

// --- Users ---

func (s *Store) User(id int) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return User{}, fmt.Errorf("user %d not found", id)
	}
	return u, nil
}

func (s *Store) Users() ([]User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]User, 0, len(s.users))
	for _, u := range s.users {
		result = append(result, u)
	}
	return result, nil
}

func (s *Store) CreateUser(u User) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u.ID = s.nextUser
	s.nextUser++
	s.users[u.ID] = u
	return u, nil
}

// --- Books ---

func (s *Store) Book(id int) (Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.books[id]
	if !ok {
		return Book{}, fmt.Errorf("book %d not found", id)
	}
	return b, nil
}

func (s *Store) Books() ([]Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Book, 0, len(s.books))
	for _, b := range s.books {
		result = append(result, b)
	}
	return result, nil
}

func (s *Store) CreateBook(b Book) (Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b.ID = s.nextBook
	s.nextBook++
	s.books[b.ID] = b
	return b, nil
}

func (s *Store) BorrowBook(bookID, userID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.books[bookID]
	if !ok {
		return fmt.Errorf("book %d not found", bookID)
	}
	if b.Borrowed != nil {
		return fmt.Errorf("book %d already borrowed by %s", bookID, b.Borrowed.Name)
	}
	u, ok := s.users[userID]
	if !ok {
		return fmt.Errorf("user %d not found", userID)
	}
	b.Borrowed = &User{ID: u.ID, Name: u.Name, Email: u.Email}
	s.books[bookID] = b
	return nil
}

func (s *Store) ReturnBook(bookID int) error {
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

// BorrowedBooks returns all books that are currently borrowed.
func (s *Store) BorrowedBooks() ([]Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []Book
	for _, b := range s.books {
		if b.Borrowed != nil {
			result = append(result, b)
		}
	}
	return result, nil
}

// --- Reviews ---

func (s *Store) Reviews(bookID int) ([]Review, error) {
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

func (s *Store) CreateReview(r Review) (Review, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r.ID = s.nextReview
	s.nextReview++
	s.reviews[r.ID] = r
	return r, nil
}
