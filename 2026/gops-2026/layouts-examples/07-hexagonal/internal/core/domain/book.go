package domain

// Book represents a library book.
type Book struct {
	ID       int
	Title    string
	Author   string
	Borrower *User
}
