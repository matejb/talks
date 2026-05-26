package library

import "fmt"

// Book represents a library book.
type Book struct {
	ID       int
	Title    string
	Author   string
	Borrowed *User
}

type BookStore interface {
	GetBook(id int) (Book, error)
	ListBooks() ([]Book, error)
	SaveBook(book Book) error
}

type BookService struct {
	Store BookStore
	Users UserStore
}

func (s *BookService) Book(id int) (Book, error) {
	return s.Store.GetBook(id)
}

func (s *BookService) Books() ([]Book, error) {
	return s.Store.ListBooks()
}

func (s *BookService) Borrow(bookID, userID int) error {
	book, err := s.Store.GetBook(bookID)
	if err != nil {
		return fmt.Errorf("find book: %w", err)
	}
	if book.Borrowed != nil {
		return fmt.Errorf("book %d already borrowed by %s", bookID, book.Borrowed.Name)
	}
	user, err := s.Users.GetUser(userID)
	if err != nil {
		return fmt.Errorf("find user: %w", err)
	}
	book.Borrowed = &User{ID: user.ID, Name: user.Name, Email: user.Email}
	return s.Store.SaveBook(book)
}

func (s *BookService) Return(bookID int) error {
	book, err := s.Store.GetBook(bookID)
	if err != nil {
		return fmt.Errorf("find book: %w", err)
	}
	book.Borrowed = nil
	return s.Store.SaveBook(book)
}
