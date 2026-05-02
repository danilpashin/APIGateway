package service

import (
	"apigateway/services/product/internal/domain"
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type MockProductRepo struct {
	createProduct func(ctx context.Context, insertData map[string]any) (*domain.Product, error)
	updateProduct func(ctx context.Context, id int, updateData map[string]any) (*domain.Product, error)
	getProduct    func(ctx context.Context, id int) (*domain.Product, error)
	listProducts  func(ctx context.Context, cursor int, limit uint64) ([]*domain.Product, int, bool, error)
	deleteProduct func(ctx context.Context, id int) error
}

func (m *MockProductRepo) CreateProduct(ctx context.Context, insertData map[string]any) (*domain.Product, error) {
	if m.createProduct != nil {
		return m.createProduct(ctx, insertData)
	}
	return nil, nil
}

func (m *MockProductRepo) UpdateProduct(ctx context.Context, id int, updateData map[string]any) (*domain.Product, error) {
	if m.updateProduct != nil {
		return m.updateProduct(ctx, id, updateData)
	}
	return nil, nil
}

func (m *MockProductRepo) GetProduct(ctx context.Context, id int) (*domain.Product, error) {
	if m.getProduct != nil {
		return m.getProduct(ctx, id)
	}
	return nil, nil
}

func (m *MockProductRepo) ListProducts(ctx context.Context, cursor int, limit uint64) ([]*domain.Product, int, bool, error) {
	if m.listProducts != nil {
		return m.listProducts(ctx, cursor, limit)
	}
	return nil, 0, false, nil
}

func (m *MockProductRepo) DeleteProduct(ctx context.Context, id int) error {
	if m.deleteProduct != nil {
		return m.deleteProduct(ctx, id)
	}
	return nil
}

type TestCreate struct {
	name     string
	input    domain.CreateProductRequest
	want     *domain.Product
	wantErr  bool
	wantResp error
}

var testsCreate = []TestCreate{
	{
		name: "general",
		input: domain.CreateProductRequest{
			Name:         "Test-product",
			Manufacturer: "test-manufacturer",
			Price:        10000,
			Amount:       10,
			Status:       true,
			Category:     "Household appliances",
		},
		want: &domain.Product{
			ID:           1,
			Name:         "Test-product",
			Manufacturer: "test-manufacturer",
			Price:        10000,
			Amount:       10,
			Status:       true,
			Category:     "Household appliances",
		},
		wantErr: false,
	},
	{
		name: "missing name",
		input: domain.CreateProductRequest{
			Manufacturer: "test-manufacturer",
			Price:        10000,
			Amount:       10,
			Status:       true,
			Category:     "Household appliances",
		},
		want: &domain.Product{
			ID:           1,
			Name:         "",
			Manufacturer: "test-manufacturer",
			Price:        10000,
			Amount:       10,
			Status:       true,
			Category:     "Household appliances",
		},
		wantErr:  true,
		wantResp: domain.ErrInvalidName,
	},
	{
		name: "missing manufacturer",
		input: domain.CreateProductRequest{
			Name:     "Test-product",
			Price:    10000,
			Amount:   10,
			Status:   true,
			Category: "Household appliances",
		},
		want: &domain.Product{
			ID:       1,
			Name:     "Test-product",
			Price:    10000,
			Amount:   10,
			Status:   true,
			Category: "Household appliances",
		},
		wantErr:  true,
		wantResp: domain.ErrInvalidManufacturer,
	},
	{
		name: "negative price",
		input: domain.CreateProductRequest{
			Name:         "Test-product",
			Manufacturer: "test-manufacturer",
			Price:        -10000,
			Amount:       10,
			Status:       true,
			Category:     "Household appliances",
		},
		want: &domain.Product{
			ID:           1,
			Name:         "Test-product",
			Manufacturer: "test-manufacturer",
			Price:        -10000,
			Amount:       10,
			Status:       true,
			Category:     "Household appliances",
		},
		wantErr:  true,
		wantResp: domain.ErrInvalidPrice,
	},
	{
		name: "null price",
		input: domain.CreateProductRequest{
			Name:         "Test-product",
			Manufacturer: "test-manufacturer",
			Price:        0,
			Amount:       10,
			Status:       true,
			Category:     "Household appliances",
		},
		want: &domain.Product{
			ID:           1,
			Name:         "Test-product",
			Manufacturer: "test-manufacturer",
			Price:        0,
			Amount:       10,
			Status:       true,
			Category:     "Household appliances",
		},
		wantErr:  true,
		wantResp: domain.ErrInvalidPrice,
	},
	{
		name: "negative amount",
		input: domain.CreateProductRequest{
			Name:         "Test-product",
			Manufacturer: "test-manufacturer",
			Price:        10000,
			Amount:       -10,
			Status:       true,
			Category:     "Household appliances",
		},
		want: &domain.Product{
			ID:           1,
			Name:         "Test-product",
			Manufacturer: "test-manufacturer",
			Price:        10000,
			Amount:       -10,
			Status:       true,
			Category:     "Household appliances",
		},
		wantErr:  true,
		wantResp: domain.ErrInvalidAmount,
	},
	{
		name: "missing or not existing category",
		input: domain.CreateProductRequest{
			Name:         "Test-product",
			Manufacturer: "test-manufacturer",
			Price:        10000,
			Amount:       10,
			Status:       true,
		},
		want: &domain.Product{
			ID:           1,
			Name:         "Test-product",
			Manufacturer: "test-manufacturer",
			Price:        10000,
			Amount:       10,
			Status:       true,
			Category:     "",
		},
		wantErr:  true,
		wantResp: domain.ErrInvalidCategory,
	},
}

func TestCreateProduct(t *testing.T) {
	for _, test := range testsCreate {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := MockProductRepo{
				createProduct: func(ctx context.Context, insertData map[string]any) (*domain.Product, error) {
					if test.wantErr {
						return nil, nil
					}

					return test.want, nil
				},
			}

			productService := NewProductService(&mockRepo)

			req := test.input

			product, err := productService.CreateProduct(context.Background(), &req)

			if test.wantErr {
				if product != nil {
					t.Errorf("expected nil, got: %v", product)
				}

				if err != test.wantResp || err == nil {
					t.Fatalf("expected %v, got: %v", test.wantResp, err)
				}
			} else {
				if diff := cmp.Diff(test.want, product); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}

				if err != nil {
					t.Fatalf("expected nil, got: %v", err)
				}
			}
		})
	}
}
