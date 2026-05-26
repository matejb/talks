package memory

import (
	"sync"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/04-bourgon-pkg/pkg/library"
)

// ReviewService is an in-memory implementation of library.ReviewService.
type ReviewService struct {
	mu      sync.RWMutex
	reviews map[int]library.Review
	nextID  int
}

// NewReviewService returns an empty ReviewService.
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
