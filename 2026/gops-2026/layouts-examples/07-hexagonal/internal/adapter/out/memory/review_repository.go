package memory

import (
	"context"
	"sync"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/core/domain"
	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/core/port"
)

// Compile-time check: ReviewRepository implements port.ReviewRepository.
var _ port.ReviewRepository = (*ReviewRepository)(nil)

// ReviewRepository is an in-memory implementation of port.ReviewRepository.
type ReviewRepository struct {
	mu      sync.RWMutex
	reviews map[int]domain.Review
	nextID  int
}

// NewReviewRepository creates a ReviewRepository with no pre-seeded data.
func NewReviewRepository() *ReviewRepository {
	return &ReviewRepository{
		reviews: make(map[int]domain.Review),
		nextID:  1,
	}
}

func (r *ReviewRepository) ByBookID(_ context.Context, bookID int) ([]domain.Review, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Review
	for _, rev := range r.reviews {
		if rev.BookID == bookID {
			result = append(result, rev)
		}
	}
	return result, nil
}

func (r *ReviewRepository) Save(_ context.Context, rev domain.Review) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	rev.ID = r.nextID
	r.nextID++
	r.reviews[rev.ID] = rev
	return nil
}
