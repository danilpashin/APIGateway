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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = validator.New(req); err != nil {
		response.WriteValidationError(w, err)
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
		http.Error(w, "empty id", http.StatusBadRequest)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "id is not a number", http.StatusBadRequest)
		return
	}

	var req *domain.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	product, err := h.productService.UpdateProduct(r.Context(), idInt, req)
	if err != nil {
		http.Error(w, "failed update", http.StatusInternalServerError)
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
		http.Error(w, "no id to get the product", http.StatusBadRequest)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "id is not a number", http.StatusBadRequest)
		return
	}

	product, err := h.productService.GetProduct(r.Context(), idInt)
	if err != nil {
		h.handleError(w, err)
		return
	}

	if product == nil {
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
		http.Error(w, "cursor is not a number", http.StatusBadRequest)
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "limit is not a number", http.StatusBadRequest)
		return
	}

	product, err := h.productService.ListProducts(r.Context(), cursorInt, uint64(limitInt))
	if err != nil {
		h.handleError(w, err)
		return
	}

	if product == nil {
		w.WriteHeader(http.StatusNoContent)
		h.handleError(w, domain.ErrProductsNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "empty id", http.StatusBadRequest)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "id is not a number", http.StatusBadRequest)
		return
	}

	err = h.productService.DeleteProduct(r.Context(), idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrProductsNotFound):
		JSONError(w, http.StatusNotFound, "Products not found")

	case errors.Is(err, domain.ErrProductExist):
		JSONError(w, http.StatusConflict, "Product already exists")

	case errors.Is(err, domain.ErrForbidden):
		JSONError(w, http.StatusForbidden, "Forbidden")

	case errors.Is(err, domain.ErrNameRequired):
		JSONError(w, http.StatusBadRequest, "Name is required")

	case errors.Is(err, domain.ErrManufacturerRequired):
		JSONError(w, http.StatusBadRequest, "Manufacturer is required")

	case errors.Is(err, domain.ErrPriceRequired):
		JSONError(w, http.StatusBadRequest, "Price is required")

	case errors.Is(err, domain.ErrAmountRequired):
		JSONError(w, http.StatusBadRequest, "Amount is required")

	case errors.Is(err, domain.ErrCategoryRequired):
		JSONError(w, http.StatusBadRequest, "Category is required")

	default:
		log.Printf("unexpected error: %v", err)
		JSONError(w, http.StatusInternalServerError, "Internal server error")
	}
}

func JSONError(w http.ResponseWriter, statusCode int, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": err,
	})
}
