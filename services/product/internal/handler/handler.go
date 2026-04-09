package handler

import (
	"apigateway/services/product/internal/domain"
	"apigateway/services/product/internal/service"
	"errors"
	"log"
	"net/http"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
}

func (h *ProductHandler) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
}

func (h *ProductHandler) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
}

func (h *ProductHandler) ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
}

func (h *ProductHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
}

func (h *ProductHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrProductNotFound):
		http.Error(w, "Product not found", http.StatusNotFound)

	case errors.Is(err, domain.ErrForbidden):
		http.Error(w, "Forbidden", http.StatusForbidden)

	default:
		log.Printf("unexpected error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
