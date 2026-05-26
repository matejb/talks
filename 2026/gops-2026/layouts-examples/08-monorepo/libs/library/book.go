package library

// Book represents a library book.
type Book struct {
	ID       int
	Title    string
	Author   string
	Borrowed *User
}
