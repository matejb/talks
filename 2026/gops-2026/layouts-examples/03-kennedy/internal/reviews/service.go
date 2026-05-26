package reviews

import (
	"sync"
)

// Service provides review business logic.
type Service struct {
	mu      sync.RWMutex
	reviews map[int]Review
	nextID  int
}

// NewService creates a review service.
func NewService() *Service {
	return &Service{
		reviews: make(map[int]Review),
		nextID:  1,
	}
}

// Reviews returns all reviews for a given book.
func (s *Service) Reviews(bookID int) ([]Review, error) {
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

// CreateReview adds a new review.
func (s *Service) CreateReview(r Review) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	r.ID = s.nextID
	s.nextID++
	s.reviews[r.ID] = r
	return nil
}
