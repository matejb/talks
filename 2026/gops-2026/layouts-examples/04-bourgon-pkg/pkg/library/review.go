package library

// Review represents a book review.
type Review struct {
	ID     int
	BookID int
	UserID int
	Text   string
	Rating int
}

// ReviewService provides operations on reviews.
type ReviewService interface {
	Reviews(bookID int) ([]Review, error)
	CreateReview(review Review) error
}
