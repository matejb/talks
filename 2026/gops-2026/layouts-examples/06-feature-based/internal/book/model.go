package book

// Book represents a library book.
type Book struct {
	ID         int
	Title      string
	Author     string
	BorrowerID int // 0 means not borrowed
}
