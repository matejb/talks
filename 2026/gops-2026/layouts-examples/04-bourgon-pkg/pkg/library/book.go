package library

// Book represents a library book.
type Book struct {
	ID       int
	Title    string
	Author   string
	Borrowed *User
}

// BookService provides operations on books.
type BookService interface {
	Book(id int) (Book, error)
	Books() ([]Book, error)
	Borrow(bookID, userID int) error
	Return(bookID int) error
}
