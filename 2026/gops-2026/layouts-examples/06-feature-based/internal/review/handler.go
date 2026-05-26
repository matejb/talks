package review

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Handler serves HTTP requests for reviews.
type Handler struct {
	svc *Service
}

// NewHandler creates a new review Handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Routes returns the route definitions for this feature.
func (h *Handler) Routes() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"GET /books/{id}/reviews": h.list,
		"POST /reviews":           h.create,
	}
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	reviews, err := h.svc.Reviews(r.Context(), bookID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(reviews)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var rv Review
	if err := json.NewDecoder(r.Body).Decode(&rv); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.svc.CreateReview(r.Context(), rv); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
