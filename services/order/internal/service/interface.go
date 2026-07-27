package service

import (
	"apigateway/services/order/internal/domain"
	"context"
)

type OrderServiceInterface interface {
	AddProductToCart(ctx context.Context, req *domain.AddProductToCartRequest) (*domain.Order, error)
}
