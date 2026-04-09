package postgres

import (
	"apigateway/services/product/internal/domain"
	"context"
)

type ProductRepoInterface interface {
	CreateProduct(ctx context.Context, insertData map[string]interface{}) (*domain.Product, error)
	UpdateProduct(ctx context.Context, id int, updateData map[string]interface{}) (*domain.Product, error)
	GetProduct(ctx context.Context, id int) (*domain.Product, error)
	ListProducts(ctx context.Context, cursor int, limit uint64) ([]*domain.Product, error)
	DeleteProduct(ctx context.Context, id int) error
}
