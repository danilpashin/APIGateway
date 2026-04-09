package postgres

import (
	"apigateway/services/product/internal/domain"
	"context"
	"database/sql"
)

type ProductRepository struct {
	db *sql.DB
}

func NewPostgresProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, insertData map[string]interface{}) (*domain.Product, error) {
	var product domain.Product

	return &product, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, id int, updateData map[string]interface{}) (*domain.Product, error) {
	var product domain.Product

	return &product, nil
}

func (r *ProductRepository) GetProduct(ctx context.Context, id int) (*domain.Product, error) {
	var product domain.Product

	return &product, nil
}

func (r *ProductRepository) ListProducts(ctx context.Context, cursor int, limit uint64) ([]*domain.Product, error) {
	listProducts := make([]*domain.Product, 0)

	return listProducts, nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id int) error {
	return nil
}
