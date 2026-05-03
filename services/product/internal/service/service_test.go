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
		name:     "missing name",
		input:    domain.CreateProductRequest{},
		wantErr:  true,
		wantResp: domain.ErrNameRequired,
	},
	{
		name: "invalid name",
		input: domain.CreateProductRequest{
			Name: "t",
		},
		wantErr:  true,
		wantResp: domain.ErrInvalidName,
	},
	{
		name: "missing manufacturer",
		input: domain.CreateProductRequest{
			Name: "Test-product",
		},
		wantErr:  true,
		wantResp: domain.ErrManufacturerRequired,
	},
	{
		name: "invalid manufacturer",
		input: domain.CreateProductRequest{
			Name:         "Test-product",
			Manufacturer: "t",
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
		},
		wantErr:  true,
		wantResp: domain.ErrInvalidAmount,
	},
	{
		name: "missing category",
		input: domain.CreateProductRequest{
			Name:         "Test-product",
			Manufacturer: "test-manufacturer",
			Price:        10000,
			Amount:       10,
			Status:       true,
		},
		wantErr:  true,
		wantResp: domain.ErrCategoryRequired,
	},
	{
		name: "invalid category",
		input: domain.CreateProductRequest{
			Name:         "Test-product",
			Manufacturer: "test-manufacturer",
			Price:        10000,
			Amount:       10,
			Status:       true,
			Category:     "h",
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

			product, err := productService.CreateProduct(context.Background(), &test.input)

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

type TestUpdate struct {
	name      string
	input     domain.UpdateProductRequest
	productID int
	want      *domain.Product
	wantErr   bool
	wantResp  error
}

var testsUpdate = []TestUpdate{
	{
		name: "general",
		input: domain.UpdateProductRequest{
			Name:         stringPtr("UPD-Test-product"),
			Manufacturer: stringPtr("UPD-Test-manufacturer"),
			Price:        intPtr(15000),
			Amount:       intPtr(15),
			Status:       boolPtr(true),
			Category:     stringPtr("Household appliances"),
		},
		productID: 1,
		want: &domain.Product{
			ID:           1,
			Name:         "UPD-Test-product",
			Manufacturer: "UPD-Test-manufacturer",
			Price:        15000,
			Amount:       15,
			Status:       false,
			Category:     "Household appliances",
		},
		wantErr: false,
	},
	{
		name: "missing values",
		input: domain.UpdateProductRequest{
			Name:         nil,
			Manufacturer: nil,
			Price:        nil,
			Amount:       nil,
			Status:       nil,
			Category:     nil,
		},
		productID: 1,
		wantErr:   true,
		wantResp:  domain.ErrNoUpdateData,
	},
	{
		name: "invalid name",
		input: domain.UpdateProductRequest{
			Name: stringPtr("t"),
		},
		productID: 1,
		wantErr:   true,
		wantResp:  domain.ErrInvalidName,
	},
	{
		name: "invalid manufacturer",
		input: domain.UpdateProductRequest{
			Name:         stringPtr("UPD-Test-product"),
			Manufacturer: stringPtr(""),
		},
		productID: 1,
		wantErr:   true,
		wantResp:  domain.ErrInvalidManufacturer,
	},
	{
		name: "invalid price",
		input: domain.UpdateProductRequest{
			Name:         stringPtr("UPD-Test-product"),
			Manufacturer: stringPtr("UPD-Test-manufacturer"),
			Price:        intPtr(-15000),
		},
		productID: 1,
		wantErr:   true,
		wantResp:  domain.ErrInvalidPrice,
	},
	{
		name: "invalid amount",
		input: domain.UpdateProductRequest{
			Name:         stringPtr("UPD-Test-product"),
			Manufacturer: stringPtr("UPD-Test-manufacturer"),
			Price:        intPtr(15000),
			Amount:       intPtr(-15),
		},
		productID: 1,
		wantErr:   true,
		wantResp:  domain.ErrInvalidAmount,
	},
	{
		name: "invalid category",
		input: domain.UpdateProductRequest{
			Name:         stringPtr("UPD-Test-product"),
			Manufacturer: stringPtr("UPD-Test-manufacturer"),
			Price:        intPtr(15000),
			Amount:       intPtr(15),
			Category:     stringPtr(""),
		},
		productID: 1,
		wantErr:   true,
		wantResp:  domain.ErrInvalidCategory,
	},
}

func TestUpdateProduct(t *testing.T) {
	for _, test := range testsUpdate {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := MockProductRepo{
				updateProduct: func(ctx context.Context, id int, updateData map[string]any) (*domain.Product, error) {
					if test.wantErr {
						return nil, nil
					}

					return test.want, nil
				},
			}

			productService := NewProductService(&mockRepo)

			product, err := productService.UpdateProduct(context.Background(), test.productID, &test.input)

			if test.wantErr {
				if product != nil {
					t.Errorf("expected nil, got: %v", product)
				}

				if err != test.wantResp {
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

func stringPtr(s string) *string { return &s }
func intPtr(i int) *int          { return &i }
func boolPtr(b bool) *bool       { return &b }

type TestGet struct {
	name      string
	productID int
	want      *domain.Product
	wantErr   bool
	wantResp  error
}

var testsGet = []TestGet{
	{
		name:      "general",
		productID: 1,
		want: &domain.Product{
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
		name:      "invalid ID",
		productID: -1,
		wantErr:   true,
		wantResp:  domain.ErrInvalidID,
	},
	{
		name:      "product not found",
		productID: 1,
		wantErr:   true,
		wantResp:  domain.ErrProductsNotFound,
	},
}

func TestGetProduct(t *testing.T) {
	for _, test := range testsGet {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := MockProductRepo{
				getProduct: func(ctx context.Context, id int) (*domain.Product, error) {
					if test.wantErr {
						return nil, domain.ErrProductsNotFound
					}

					return test.want, nil
				},
			}

			productService := NewProductService(&mockRepo)

			product, err := productService.GetProduct(context.Background(), test.productID)

			if test.wantErr {
				if product != nil {
					t.Errorf("expected nil, got: %v", product)
				}

				if err != test.wantResp {
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

type TestList struct {
	name           string
	cursor         int
	limit          uint64
	want           []*domain.Product
	wantPagination *domain.Pagination
	wantErr        bool
	wantResp       error
}

var testsList = []TestList{
	{
		name:   "general",
		cursor: 1,
		limit:  2,
		want: []*domain.Product{
			{
				ID:           1,
				Name:         "test-product",
				Manufacturer: "test-manufacturer",
				Price:        10000,
				Amount:       10,
				Status:       true,
				Category:     "Household appliances",
			},
			{
				ID:           2,
				Name:         "test-product",
				Manufacturer: "test-manufacturer",
				Price:        10000,
				Amount:       10,
				Status:       true,
				Category:     "Household appliances",
			},
			{
				ID:           3,
				Name:         "test-product",
				Manufacturer: "test-manufacturer",
				Price:        10000,
				Amount:       10,
				Status:       true,
				Category:     "Household appliances",
			},
		},
		wantPagination: &domain.Pagination{
			NextCursor: 3,
			HasMore:    true,
			Limit:      2,
		},
		wantErr: false,
	},
	{
		name:   "list end",
		cursor: 2,
		limit:  3,
		want: []*domain.Product{
			{
				ID:           1,
				Name:         "test-product",
				Manufacturer: "test-manufacturer",
				Price:        10000,
				Amount:       10,
				Status:       true,
				Category:     "Household appliances",
			},
			{
				ID:           2,
				Name:         "test-product",
				Manufacturer: "test-manufacturer",
				Price:        10000,
				Amount:       10,
				Status:       true,
				Category:     "Household appliances",
			},
			{
				ID:           3,
				Name:         "test-product",
				Manufacturer: "test-manufacturer",
				Price:        10000,
				Amount:       10,
				Status:       true,
				Category:     "Household appliances",
			},
		},
		wantPagination: &domain.Pagination{
			NextCursor: 0,
			HasMore:    false,
			Limit:      3,
		},
		wantErr: false,
	},
	{
		name:   "negative cursor and null limit",
		cursor: -2,
		limit:  0,
		want: []*domain.Product{
			{
				ID:           1,
				Name:         "test-product",
				Manufacturer: "test-manufacturer",
				Price:        10000,
				Amount:       10,
				Status:       true,
				Category:     "Household appliances",
			},
			{
				ID:           2,
				Name:         "test-product",
				Manufacturer: "test-manufacturer",
				Price:        10000,
				Amount:       10,
				Status:       true,
				Category:     "Household appliances",
			},
			{
				ID:           3,
				Name:         "test-product",
				Manufacturer: "test-manufacturer",
				Price:        10000,
				Amount:       10,
				Status:       true,
				Category:     "Household appliances",
			},
		},
		wantPagination: &domain.Pagination{
			NextCursor: 0,
			HasMore:    false,
			Limit:      10,
		},
		wantErr: false,
	},
	{
		name:   "50+ limit",
		cursor: 1,
		limit:  100,
		want: []*domain.Product{
			{
				ID:           1,
				Name:         "test-product",
				Manufacturer: "test-manufacturer",
				Price:        10000,
				Amount:       10,
				Status:       true,
				Category:     "Household appliances",
			},
			{
				ID:           2,
				Name:         "test-product",
				Manufacturer: "test-manufacturer",
				Price:        10000,
				Amount:       10,
				Status:       true,
				Category:     "Household appliances",
			},
			{
				ID:           3,
				Name:         "test-product",
				Manufacturer: "test-manufacturer",
				Price:        10000,
				Amount:       10,
				Status:       true,
				Category:     "Household appliances",
			},
		},
		wantPagination: &domain.Pagination{
			NextCursor: 0,
			HasMore:    false,
			Limit:      50,
		},
		wantErr: false,
	},
	{
		name:    "repo error",
		cursor:  1,
		limit:   2,
		wantErr: true,
	},
}

func TestListProducts(t *testing.T) {
	for _, test := range testsList {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := MockProductRepo{
				listProducts: func(ctx context.Context, cursor int, limit uint64) ([]*domain.Product, int, bool, error) {
					if test.wantErr {
						return nil, 0, false, domain.ErrListQuery
					}

					return test.want, test.wantPagination.NextCursor, test.wantPagination.HasMore, nil
				},
			}

			productService := NewProductService(&mockRepo)

			products, pagination, err := productService.ListProducts(context.Background(), test.cursor, test.limit)

			if test.wantErr {
				if products != nil {
					t.Errorf("expected nil, got: %v", products)
				}

				if err == nil {
					t.Fatalf("expected nil, got: %v", err)
				}
			} else {
				if diff := cmp.Diff(test.want, products); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}

				if diff := cmp.Diff(test.wantPagination, pagination); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}

				if err != nil {
					t.Fatalf("expected nil, got: %v", err)
				}
			}
		})
	}
}

type TestDelete struct {
	name      string
	productID int
	wantErr   bool
	wantResp  error
}

var testsDelete = []TestDelete{
	{
		name:      "general",
		productID: 1,
		wantErr:   false,
	},
	{
		name:      "invalid ID",
		productID: -1,
		wantErr:   true,
		wantResp:  domain.ErrInvalidID,
	},
}

func TestDeleteProduct(t *testing.T) {
	for _, test := range testsDelete {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := MockProductRepo{
				deleteProduct: func(ctx context.Context, id int) error {
					return nil
				},
			}

			productService := NewProductService(&mockRepo)

			err := productService.DeleteProduct(context.Background(), test.productID)

			if test.wantErr {
				if err != test.wantResp {
					t.Fatalf("expected %v, got: %v", test.wantResp, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected nil, got: %v", err)
				}
			}
		})
	}
}
