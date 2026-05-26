package main

import (
	"fmt"
	"sync"
)

type reviewService struct {
	mu      sync.RWMutex
	reviews map[int]Review
	nextID  int
}

func newReviewService() *reviewService {
	return &reviewService{
		reviews: make(map[int]Review),
		nextID:  1,
	}
}

func (s *reviewService) Reviews(bookID int) ([]Review, error) {
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

func (s *reviewService) CreateReview(review Review) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	review.ID = s.nextID
	s.nextID++
	s.reviews[review.ID] = review
	return nil
}

func (s *reviewService) Review(id int) (Review, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.reviews[id]
	if !ok {
		return Review{}, fmt.Errorf("review %d not found", id)
	}
	return r, nil
}
