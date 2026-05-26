package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/core/domain"
	"github.com/matejb/talks/2026/gops-2026/layouts-examples/07-hexagonal/internal/core/service"
)

// ReviewHandler handles HTTP requests for reviews.
type ReviewHandler struct {
	svc *service.ReviewService
}

// NewReviewHandler creates a new ReviewHandler.
func NewReviewHandler(svc *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{svc: svc}
}

// Routes returns the HTTP route patterns for reviews.
func (h *ReviewHandler) Routes() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"GET /books/{id}/reviews": h.listByBook,
		"POST /reviews":           h.create,
	}
}

func (h *ReviewHandler) listByBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	reviews, err := h.svc.ByBook(r.Context(), bookID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(reviews)
}

func (h *ReviewHandler) create(w http.ResponseWriter, r *http.Request) {
	var rev domain.Review
	if err := json.NewDecoder(r.Body).Decode(&rev); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.svc.Create(r.Context(), rev); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
