package postgres

import (
	"apigateway/services/order/internal/domain"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

func (r *OrderRepository) AddOrder(ctx context.Context, order *domain.Order) error {
	return nil
}
