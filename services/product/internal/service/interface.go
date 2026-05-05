package service

import (
	"apigateway/services/product/internal/domain"
	"context"
)

type ProductServiceInterface interface {
	CreateProduct(ctx context.Context, req *domain.CreateProductRequest) (*domain.Product, error)
	UpdateProduct(ctx context.Context, id int, req *domain.UpdateProductRequest) (*domain.Product, error)
	GetProduct(ctx context.Context, id int) (*domain.Product, error)
	ListProducts(ctx context.Context, cursor int, limit uint64) ([]*domain.Product, *domain.Pagination, error)
	DeleteProduct(ctx context.Context, id int) error
}
