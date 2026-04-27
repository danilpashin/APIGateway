package handler

import (
	"apigateway/services/user/internal/domain"
	"apigateway/services/user/internal/service"
	"encoding/json"
	"errors"
	"log"
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
		h.handleError(w, domain.ErrInvalidJSON)
		return
	}

	user, err := h.service.CreateUser(r.Context(), &req)
	if err != nil {
		h.handleError(w, err)
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
		h.handleError(w, domain.ErrIDRequired)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.handleError(w, domain.ErrInvalidID)
		return
	}

	var req domain.UpdateUserRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.handleError(w, domain.ErrInvalidJSON)
		return
	}

	user, err := h.service.UpdateUser(r.Context(), idInt, &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	resp := domain.UpdateUserResponse{Username: user.Username, Email: user.Email}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id := chi.URLParam(r, "id")
	if id == "" {
		h.handleError(w, domain.ErrIDRequired)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.handleError(w, domain.ErrInvalidID)
		return
	}

	user, err := h.service.GetUser(r.Context(), idInt)
	if err != nil {
		h.handleError(w, err)
		return
	}

	resp := domain.GetUserResponse{Username: user.Username, CreatedAt: user.CreatedAt}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		http.Error(w, domain.ErrUserNotFound.Error(), http.StatusNotFound)
		return

	case errors.Is(err, domain.ErrUserExist):
		http.Error(w, domain.ErrUserExist.Error(), http.StatusConflict)
		return

	case errors.Is(err, domain.ErrRoleNotFound):
		http.Error(w, domain.ErrRoleNotFound.Error(), http.StatusNotFound)
		return

	case errors.Is(err, domain.ErrIDRequired):
		http.Error(w, domain.ErrIDRequired.Error(), http.StatusBadRequest)
		return

	case errors.Is(err, domain.ErrWrongPassword):
		http.Error(w, domain.ErrWrongPassword.Error(), http.StatusBadRequest)
		return

	case errors.Is(err, domain.ErrNoInsertData):
		http.Error(w, domain.ErrNoInsertData.Error(), http.StatusBadRequest)
		return

	case errors.Is(err, domain.ErrInvalidID):
		http.Error(w, domain.ErrInvalidID.Error(), http.StatusBadRequest)
		return

	case errors.Is(err, domain.ErrInvalidJSON):
		http.Error(w, domain.ErrInvalidJSON.Error(), http.StatusBadRequest)
		return

	default:
		log.Print("internal server error: ", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
