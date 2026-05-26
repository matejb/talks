package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/core/service"
)

// BookHandler handles HTTP requests for books.
type BookHandler struct {
	svc *service.BookService
}

// NewBookHandler creates a new BookHandler.
func NewBookHandler(svc *service.BookService) *BookHandler {
	return &BookHandler{svc: svc}
}

// Routes returns the HTTP route patterns for books.
func (h *BookHandler) Routes() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"GET /books":              h.list,
		"GET /books/{id}":         h.get,
		"POST /books/{id}/borrow": h.borrow,
		"POST /books/{id}/return": h.returnBook,
	}
}

func (h *BookHandler) list(w http.ResponseWriter, r *http.Request) {
	books, err := h.svc.All(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	book, err := h.svc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) borrow(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var body struct{ UserID int }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.svc.Borrow(r.Context(), bookID, body.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *BookHandler) returnBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.svc.Return(r.Context(), bookID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
