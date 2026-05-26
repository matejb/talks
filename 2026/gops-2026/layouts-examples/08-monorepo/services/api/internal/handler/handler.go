package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/08-monorepo/libs/library"
)

type Handler struct {
	Store *library.Store
}

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
	users, _ := h.Store.Users()
	writeJSON(w, users)
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	user, err := h.Store.User(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, user)
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var u library.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	created, err := h.Store.CreateUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, created)
}

func (h *Handler) listBooks(w http.ResponseWriter, _ *http.Request) {
	books, _ := h.Store.Books()
	writeJSON(w, books)
}

func (h *Handler) getBook(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	book, err := h.Store.Book(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, book)
}

func (h *Handler) borrowBook(w http.ResponseWriter, r *http.Request) {
	bookID, _ := strconv.Atoi(r.PathValue("id"))
	var body struct{ UserID int }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.Store.BorrowBook(bookID, body.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) returnBook(w http.ResponseWriter, r *http.Request) {
	bookID, _ := strconv.Atoi(r.PathValue("id"))
	if err := h.Store.ReturnBook(bookID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) listReviews(w http.ResponseWriter, r *http.Request) {
	bookID, _ := strconv.Atoi(r.PathValue("id"))
	reviews, _ := h.Store.Reviews(bookID)
	writeJSON(w, reviews)
}

func (h *Handler) createReview(w http.ResponseWriter, r *http.Request) {
	var rev library.Review
	if err := json.NewDecoder(r.Body).Decode(&rev); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	created, err := h.Store.CreateReview(rev)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, created)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
