package handler

import (
	"apigateway/services/product/internal/domain"
	"apigateway/services/product/internal/service"
	"apigateway/services/product/internal/validator"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"pkg/response"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v3"
)

type ProductHandler struct {
	service service.ProductServiceInterface
}

func NewProductHandler(service service.ProductServiceInterface) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	httplog.SetAttrs(r.Context(), slog.String("handler", "CreateProduct"))

	var req domain.CreateProductRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.handleError(r.Context(), w, domain.ErrInvalidJSON)
		return
	}

	if err = validator.New(req); err != nil {
		errResp := domain.ErrorResponse{Message: "validation error", Details: response.FormatValidationError(err)}
		h.handleError(r.Context(), w, errResp)
		return
	}

	product, err := h.service.CreateProduct(r.Context(), &req)
	if err != nil {
		h.handleError(r.Context(), w, err, slog.Any("req", req))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	httplog.SetAttrs(r.Context(), slog.String("handler", "UpdateProduct"))

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		h.handleError(r.Context(), w, domain.ErrIDRequired)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.handleError(r.Context(), w, domain.ErrInvalidID)
		return
	}

	var req *domain.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(r.Context(), w, domain.ErrInvalidJSON)
		return
	}

	product, err := h.service.UpdateProduct(r.Context(), id, req)
	if err != nil {
		h.handleError(r.Context(), w, err, slog.Int("id", id))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	httplog.SetAttrs(r.Context(), slog.String("handler", "GetProduct"))

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		h.handleError(r.Context(), w, domain.ErrIDRequired)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.handleError(r.Context(), w, domain.ErrInvalidID)
	}

	product, err := h.service.GetProduct(r.Context(), id)
	if err != nil {
		h.handleError(r.Context(), w, err, slog.Int("id", id))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	httplog.SetAttrs(r.Context(), slog.String("handler", "ListProducts"))

	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "10"
	}
	cursor, err := strconv.Atoi(cursorStr)
	if err != nil {
		h.handleError(r.Context(), w, domain.ErrInvalidCursor)
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.handleError(r.Context(), w, domain.ErrInvalidLimit)
		return
	}

	products, pagination, err := h.service.ListProducts(r.Context(), cursor, uint64(limit))
	if err != nil {
		h.handleError(r.Context(), w, err, slog.Int("cursor", cursor), slog.Int("limit", limit))
		return
	}

	listProductsResponse := domain.ListProductsResponse{Products: products, PaginationParams: pagination, Total: len(products)}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(listProductsResponse)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	httplog.SetAttrs(r.Context(), slog.String("handler", "DeleteProduct"))

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

	err = h.service.DeleteProduct(r.Context(), idInt)
	if err != nil {
		h.handleError(r.Context(), w, err, slog.Int("id", idInt))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) handleError(ctx context.Context, w http.ResponseWriter, err error, attrs ...slog.Attr) {
	var statusCode int
	var errResp domain.ErrorResponse

	switch {
	case errors.As(err, &errResp):
		statusCode = http.StatusBadRequest

	case errors.Is(err, domain.ErrProductsNotFound):
		statusCode = http.StatusNotFound
		errResp = domain.ErrorResponse{Message: err.Error()}

	case errors.Is(err, domain.ErrProductExist):
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

	case errors.Is(err, domain.ErrListQuery):
		statusCode = http.StatusBadRequest
		errResp = domain.ErrorResponse{Message: err.Error()}

	case errors.Is(err, domain.ErrNoUpdateData):
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
