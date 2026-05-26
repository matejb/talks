package library

// UserService provides access to library users.
type UserService interface {
	User(id int) (User, error)
	Users() ([]User, error)
	CreateUser(user User) error
}

// BookService provides access to library books and borrowing operations.
type BookService interface {
	Book(id int) (Book, error)
	Books() ([]Book, error)
	Borrow(bookID, userID int) error
	Return(bookID int) error
}

// ReviewService provides access to book reviews.
type ReviewService interface {
	Reviews(bookID int) ([]Review, error)
	CreateReview(review Review) error
}
