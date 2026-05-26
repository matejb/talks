package review

import (
	"sync"
)

// Repository defines the data access interface for reviews.
type Repository interface {
	ByBookID(bookID int) ([]Review, error)
	Create(r Review) error
}

// In-memory implementation of Repository.
type memoryRepository struct {
	mu      sync.RWMutex
	reviews map[int]Review
	nextID  int
}

// NewMemoryRepository returns a Repository backed by memory.
func NewMemoryRepository() Repository {
	return &memoryRepository{
		reviews: make(map[int]Review),
		nextID:  1,
	}
}

func (r *memoryRepository) ByBookID(bookID int) ([]Review, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []Review
	for _, rv := range r.reviews {
		if rv.BookID == bookID {
			result = append(result, rv)
		}
	}
	return result, nil
}

func (r *memoryRepository) Create(rv Review) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	rv.ID = r.nextID
	r.nextID++
	r.reviews[rv.ID] = rv
	return nil
}

// Compile-time check.
var _ Repository = (*memoryRepository)(nil)
