package service

import (
	"apigateway/services/product/internal/domain"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) CreateProduct(ctx context.Context, insertData map[string]any) (*domain.Product, error) {
	args := m.Called(ctx, insertData)
	p, _ := args[0].(*domain.Product)
	return p, args.Error(1)
}

func (m *MockRepo) UpdateProduct(ctx context.Context, id int, updateData map[string]any) (*domain.Product, error) {
	args := m.Called(ctx, id, updateData)
	p, _ := args[0].(*domain.Product)
	return p, args.Error(1)
}

func (m *MockRepo) GetProduct(ctx context.Context, id int) (*domain.Product, error) {
	args := m.Called(ctx, id)
	p, _ := args[0].(*domain.Product)
	return p, args.Error(1)
}

func (m *MockRepo) ListProducts(ctx context.Context, cursor int, limit uint64) ([]*domain.Product, int, bool, error) {
	args := m.Called(ctx, cursor, limit)
	p, _ := args[0].([]*domain.Product)
	return p, args.Int(1), args.Bool(2), args.Error(3)
}

func (m *MockRepo) DeleteProduct(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type TestCreate struct {
	name       string
	input      domain.CreateProductRequest
	mockInput  map[string]any
	mockReturn *domain.Product
	mockError  error
	want       *domain.Product
	wantErr    bool
	wantResp   error
}

var testsCreate = []TestCreate{
	{
		name: "success: all values",
		input: domain.CreateProductRequest{
			Name:         "Test-product",
			Manufacturer: "test-manufacturer",
			Price:        10000,
			Amount:       10,
			Status:       true,
			Category:     "Household appliances",
		},
		mockInput: map[string]any{
			"name":         "Test-product",
			"manufacturer": "test-manufacturer",
			"price":        10000,
			"amount":       10,
			"status":       true,
			"category":     "Household appliances",
		},
		mockReturn: &domain.Product{
			ID:           1,
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
		name:     "error: missing required field: name",
		input:    domain.CreateProductRequest{},
		wantErr:  true,
		wantResp: domain.ErrNameRequired,
	},
	{
		name: "error: invalid name format",
		input: domain.CreateProductRequest{
			Name: "t",
		},
		wantErr:  true,
		wantResp: domain.ErrInvalidName,
	},
	{
		name: "error: missing required field: manufacturer",
		input: domain.CreateProductRequest{
			Name: "Test-product",
		},
		wantErr:  true,
		wantResp: domain.ErrManufacturerRequired,
	},
	{
		name: "error: invalid manufacturer format",
		input: domain.CreateProductRequest{
			Name:         "Test-product",
			Manufacturer: "t",
		},
		wantErr:  true,
		wantResp: domain.ErrInvalidManufacturer,
	},
	{
		name: "error: negative price",
		input: domain.CreateProductRequest{
			Name:         "Test-product",
			Manufacturer: "test-manufacturer",
			Price:        -10000,
		},
		wantErr:  true,
		wantResp: domain.ErrInvalidPrice,
	},
	{
		name: "error: null price",
		input: domain.CreateProductRequest{
			Name:         "Test-product",
			Manufacturer: "test-manufacturer",
			Price:        0,
		},
		wantErr:  true,
		wantResp: domain.ErrInvalidPrice,
	},
	{
		name: "error: negative amount",
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
		name: "error: missing required field: category",
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
		name: "error: invalid category format",
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

func TestProductService_Create(t *testing.T) {
	for _, tt := range testsCreate {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			if !tt.wantErr {
				mockRepo.On("CreateProduct", mock.Anything, tt.mockInput).Return(tt.mockReturn, tt.wantResp)
			}

			productService := NewProductService(mockRepo)

			product, err := productService.CreateProduct(context.Background(), &tt.input)

			if tt.wantErr {
				assert.Equal(t, tt.want, product)
				assert.Equal(t, tt.wantResp, err)
			} else {
				assert.Equal(t, tt.want, product)
				assert.Equal(t, nil, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

type TestUpdate struct {
	name       string
	input      domain.UpdateProductRequest
	productID  int
	mockInput  map[string]any
	mockReturn *domain.Product
	want       *domain.Product
	wantErr    bool
	wantResp   error
}

var testsUpdate = []TestUpdate{
	{
		name: "success",
		input: domain.UpdateProductRequest{
			Name:         stringPtr("UPD-Test-product"),
			Manufacturer: stringPtr("UPD-Test-manufacturer"),
			Price:        intPtr(15000),
			Amount:       intPtr(15),
			Status:       boolPtr(true),
			Category:     stringPtr("Household appliances"),
		},
		productID: 1,
		mockInput: map[string]any{
			"name":         "UPD-Test-product",
			"manufacturer": "UPD-Test-manufacturer",
			"price":        15000,
			"amount":       15,
			"status":       true,
			"category":     "Household appliances",
		},
		mockReturn: &domain.Product{
			ID:           1,
			Name:         "UPD-Test-product",
			Manufacturer: "UPD-Test-manufacturer",
			Price:        15000,
			Amount:       15,
			Status:       true,
			Category:     "Household appliances",
		},
		want: &domain.Product{
			ID:           1,
			Name:         "UPD-Test-product",
			Manufacturer: "UPD-Test-manufacturer",
			Price:        15000,
			Amount:       15,
			Status:       true,
			Category:     "Household appliances",
		},
		wantErr: false,
	},
	{
		name: "error: missing all update values",
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
		name: "error: invalid name format",
		input: domain.UpdateProductRequest{
			Name: stringPtr("t"),
		},
		productID: 1,
		wantErr:   true,
		wantResp:  domain.ErrInvalidName,
	},
	{
		name: "error: invalid manufacturer format",
		input: domain.UpdateProductRequest{
			Name:         stringPtr("UPD-Test-product"),
			Manufacturer: stringPtr(""),
		},
		productID: 1,
		wantErr:   true,
		wantResp:  domain.ErrInvalidManufacturer,
	},
	{
		name: "error: invalid price format",
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
		name: "error: invalid amount format",
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
		name: "error: invalid category format",
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

func TestProductService_Update(t *testing.T) {
	for _, tt := range testsUpdate {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			if !tt.wantErr {
				mockRepo.On("UpdateProduct", mock.Anything, tt.productID, tt.mockInput).Return(tt.mockReturn, tt.wantResp)
			}

			productService := NewProductService(mockRepo)

			product, err := productService.UpdateProduct(context.Background(), tt.productID, &tt.input)

			if tt.wantErr {
				assert.Equal(t, tt.want, product)
				assert.Equal(t, tt.wantResp, err)
			} else {
				assert.Equal(t, tt.want, product)
				assert.Equal(t, nil, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func stringPtr(s string) *string { return &s }
func intPtr(i int) *int          { return &i }
func boolPtr(b bool) *bool       { return &b }

type TestGet struct {
	name       string
	productID  int
	mockReturn *domain.Product
	want       *domain.Product
	wantErr    bool
	wantResp   error
}

var testsGet = []TestGet{
	{
		name:      "success",
		productID: 1,
		mockReturn: &domain.Product{
			Name:         "Test-product",
			Manufacturer: "test-manufacturer",
			Price:        10000,
			Amount:       10,
			Status:       true,
			Category:     "Household appliances",
		},
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
		name:      "error: invalid ID format",
		productID: -1,
		wantErr:   true,
		wantResp:  domain.ErrInvalidID,
	},
}

func TestProductService_Get(t *testing.T) {
	for _, tt := range testsGet {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			if !tt.wantErr {
				mockRepo.On("GetProduct", mock.Anything, tt.productID).Return(tt.mockReturn, tt.wantResp)
			}

			productService := NewProductService(mockRepo)

			product, err := productService.GetProduct(context.Background(), tt.productID)

			if tt.wantErr {
				assert.Equal(t, tt.want, product)
				assert.Equal(t, tt.wantResp, err)
			} else {
				assert.Equal(t, tt.want, product)
				assert.Equal(t, nil, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

type TestList struct {
	name           string
	cursor         int
	limit          uint64
	mockNextCursor int
	mockHasMore    bool
	mockCursor     int
	mockLimit      uint64
	want           []*domain.Product
	wantPagination *domain.Pagination
	wantErr        bool
	wantResp       error
}

var testsList = []TestList{
	{
		name:           "success: first two products (cursor=1, limit=2)",
		cursor:         1,
		limit:          2,
		mockNextCursor: 3,
		mockHasMore:    true,
		mockCursor:     1,
		mockLimit:      2,
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
		name:       "success: end of list (no more)",
		cursor:     3,
		limit:      3,
		mockCursor: 3,
		mockLimit:  3,
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
		name:       "success: negative cursor and null limit",
		cursor:     -2,
		limit:      0,
		mockCursor: 0,
		mockLimit:  10,
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
		name:       "success: limit > 50 clamped to 50",
		cursor:     1,
		limit:      100,
		mockCursor: 1,
		mockLimit:  50,
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
		name:       "error: repository error",
		cursor:     1,
		limit:      2,
		mockCursor: 1,
		mockLimit:  2,
		wantErr:    true,
		wantResp:   errors.New("repo error"),
	},
}

func TestProductService_List(t *testing.T) {
	for _, tt := range testsList {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			mockRepo.On("ListProducts", mock.Anything, tt.mockCursor, tt.mockLimit).Return(tt.want, tt.mockNextCursor, tt.mockHasMore, tt.wantResp)

			productService := NewProductService(mockRepo)

			products, pagination, err := productService.ListProducts(context.Background(), tt.cursor, tt.limit)

			if tt.wantErr {
				assert.Equal(t, tt.want, products)
				assert.Equal(t, tt.wantResp, err)
			} else {
				assert.Equal(t, tt.want, products)
				assert.Equal(t, tt.wantPagination, pagination)
				assert.Equal(t, nil, err)
			}
			mockRepo.AssertExpectations(t)
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
		name:      "success",
		productID: 1,
		wantErr:   false,
	},
	{
		name:      "error: invalid ID",
		productID: -1,
		wantErr:   true,
		wantResp:  domain.ErrInvalidID,
	},
}

func TestProductService_Delete(t *testing.T) {
	for _, tt := range testsDelete {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			if !tt.wantErr {
				mockRepo.On("DeleteProduct", mock.Anything, tt.productID).Return(tt.wantResp)
			}

			productService := NewProductService(mockRepo)

			err := productService.DeleteProduct(context.Background(), tt.productID)

			if tt.wantErr {
				assert.Equal(t, tt.wantResp, err)
			} else {
				assert.Equal(t, nil, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
