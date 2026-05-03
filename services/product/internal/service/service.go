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
	insertData := make(map[string]any, 0)

	if req.Name != "" {
		err := regex.ValidateName(req.Name)
		if err != nil {
			return nil, domain.ErrInvalidName
		}
		insertData["name"] = req.Name
	} else {
		return nil, domain.ErrNameRequired
	}

	if req.Manufacturer != "" {
		err := regex.ValidateManufacturer(req.Manufacturer)
		if err != nil {
			return nil, domain.ErrInvalidManufacturer
		}
		insertData["manufacturer"] = req.Manufacturer
	} else {
		return nil, domain.ErrManufacturerRequired
	}

	if req.Price > 0 {
		insertData["price"] = req.Price
	} else {
		return nil, domain.ErrInvalidPrice
	}

	if req.Amount > 0 {
		insertData["amount"] = req.Amount
	} else {
		return nil, domain.ErrInvalidAmount
	}

	if req.Category != "" {
		err := regex.ValidateCategory(req.Category)
		if err != nil {
			return nil, domain.ErrInvalidCategory
		}
		insertData["category"] = req.Category
	} else {
		return nil, domain.ErrCategoryRequired
	}

	insertData["status"] = req.Status

	return s.repo.CreateProduct(ctx, insertData)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id int, req *domain.UpdateProductRequest) (*domain.Product, error) {
	updateData := make(map[string]any, 0)

	if req.Name != nil {
		err := regex.ValidateName(*req.Name)
		if err != nil {
			return nil, domain.ErrInvalidName
		}
		updateData["name"] = *req.Name
	}
	if req.Manufacturer != nil {
		err := regex.ValidateManufacturer(*req.Manufacturer)
		if err != nil {
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
		err := regex.ValidateCategory(*req.Category)
		if err != nil {
			return nil, domain.ErrInvalidCategory
		}
		updateData["category"] = *req.Category
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
