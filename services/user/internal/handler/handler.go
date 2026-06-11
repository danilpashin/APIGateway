package handler

import (
	"apigateway/services/user/internal/domain"
	"apigateway/services/user/internal/service"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v3"
)

type UserHandler struct {
	service service.UserServiceInterface
}

func NewUserHandler(service service.UserServiceInterface) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	httplog.SetAttrs(r.Context(), slog.String("handler", "CreateUser"))

	var req domain.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.handleError(r.Context(), w, domain.ErrInvalidJSON)
		return
	}

	user, err := h.service.CreateUser(r.Context(), &req)
	if err != nil {
		h.handleError(r.Context(), w, err, slog.Any("req", req))
		return
	}

	resp := domain.CreateUserResponse{Username: user.Username, Email: user.Email, CreatedAt: user.CreatedAt}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	httplog.SetAttrs(r.Context(), slog.String("handler", "UpdateUser"))

	id := chi.URLParam(r, "id")
	if id == "" {
		h.handleError(r.Context(), w, domain.ErrIDRequired)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.handleError(r.Context(), w, domain.ErrInvalidID)
		return
	}

	var req domain.UpdateUserRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.handleError(r.Context(), w, domain.ErrInvalidJSON)
		return
	}

	user, err := h.service.UpdateUser(r.Context(), idInt, &req)
	if err != nil {
		h.handleError(r.Context(), w, err, slog.Int("id", idInt), slog.Any("req", req))
		return
	}

	resp := domain.UpdateUserResponse{Username: user.Username, Email: user.Email}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	httplog.SetAttrs(r.Context(), slog.String("handler", "GetUser"))

	id := chi.URLParam(r, "id")
	if id == "" {
		h.handleError(r.Context(), w, domain.ErrIDRequired)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.handleError(r.Context(), w, domain.ErrInvalidID)
		return
	}

	user, err := h.service.GetUser(r.Context(), idInt)
	if err != nil {
		h.handleError(r.Context(), w, err, slog.Int("id", idInt))
		return
	}

	resp := domain.GetUserResponse{Username: user.Username, CreatedAt: user.CreatedAt}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	httplog.SetAttrs(r.Context(), slog.String("handler", "ListUsers"))

	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}
	cursor, err := strconv.Atoi(cursorStr)
	if err != nil {
		h.handleError(r.Context(), w, domain.ErrInvalidCursor)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "10"
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.handleError(r.Context(), w, domain.ErrInvalidLimit)
		return
	}

	listUsers, pagination, err := h.service.ListUsers(r.Context(), cursor, uint64(limit))
	if err != nil {
		h.handleError(r.Context(), w, err, slog.Int("cursor", cursor), slog.Int("limit", limit))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.ListUsersResponse{
		Users:      listUsers,
		Pagination: pagination,
	})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	httplog.SetAttrs(r.Context(), slog.String("handler", "DeleteUser"))

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		h.handleError(r.Context(), w, domain.ErrIDRequired)
		return
	}
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		h.handleError(r.Context(), w, domain.ErrInvalidID)
		return
	}

	err = h.service.DeleteUser(r.Context(), idInt)
	if err != nil {
		h.handleError(r.Context(), w, err, slog.Int("userID", idInt))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) handleError(ctx context.Context, w http.ResponseWriter, err error, attrs ...slog.Attr) {
	var statusCode int
	var errResp domain.ErrorResponse

	switch {
	case errors.As(err, &errResp):
		statusCode = http.StatusBadRequest

	case errors.Is(err, domain.ErrUserNotFound):
		statusCode = http.StatusNotFound
		errResp = domain.ErrorResponse{Message: err.Error()}

	case errors.Is(err, domain.ErrUserExist):
		statusCode = http.StatusConflict
		errResp = domain.ErrorResponse{Message: err.Error()}

	case errors.Is(err, domain.ErrInvalidJSON):
		statusCode = http.StatusBadRequest
		errResp = domain.ErrorResponse{Message: err.Error()}

	case errors.Is(err, domain.ErrIDRequired):
		statusCode = http.StatusBadRequest
		errResp = domain.ErrorResponse{Message: err.Error()}

	case errors.Is(err, domain.ErrInvalidID):
		statusCode = http.StatusBadRequest
		errResp = domain.ErrorResponse{Message: err.Error()}

	case errors.Is(err, domain.ErrInvalidCursor):
		statusCode = http.StatusBadRequest
		errResp = domain.ErrorResponse{Message: err.Error()}

	case errors.Is(err, domain.ErrInvalidLimit):
		statusCode = http.StatusBadRequest
		errResp = domain.ErrorResponse{Message: err.Error()}

	case errors.Is(err, domain.ErrNoInsertData):
		statusCode = http.StatusBadRequest
		errResp = domain.ErrorResponse{Message: err.Error()}

	default:
		statusCode = http.StatusInternalServerError
		errResp = domain.ErrorResponse{Message: "internal server error"}

		if len(attrs) > 0 {
			httplog.SetAttrs(ctx, attrs...)
		}
		httplog.SetError(ctx, err)
	}

	JSONError(w, statusCode, errResp)
}

func JSONError(w http.ResponseWriter, statusCode int, err domain.ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(err)
}
