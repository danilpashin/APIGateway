package service

import (
	"apigateway/services/product/internal/domain"
	"apigateway/services/product/internal/repository/postgres"
	"context"
	"errors"
	"time"
)

type ProductService struct {
	repo postgres.ProductRepoInterface
}

func NewProductService(repo postgres.ProductRepoInterface) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, req *domain.CreateProductRequest) (*domain.Product, error) {
	insertData := make(map[string]interface{}, 0)

	return s.repo.CreateProduct(ctx, insertData)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id int, req *domain.UpdateProductRequest) (*domain.Product, error) {
	updateData := make(map[string]interface{}, 0)

	updateData["updated_at"] = time.Now()

	return s.repo.UpdateProduct(ctx, id, updateData)
}

func (s *ProductService) GetProduct(ctx context.Context, id int) (*domain.Product, error) {
	if id <= 0 {
		return nil, errors.New("id must be positive")
	}

	return s.repo.GetProduct(ctx, id)
}

func (s *ProductService) ListProducts(ctx context.Context, cursor int, limit uint64) ([]*domain.Product, error) {
	if cursor < 0 {
		cursor = 0
	}
	if limit <= 0 {
		limit = 10
	} else if limit > 50 {
		limit = 50
	}

	return s.repo.ListProducts(ctx, cursor, limit)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("id must be positive")
	}

	return s.repo.DeleteProduct(ctx, id)
}
