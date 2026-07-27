package postgres

import (
	"apigateway/services/order/internal/domain"
	"context"
)

type OrderRepositoryInterface interface {
	AddOrder(ctx context.Context, order *domain.Order) error
}
