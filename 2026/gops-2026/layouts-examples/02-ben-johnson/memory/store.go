package memory

import (
	"fmt"
	"sync"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/02-ben-johnson"
)

// Compile-time interface checks.
var (
	_ library.UserService   = (*UserService)(nil)
	_ library.BookService   = (*BookService)(nil)
	_ library.ReviewService = (*ReviewService)(nil)
)

// --- UserService ---

// UserService is an in-memory implementation of library.UserService.
type UserService struct {
	mu     sync.RWMutex
	users  map[int]library.User
	nextID int
}

// NewUserService creates a UserService pre-seeded with sample data.
func NewUserService() *UserService {
	return &UserService{
		users: map[int]library.User{
			1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
			2: {ID: 2, Name: "Bob", Email: "bob@example.com"},
		},
		nextID: 3,
	}
}

func (s *UserService) User(id int) (library.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return library.User{}, fmt.Errorf("user %d not found", id)
	}
	return u, nil
}

func (s *UserService) Users() ([]library.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]library.User, 0, len(s.users))
	for _, u := range s.users {
		result = append(result, u)
	}
	return result, nil
}

func (s *UserService) CreateUser(user library.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	user.ID = s.nextID
	s.nextID++
	s.users[user.ID] = user
	return nil
}

// --- BookService ---

// BookService is an in-memory implementation of library.BookService.
type BookService struct {
	mu     sync.RWMutex
	books  map[int]library.Book
	nextID int
	users  *UserService
}

// NewBookService creates a BookService pre-seeded with sample data.
func NewBookService(us *UserService) *BookService {
	return &BookService{
		books: map[int]library.Book{
			1: {ID: 1, Title: "The Go Programming Language", Author: "Donovan & Kernighan"},
			2: {ID: 2, Title: "Concurrency in Go", Author: "Katherine Cox-Buday"},
		},
		nextID: 3,
		users:  us,
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

// --- ReviewService ---

// ReviewService is an in-memory implementation of library.ReviewService.
type ReviewService struct {
	mu      sync.RWMutex
	reviews map[int]library.Review
	nextID  int
}

// NewReviewService creates a ReviewService with empty state.
func NewReviewService() *ReviewService {
	return &ReviewService{
		reviews: make(map[int]library.Review),
		nextID:  1,
	}
}

func (s *ReviewService) Reviews(bookID int) ([]library.Review, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []library.Review
	for _, r := range s.reviews {
		if r.BookID == bookID {
			result = append(result, r)
		}
	}
	return result, nil
}

func (s *ReviewService) CreateReview(review library.Review) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	review.ID = s.nextID
	s.nextID++
	s.reviews[review.ID] = review
	return nil
}
