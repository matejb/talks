package books

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Handler provides HTTP endpoints for books.
type Handler struct {
	Service *Service
}

// RegisterRoutes registers book routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /books", h.list)
	mux.HandleFunc("GET /books/{id}", h.get)
	mux.HandleFunc("POST /books/{id}/borrow", h.borrow)
	mux.HandleFunc("POST /books/{id}/return", h.returnBook)
}

func (h *Handler) list(w http.ResponseWriter, _ *http.Request) {
	books, _ := h.Service.Books()
	json.NewEncoder(w).Encode(books)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	book, err := h.Service.Book(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(book)
}

func (h *Handler) borrow(w http.ResponseWriter, r *http.Request) {
	bookID, _ := strconv.Atoi(r.PathValue("id"))
	var body struct{ UserID int }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.Service.Borrow(bookID, body.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) returnBook(w http.ResponseWriter, r *http.Request) {
	bookID, _ := strconv.Atoi(r.PathValue("id"))
	if err := h.Service.Return(bookID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
