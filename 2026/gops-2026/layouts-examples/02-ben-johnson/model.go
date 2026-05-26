package library

// User represents a library member.
type User struct {
	ID    int
	Name  string
	Email string
}

// Book represents a library book.
type Book struct {
	ID       int
	Title    string
	Author   string
	Borrowed *User
}

// Review represents a book review.
type Review struct {
	ID     int
	BookID int
	UserID int
	Text   string
	Rating int
}
