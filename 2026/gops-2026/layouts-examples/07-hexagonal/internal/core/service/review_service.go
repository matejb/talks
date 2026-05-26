package service

import (
	"context"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/core/domain"
	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/core/port"
)

// ReviewService implements the review-related use cases.
type ReviewService struct {
	reviews port.ReviewRepository
}

// NewReviewService creates a new ReviewService with the given repository.
func NewReviewService(reviews port.ReviewRepository) *ReviewService {
	return &ReviewService{reviews: reviews}
}

// ByBook returns all reviews for a given book.
func (s *ReviewService) ByBook(ctx context.Context, bookID int) ([]domain.Review, error) {
	return s.reviews.ByBookID(ctx, bookID)
}

// Create adds a new review.
func (s *ReviewService) Create(ctx context.Context, r domain.Review) error {
	return s.reviews.Save(ctx, r)
}
