package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/core/domain"
	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/core/port"
)

// Compile-time check: BookRepository implements port.BookRepository.
var _ port.BookRepository = (*BookRepository)(nil)

// BookRepository is an in-memory implementation of port.BookRepository.
type BookRepository struct {
	mu     sync.RWMutex
	books  map[int]domain.Book
	nextID int
}

// NewBookRepository creates a BookRepository with pre-seeded data.
func NewBookRepository() *BookRepository {
	return &BookRepository{
		books: map[int]domain.Book{
			1: {ID: 1, Title: "The Go Programming Language", Author: "Donovan & Kernighan"},
			2: {ID: 2, Title: "Concurrency in Go", Author: "Katherine Cox-Buday"},
		},
		nextID: 3,
	}
}

func (r *BookRepository) ByID(_ context.Context, id int) (domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.books[id]
	if !ok {
		return domain.Book{}, fmt.Errorf("book %d not found", id)
	}
	return b, nil
}

func (r *BookRepository) All(_ context.Context) ([]domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domain.Book, 0, len(r.books))
	for _, b := range r.books {
		result = append(result, b)
	}
	return result, nil
}

func (r *BookRepository) Save(_ context.Context, b domain.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	b.ID = r.nextID
	r.nextID++
	r.books[b.ID] = b
	return nil
}

func (r *BookRepository) Update(_ context.Context, b domain.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.books[b.ID]; !ok {
		return fmt.Errorf("book %d not found", b.ID)
	}
	r.books[b.ID] = b
	return nil
}
