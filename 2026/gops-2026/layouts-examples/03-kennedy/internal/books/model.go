package books

// Book represents a library book.
// In the Kennedy layout, domain types live in their domain package.
// BorrowerID references the user who borrowed the book (0 = available).
type Book struct {
	ID           int
	Title        string
	Author       string
	BorrowerID   int    // 0 means not borrowed
	BorrowerName string // denormalized for display
}
