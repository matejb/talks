package httphandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/04-bourgon-pkg/pkg/library"
)

// Handler holds references to all services and provides HTTP routes.
type Handler struct {
	Users   library.UserService
	Books   library.BookService
	Reviews library.ReviewService
}

// RegisterRoutes registers all routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /users", h.listUsers)
	mux.HandleFunc("GET /users/{id}", h.getUser)
	mux.HandleFunc("POST /users", h.createUser)
	mux.HandleFunc("GET /books", h.listBooks)
	mux.HandleFunc("GET /books/{id}", h.getBook)
	mux.HandleFunc("POST /books/{id}/borrow", h.borrowBook)
	mux.HandleFunc("POST /books/{id}/return", h.returnBook)
	mux.HandleFunc("GET /books/{id}/reviews", h.listReviews)
	mux.HandleFunc("POST /reviews", h.createReview)
}

func (h *Handler) listUsers(w http.ResponseWriter, _ *http.Request) {
	users, _ := h.Users.Users()
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	user, err := h.Users.User(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var u library.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.Users.CreateUser(u)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) listBooks(w http.ResponseWriter, _ *http.Request) {
	books, _ := h.Books.Books()
	json.NewEncoder(w).Encode(books)
}

func (h *Handler) getBook(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	book, err := h.Books.Book(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(book)
}

func (h *Handler) borrowBook(w http.ResponseWriter, r *http.Request) {
	bookID, _ := strconv.Atoi(r.PathValue("id"))
	var body struct{ UserID int }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.Books.Borrow(bookID, body.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) returnBook(w http.ResponseWriter, r *http.Request) {
	bookID, _ := strconv.Atoi(r.PathValue("id"))
	if err := h.Books.Return(bookID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) listReviews(w http.ResponseWriter, r *http.Request) {
	bookID, _ := strconv.Atoi(r.PathValue("id"))
	reviews, _ := h.Reviews.Reviews(bookID)
	json.NewEncoder(w).Encode(reviews)
}

func (h *Handler) createReview(w http.ResponseWriter, r *http.Request) {
	var rev library.Review
	if err := json.NewDecoder(r.Body).Decode(&rev); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.Reviews.CreateReview(rev)
	w.WriteHeader(http.StatusCreated)
}
