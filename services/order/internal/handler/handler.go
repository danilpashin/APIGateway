package handler

import (
	"apigateway/services/order/internal/domain"
	"apigateway/services/order/internal/service"
	"encoding/json"
	"net/http"
)

type OrderHandler struct {
	service service.OrderServiceInterface
}

func NewOrderHandler(service service.OrderServiceInterface) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) AddProductToCart(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req *domain.AddProductToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	order, err := h.service.AddProductToCart(r.Context(), &domain.AddProductToCartRequest{UserID: req.UserID, ProductID: req.ProductID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}
