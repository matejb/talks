package library

import "fmt"

// Review represents a book review.
type Review struct {
	ID     int
	BookID int
	UserID int
	Text   string
	Rating int
}

type ReviewStore interface {
	ListReviews(bookID int) ([]Review, error)
	CreateReview(review Review) error
}

type ReviewService struct {
	Store ReviewStore
}

func (s *ReviewService) Reviews(bookID int) ([]Review, error) {
	return s.Store.ListReviews(bookID)
}

func (s *ReviewService) CreateReview(review Review) error {
	if review.BookID == 0 {
		return fmt.Errorf("book ID is required")
	}
	if review.Rating < 1 || review.Rating > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}
	return s.Store.CreateReview(review)
}
