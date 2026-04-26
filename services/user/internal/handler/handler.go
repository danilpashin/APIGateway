package handler

import (
	"apigateway/services/user/internal/domain"
	"apigateway/services/user/internal/service"
	"context"
	"encoding/json"
	"net/http"
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

	user, err := h.service.CreateUser(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := domain.CreateUserResponse{Username: user.Username, Email: user.Email, CreatedAt: user.CreatedAt}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {}
