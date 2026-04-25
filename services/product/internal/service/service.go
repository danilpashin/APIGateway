package service

import (
	"apigateway/services/product/internal/domain"
	"apigateway/services/product/internal/helpers/regex"
	"apigateway/services/product/internal/repository/postgres"
	"context"
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

	if req.Name != "" {
		err := regex.ValidateProductName(req.Name)
		if err != nil {
			return nil, domain.ErrInvalidName
		}
		insertData["name"] = req.Name
	}
	if req.Manufacturer != "" {
		insertData["manufacturer"] = req.Manufacturer
	}
	if req.Price > 0 {
		insertData["price"] = req.Price
	}
	if req.Amount > 0 {
		insertData["amount"] = req.Amount
	}
	if req.Status {
		insertData["status"] = req.Status
	} else {
		insertData["status"] = true
	}
	if req.Category != "" {
		insertData["category"] = req.Category
	}

	return s.repo.CreateProduct(ctx, insertData)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id int, req *domain.UpdateProductRequest) (*domain.Product, error) {
	updateData := make(map[string]interface{}, 0)

	if req.Name != nil {
		err := regex.ValidateProductName(*req.Name)
		if err != nil {
			return nil, err
		}
		updateData["name"] = *req.Name
	}
	if req.Manufacturer != nil {
		if *req.Manufacturer == "" {
			return nil, domain.ErrInvalidManufacturer
		}
		updateData["manufacturer"] = *req.Manufacturer
	}
	if req.Amount != nil {
		if *req.Amount < 0 {
			return nil, domain.ErrInvalidAmount
		}
		updateData["amount"] = *req.Amount
	}
	if req.Price != nil {
		if *req.Price <= 0 {
			return nil, domain.ErrInvalidPrice
		}
		updateData["price"] = *req.Price
	}
	if req.Category != nil {
		if *req.Category == "" {
			return nil, domain.ErrInvalidPrice
		}
		updateData["price"] = *req.Price
	}

	if len(updateData) == 0 {
		return nil, domain.ErrNoUpdateData
	}

	updateData["updated_at"] = time.Now()

	return s.repo.UpdateProduct(ctx, id, updateData)
}

func (s *ProductService) GetProduct(ctx context.Context, id int) (*domain.Product, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidID
	}

	return s.repo.GetProduct(ctx, id)
}

func (s *ProductService) ListProducts(ctx context.Context, cursor int, limit uint64) ([]*domain.Product, *domain.Pagination, error) {
	if cursor < 0 {
		cursor = 0
	}
	if limit <= 0 {
		limit = 10
	} else if limit > 50 {
		limit = 50
	}

	listProducts, nextCursor, hasMore, err := s.repo.ListProducts(ctx, cursor, limit)
	if err != nil {
		return nil, nil, err
	}

	return listProducts, &domain.Pagination{NextCursor: nextCursor, HasMore: hasMore, Limit: limit}, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	if id <= 0 {
		return domain.ErrInvalidID
	}

	return s.repo.DeleteProduct(ctx, id)
}
