package reviews

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Handler provides HTTP endpoints for reviews.
type Handler struct {
	Service *Service
}

// RegisterRoutes registers review routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /books/{id}/reviews", h.list)
	mux.HandleFunc("POST /reviews", h.create)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	bookID, _ := strconv.Atoi(r.PathValue("id"))
	reviews, _ := h.Service.Reviews(bookID)
	json.NewEncoder(w).Encode(reviews)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var rev Review
	if err := json.NewDecoder(r.Body).Decode(&rev); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.Service.CreateReview(rev)
	w.WriteHeader(http.StatusCreated)
}
