package book

import (
	"fmt"
	"sync"
)

// Repository defines the data access interface for books.
type Repository interface {
	ByID(id int) (Book, error)
	All() ([]Book, error)
	Save(b Book) error
}

// In-memory implementation of Repository.
type memoryRepository struct {
	mu     sync.RWMutex
	books  map[int]Book
	nextID int
}

// NewMemoryRepository returns a Repository backed by memory, pre-seeded with sample data.
func NewMemoryRepository() Repository {
	return &memoryRepository{
		books: map[int]Book{
			1: {ID: 1, Title: "The Go Programming Language", Author: "Donovan & Kernighan"},
			2: {ID: 2, Title: "Concurrency in Go", Author: "Katherine Cox-Buday"},
		},
		nextID: 3,
	}
}

func (r *memoryRepository) ByID(id int) (Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.books[id]
	if !ok {
		return Book{}, fmt.Errorf("book %d not found", id)
	}
	return b, nil
}

func (r *memoryRepository) All() ([]Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]Book, 0, len(r.books))
	for _, b := range r.books {
		result = append(result, b)
	}
	return result, nil
}

func (r *memoryRepository) Save(b Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.books[b.ID] = b
	return nil
}
