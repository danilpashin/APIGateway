package handler

import (
	"apigateway/services/user/internal/domain"
	"apigateway/services/user/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Basic handler check"))
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req domain.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.service.CreateUser(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := domain.CreateUserResponse{Username: user.Username, Email: user.Email, CreatedAt: user.CreatedAt}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "empty id", http.StatusBadRequest)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req domain.UpdateUserRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.service.UpdateUser(r.Context(), idInt, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := domain.UpdateUserResponse{Username: user.Username, Email: user.Email}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "empty id", http.StatusBadRequest)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUser(r.Context(), idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := domain.GetUserResponse{Username: user.Username, CreatedAt: user.CreatedAt}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
