package users

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Handler provides HTTP endpoints for users.
type Handler struct {
	Service *Service
}

// RegisterRoutes registers user routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /users", h.list)
	mux.HandleFunc("GET /users/{id}", h.get)
	mux.HandleFunc("POST /users", h.create)
}

func (h *Handler) list(w http.ResponseWriter, _ *http.Request) {
	users, _ := h.Service.Users()
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	user, err := h.Service.User(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.Service.CreateUser(u)
	w.WriteHeader(http.StatusCreated)
}
