package library

// Review represents a book review.
type Review struct {
	ID     int
	BookID int
	UserID int
	Text   string
	Rating int
}
