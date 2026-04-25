package handler

import (
	"apigateway/services/product/internal/domain"
	"apigateway/services/product/internal/service"
	"apigateway/services/product/internal/validator"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"pkg/response"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req domain.CreateProductRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.handleError(w, domain.ErrInvalidJSON)
		return
	}

	if err = validator.New(req); err != nil {
		errResp := domain.ErrorResponse{Message: "validation error", Details: response.FormatValidationError(err)}
		h.handleError(w, errResp)
		return
	}

	product, err := h.productService.CreateProduct(r.Context(), &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
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

	var req *domain.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, domain.ErrInvalidJSON)
		return
	}

	product, err := h.productService.UpdateProduct(r.Context(), idInt, req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id := chi.URLParam(r, "id")
	if id == "" {
		h.handleError(w, domain.ErrIDRequired)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.handleError(w, domain.ErrInvalidID)
	}

	product, err := h.productService.GetProduct(r.Context(), idInt)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cursor := r.URL.Query().Get("cursor")
	if cursor == "" {
		cursor = "0"
	}
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "10"
	}
	cursorInt, err := strconv.Atoi(cursor)
	if err != nil {
		h.handleError(w, domain.ErrInvalidCursor)
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		h.handleError(w, domain.ErrInvalidLimit)
		return
	}

	products, pagination, err := h.productService.ListProducts(r.Context(), cursorInt, uint64(limitInt))
	if err != nil {
		h.handleError(w, err)
		return
	}

	listProductsResponse := domain.ListProductsResponse{Products: products, PaginationParams: pagination, Total: len(products)}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(listProductsResponse)
}

func (h *ProductHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
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

	err = h.productService.DeleteProduct(r.Context(), idInt)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) handleError(w http.ResponseWriter, err error) {
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
		log.Printf("unexpected error: %v", err)
	}

	JSONError(w, statusCode, errResp)
}

func JSONError(w http.ResponseWriter, statusCode int, err domain.ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(err)
}
