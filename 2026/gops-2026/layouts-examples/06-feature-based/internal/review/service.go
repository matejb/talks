package review

import "context"

// BookLookup is a function that verifies a book exists.
type BookLookup func(ctx context.Context, bookID int) error

// Service provides review business logic.
type Service struct {
	repo   Repository
	lookup BookLookup
}

// NewService creates a new review Service.
func NewService(repo Repository, lookup BookLookup) *Service {
	return &Service{repo: repo, lookup: lookup}
}

// Reviews returns all reviews for a given book.
func (s *Service) Reviews(_ context.Context, bookID int) ([]Review, error) {
	return s.repo.ByBookID(bookID)
}

// CreateReview adds a new review.
func (s *Service) CreateReview(ctx context.Context, rv Review) error {
	if err := s.lookup(ctx, rv.BookID); err != nil {
		return err
	}
	return s.repo.Create(rv)
}
