package handler

import (
	"apigateway/services/product/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (s *MockService) CreateProduct(ctx context.Context, req *domain.CreateProductRequest) (*domain.Product, error) {
	args := s.Called(ctx, req)
	p, _ := args.Get(0).(*domain.Product)
	return p, args.Error(1)
}

func (s *MockService) UpdateProduct(ctx context.Context, id int, req *domain.UpdateProductRequest) (*domain.Product, error) {
	args := s.Called(ctx, id, req)
	p, _ := args.Get(0).(*domain.Product)
	return p, args.Error(1)
}

func (s *MockService) GetProduct(ctx context.Context, id int) (*domain.Product, error) {
	args := s.Called(ctx, id)
	p, _ := args.Get(0).(*domain.Product)
	return p, args.Error(1)
}

func (s *MockService) ListProducts(ctx context.Context, cursor int, limit uint64) ([]*domain.Product, *domain.Pagination, error) {
	args := s.Called(ctx, cursor, limit)
	p, _ := args.Get(0).([]*domain.Product)
	pgn, _ := args.Get(1).(*domain.Pagination)
	return p, pgn, args.Error(2)
}

func (s *MockService) DeleteProduct(ctx context.Context, id int) error {
	args := s.Called(ctx, id)
	return args.Error(0)
}

type TestCreate struct {
	name        string
	inputJSON   string
	mockInput   *domain.CreateProductRequest
	mockReturn  *domain.Product
	mockError   error
	resp        string
	skipService bool
	wantErr     bool
	wantStatus  int
	wantResp    error
}

var testsCreate = []TestCreate{
	{
		name:      "success: all values",
		inputJSON: `{"name": "Laptop HUAWEI D16 2024", "manufacturer": "HUAWEI", "price": 57499, "amount": 21, "status": true, "category": "PCs, laptops, peripherals"}`,
		mockInput: &domain.CreateProductRequest{
			Name:         "Laptop HUAWEI D16 2024",
			Manufacturer: "HUAWEI",
			Price:        57499,
			Amount:       21,
			Status:       true,
			Category:     "PCs, laptops, peripherals",
		},
		mockReturn: &domain.Product{
			Name:         "Laptop HUAWEI D16 2024",
			Manufacturer: "HUAWEI",
			Price:        57499,
			Amount:       21,
			Status:       true,
			Category:     "PCs, laptops, peripherals",
		},
		resp:       `{"name": "Laptop HUAWEI D16 2024", "manufacturer": "HUAWEI", "price": 57499, "amount": 21, "status": true, "category": "PCs, laptops, peripherals"}`,
		wantErr:    false,
		wantStatus: 201,
	},
	{
		name:      "success: without status",
		inputJSON: `{"name": "Laptop HUAWEI D16 2024", "manufacturer": "HUAWEI", "price": 57499, "amount": 21, "category": "PCs, laptops, peripherals"}`,
		mockInput: &domain.CreateProductRequest{
			Name:         "Laptop HUAWEI D16 2024",
			Manufacturer: "HUAWEI",
			Price:        57499,
			Amount:       21,
			Category:     "PCs, laptops, peripherals",
		},
		mockReturn: &domain.Product{
			Name:         "Laptop HUAWEI D16 2024",
			Manufacturer: "HUAWEI",
			Price:        57499,
			Amount:       21,
			Status:       false,
			Category:     "PCs, laptops, peripherals",
		},
		resp:       `{"name": "Laptop HUAWEI D16 2024", "manufacturer": "HUAWEI", "price": 57499, "amount": 21, "status": true, "category": "PCs, laptops, peripherals"}`,
		wantErr:    false,
		wantStatus: 201,
	},
	{
		name:      "error: already created",
		inputJSON: `{"name": "Laptop HUAWEI D16 2024", "manufacturer": "HUAWEI", "price": 57499, "amount": 21, "category": "PCs, laptops, peripherals"}`,
		mockInput: &domain.CreateProductRequest{
			Name:         "Laptop HUAWEI D16 2024",
			Manufacturer: "HUAWEI",
			Price:        57499,
			Amount:       21,
			Category:     "PCs, laptops, peripherals",
		},
		mockReturn: nil,
		mockError:  domain.ErrProductExist,
		wantErr:    true,
		wantStatus: 409,
		wantResp:   domain.ErrorResponse{Message: domain.ErrProductExist.Error()},
	},
	{
		name:      "error: no values",
		inputJSON: `{"name": "", "manufacturer": "", "price": 0, "amount": 0, "category": ""}`,
		mockInput: &domain.CreateProductRequest{
			Name:         "",
			Manufacturer: "",
			Price:        0,
			Amount:       0,
			Category:     "",
		},
		mockReturn:  nil,
		skipService: true,
		wantErr:     true,
		wantStatus:  400,
		wantResp: domain.ErrorResponse{
			Message: "validation error",
			Details: map[string]string{"Amount": "this field is required", "Category": "this field is required", "Manufacturer": "this field is required", "Name": "this field is required", "Price": "this field is required"},
		},
	},
	{
		name:        "error: invalid status value",
		inputJSON:   `{"name": "Laptop HUAWEI D16 2024", "manufacturer": "HUAWEI", "price": 57499, "amount": 21, "status": 5123, "category": "PCs, laptops, peripherals"}`,
		mockInput:   nil,
		mockReturn:  nil,
		skipService: true,
		wantErr:     true,
		wantStatus:  400,
		wantResp:    domain.ErrorResponse{Message: domain.ErrInvalidJSON.Error()},
	},
	{
		name:      "error: internal server error",
		inputJSON: `{"name": "Laptop HUAWEI D16 2024", "manufacturer": "HUAWEI", "price": 57499, "amount": 21, "status": true, "category": "PCs, laptops, peripherals"}`,
		mockInput: &domain.CreateProductRequest{
			Name:         "Laptop HUAWEI D16 2024",
			Manufacturer: "HUAWEI",
			Price:        57499,
			Amount:       21,
			Status:       true,
			Category:     "PCs, laptops, peripherals",
		},
		mockReturn: nil,
		mockError:  errors.New("internal server error"),
		wantErr:    true,
		wantStatus: 500,
		wantResp:   domain.ErrorResponse{Message: "internal server error"},
	},
}

func TestProductHandler_Create(t *testing.T) {
	for _, tt := range testsCreate {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if !tt.skipService {
				mockSvc.On("CreateProduct", mock.Anything, tt.mockInput).Return(tt.mockReturn, tt.mockError)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/products", strings.NewReader(tt.inputJSON))
			req.Header.Set("Content-Type", "application/json")

			productHandler := NewProductHandler(mockSvc)

			productHandler.CreateProduct(w, req)

			if tt.wantErr {
				var errResp domain.ErrorResponse
				assert.NoError(t, json.NewDecoder(w.Body).Decode(&errResp))
				assert.Equal(t, tt.wantResp, errResp)
			} else {
				var product *domain.Product
				assert.NoError(t, json.NewDecoder(w.Body).Decode(&product))
				assert.Equal(t, tt.mockReturn, product)
			}
			assert.Equal(t, tt.wantStatus, w.Code)
			mockSvc.AssertExpectations(t)
		})
	}
}

type TestUpdate struct {
	name        string
	idStr       string
	idInt       int
	inputJSON   string
	mockInput   *domain.UpdateProductRequest
	mockReturn  *domain.Product
	mockError   error
	skipService bool
	wantErr     bool
	wantStatus  int
	wantResp    domain.ErrorResponse
}

var testsUpdate = []TestUpdate{
	{
		name:      "success: all values",
		idStr:     "1",
		idInt:     1,
		inputJSON: `{"name": "Laptop HUAWEI D16 2025", "manufacturer": "HUAWEI", "price": 65999, "amount": 25, "status":false, "category": "Laptops"}`,
		mockInput: &domain.UpdateProductRequest{
			Name:         stringPtr("Laptop HUAWEI D16 2025"),
			Manufacturer: stringPtr("HUAWEI"),
			Price:        intPtr(65999),
			Amount:       intPtr(25),
			Status:       boolPtr(false),
			Category:     stringPtr("Laptops"),
		},
		mockReturn: &domain.Product{
			ID:           1,
			Name:         "Laptop HUAWEI D16 2025",
			Manufacturer: "HUAWEI",
			Price:        65999,
			Amount:       25,
			Status:       true,
			Category:     "Laptops",
		},
		wantErr:    false,
		wantStatus: 200,
	},
	{
		name:      "success: no new values",
		idStr:     "1",
		idInt:     1,
		inputJSON: `{"name": "Laptop HUAWEI D16 2024", "manufacturer": "HUAWEI", "price": 57499, "amount": 21, "status": true, "category": "PCs, laptops, peripherals"}`,
		mockInput: &domain.UpdateProductRequest{
			Name:         stringPtr("Laptop HUAWEI D16 2024"),
			Manufacturer: stringPtr("HUAWEI"),
			Price:        intPtr(57499),
			Amount:       intPtr(21),
			Status:       boolPtr(true),
			Category:     stringPtr("PCs, laptops, peripherals"),
		},
		mockReturn: &domain.Product{
			ID:           1,
			Name:         "Laptop HUAWEI D16 2024",
			Manufacturer: "HUAWEI",
			Price:        57499,
			Amount:       21,
			Status:       true,
			Category:     "PCs, laptops, peripherals",
		},
		wantErr:    false,
		wantStatus: 200,
	},
	{
		name:      "error: product does not exist",
		idStr:     "2",
		idInt:     2,
		inputJSON: `{"name": "Laptop HUAWEI D16 2025", "manufacturer": "HUAWEI", "price": 65999, "amount": 21, "category": "PCs, laptops, peripherals"}`,
		mockInput: &domain.UpdateProductRequest{
			Name:         stringPtr("Laptop HUAWEI D16 2025"),
			Manufacturer: stringPtr("HUAWEI"),
			Price:        intPtr(65999),
			Amount:       intPtr(21),
			Category:     stringPtr("PCs, laptops, peripherals"),
		},
		mockReturn: nil,
		mockError:  domain.ErrProductsNotFound,
		wantErr:    true,
		wantStatus: 404,
		wantResp:   domain.ErrorResponse{Message: domain.ErrProductsNotFound.Error()},
	},
	{
		name:       "error: no update data",
		idStr:      "1",
		idInt:      1,
		inputJSON:  `{}`,
		mockInput:  &domain.UpdateProductRequest{},
		mockReturn: nil,
		mockError:  domain.ErrNoUpdateData,
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrNoUpdateData.Error()},
	},
	{
		name:        "error: invalid request",
		idStr:       "1",
		idInt:       1,
		inputJSON:   `1{-&{(}}`,
		mockInput:   nil,
		mockReturn:  nil,
		skipService: true,
		wantErr:     true,
		wantStatus:  400,
		wantResp:    domain.ErrorResponse{Message: domain.ErrInvalidJSON.Error()},
	},
	{
		name:        "error: invalid id",
		idStr:       "-1.5",
		inputJSON:   `{}`,
		mockInput:   &domain.UpdateProductRequest{},
		mockReturn:  nil,
		mockError:   domain.ErrInvalidID,
		skipService: true,
		wantErr:     true,
		wantStatus:  400,
		wantResp:    domain.ErrorResponse{Message: domain.ErrInvalidID.Error()},
	},
	{
		name:        "error: no id provided",
		idStr:       "",
		inputJSON:   `{}`,
		mockInput:   &domain.UpdateProductRequest{},
		mockReturn:  nil,
		mockError:   domain.ErrIDRequired,
		skipService: true,
		wantErr:     true,
		wantStatus:  400,
		wantResp:    domain.ErrorResponse{Message: domain.ErrIDRequired.Error()},
	},
	{
		name:      "error: internal server error",
		idStr:     "1",
		idInt:     1,
		inputJSON: `{"name": "Laptop HUAWEI D16 2025", "manufacturer": "HUAWEI", "price": 65999, "amount": 25, "category": "Laptops"}`,
		mockInput: &domain.UpdateProductRequest{
			Name:         stringPtr("Laptop HUAWEI D16 2025"),
			Manufacturer: stringPtr("HUAWEI"),
			Price:        intPtr(65999),
			Amount:       intPtr(25),
			Category:     stringPtr("Laptops"),
		},
		mockReturn: nil,
		mockError:  errors.New("internal server error"),
		wantErr:    true,
		wantStatus: 500,
		wantResp:   domain.ErrorResponse{Message: "internal server error"},
	},
}

func TestProductHandler_Update(t *testing.T) {
	for _, tt := range testsUpdate {
		t.Run(tt.name, func(t *testing.T) {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.idStr)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/products/{id}", strings.NewReader(tt.inputJSON))
			req.Header.Set("Content-Type", "application/json")

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			mockSvc := new(MockService)
			if !tt.skipService {
				mockSvc.On("UpdateProduct", req.Context(), tt.idInt, tt.mockInput).Return(tt.mockReturn, tt.mockError)
			}

			productHandler := NewProductHandler(mockSvc)

			productHandler.UpdateProduct(w, req)

			if tt.wantErr {
				var errResp domain.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errResp)
				assert.NoError(t, err, "failed to unmarshal w.Body")
				assert.Equal(t, tt.wantResp, errResp)
			} else {
				var product *domain.Product
				err := json.Unmarshal(w.Body.Bytes(), &product)
				assert.NoError(t, err, "failed to unmarshal w.Body")
				assert.Equal(t, tt.mockReturn, product)
			}
			assert.Equal(t, tt.wantStatus, w.Code)
			mockSvc.AssertExpectations(t)
		})
	}
}

func stringPtr(s string) *string { return &s }
func intPtr(i int) *int          { return &i }
func boolPtr(b bool) *bool       { return &b }

type TestGet struct {
	name        string
	idStr       string
	idInt       int
	url         string
	mockReturn  *domain.Product
	mockError   error
	skipService bool
	wantErr     bool
	wantStatus  int
	wantResp    domain.ErrorResponse
}

var testsGet = []TestGet{
	{
		name:  "success",
		idStr: "1",
		idInt: 1,
		url:   "/products/{id}",
		mockReturn: &domain.Product{
			ID:           1,
			Name:         "Laptop HUAWEI D16 2024",
			Manufacturer: "HUAWEI",
			Price:        57499,
			Amount:       21,
			Status:       true,
			Category:     "PCs, laptops, peripherals",
		},
		wantErr:    false,
		wantStatus: 200,
	},
	{
		name:       "error: product not found",
		idStr:      "2",
		idInt:      2,
		url:        "/products/{id}",
		mockReturn: nil,
		mockError:  domain.ErrProductsNotFound,
		wantErr:    true,
		wantStatus: 404,
		wantResp:   domain.ErrorResponse{Message: domain.ErrProductsNotFound.Error()},
	},
	{
		name:        "error: no id provided",
		idStr:       "",
		url:         "/products/{id}",
		mockReturn:  nil,
		mockError:   domain.ErrIDRequired,
		skipService: true,
		wantErr:     true,
		wantStatus:  400,
		wantResp:    domain.ErrorResponse{Message: domain.ErrIDRequired.Error()},
	},
	{
		name:       "error: invalid id",
		idStr:      "-1.5",
		url:        "/products/{id}",
		mockReturn: nil,
		mockError:  domain.ErrInvalidID,
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrInvalidID.Error()},
	},
	{
		name:       "error: internal server error",
		idStr:      "1",
		idInt:      1,
		url:        "/products/{id}",
		mockReturn: nil,
		mockError:  errors.New("internal server error"),
		wantErr:    true,
		wantStatus: 500,
		wantResp:   domain.ErrorResponse{Message: "internal server error"},
	},
}

func TestProductHandler_Get(t *testing.T) {
	for _, tt := range testsGet {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", tt.url, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.idStr)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			mockSvc := new(MockService)
			if !tt.skipService {
				mockSvc.On("GetProduct", req.Context(), tt.idInt).Return(tt.mockReturn, tt.mockError)
			}

			productHandler := NewProductHandler(mockSvc)

			productHandler.GetProduct(w, req)

			if tt.wantErr {
				var errResp domain.ErrorResponse
				err := json.NewDecoder(w.Body).Decode(&errResp)
				assert.NoError(t, err, "failed to decode w.Body")
				assert.Equal(t, tt.wantResp, errResp)
			} else {
				var product *domain.Product
				err := json.NewDecoder(w.Body).Decode(&product)
				assert.NoError(t, err, "failed to decode w.Body")
				assert.Equal(t, tt.mockReturn, product)
			}
			assert.Equal(t, tt.wantStatus, w.Code)
			mockSvc.AssertExpectations(t)
		})
	}
}

type TestList struct {
	name           string
	url            string
	mockCursor     int
	mockLimit      uint64
	mockReturn     []*domain.Product
	mockPagination *domain.Pagination
	mockError      error
	skipService    bool
	wantErr        bool
	wantStatus     int
	wantResp       domain.ErrorResponse
}

var testsList = []TestList{
	{
		name:       "success: 1 page of 3 (cursor=1 limit=2)",
		url:        "/products?cursor=1&limit=2",
		mockCursor: 1,
		mockLimit:  2,
		mockReturn: []*domain.Product{
			{
				ID:           1,
				Name:         "Laptop HUAWEI D16 2024",
				Manufacturer: "HUAWEI",
				Price:        57499,
				Amount:       21,
				Status:       true,
				Category:     "PCs, laptops, peripherals",
			},
			{
				ID:           2,
				Name:         "Microphone Fifine AM8",
				Manufacturer: "Fifine",
				Price:        4499,
				Amount:       18,
				Status:       true,
				Category:     "PC accessories",
			},
		},
		mockPagination: &domain.Pagination{NextCursor: 3, HasMore: true, Limit: 2},
		wantErr:        false,
		wantStatus:     200,
	},
	{
		name:       "success: 2 page of 3 (cursor=3 limit=2)",
		url:        "/products?cursor=3&limit=2",
		mockCursor: 3,
		mockLimit:  2,
		mockReturn: []*domain.Product{
			{
				ID:           3,
				Name:         "Apple iPhone 15 128 GB",
				Manufacturer: "Apple",
				Price:        56999,
				Amount:       23,
				Status:       true,
				Category:     "Smartphones and photographic equipment",
			},
			{
				ID:           4,
				Name:         "TV Samsung UE43U8000FUXRU",
				Manufacturer: "Samsung",
				Price:        30499,
				Amount:       5,
				Status:       true,
				Category:     "TV, consoles, and audio",
			},
		},
		mockPagination: &domain.Pagination{NextCursor: 5, HasMore: true, Limit: 2},
		wantErr:        false,
		wantStatus:     200,
	},
	{
		name:       "success: 3 page of 3 (cursor=5 limit=2)",
		url:        "/products?cursor=5&limit=2",
		mockCursor: 5,
		mockLimit:  2,
		mockReturn: []*domain.Product{
			{
				ID:           5,
				Name:         "LFF Server HDD Toshiba MG09",
				Manufacturer: "Toshiba",
				Price:        68799,
				Amount:       12,
				Status:       true,
				Category:     "Network equipment",
			},
		},
		mockPagination: &domain.Pagination{NextCursor: 0, HasMore: false, Limit: 2},
		wantErr:        false,
		wantStatus:     200,
	},
	{
		name:           "success: no products",
		url:            "/products?cursor=1&limit=1",
		mockCursor:     1,
		mockLimit:      1,
		mockReturn:     []*domain.Product{},
		mockPagination: &domain.Pagination{NextCursor: 0, HasMore: false, Limit: 1},
		wantErr:        false,
		wantStatus:     200,
	},
	{
		name:       "success: only one page (cursor=1 limit=2)",
		url:        "/products?cursor=1&limit=2",
		mockCursor: 1,
		mockLimit:  2,
		mockReturn: []*domain.Product{
			{
				ID:           3,
				Name:         "Apple iPhone 15 128 GB",
				Manufacturer: "Apple",
				Price:        56999,
				Amount:       23,
				Status:       true,
				Category:     "Smartphones and photographic equipment",
			},
		},
		mockPagination: &domain.Pagination{NextCursor: 0, HasMore: false, Limit: 2},
		wantErr:        false,
		wantStatus:     200,
	},
	{
		name:       "success: negative cursor (cursor=-1 limit=1)",
		url:        "/products?cursor=-1&limit=1",
		mockCursor: -1,
		mockLimit:  1,
		mockReturn: []*domain.Product{
			{
				ID:           3,
				Name:         "Apple iPhone 15 128 GB",
				Manufacturer: "Apple",
				Price:        56999,
				Amount:       23,
				Status:       true,
				Category:     "Smartphones and photographic equipment",
			},
		},
		mockPagination: &domain.Pagination{NextCursor: 0, HasMore: false, Limit: 1},
		wantErr:        false,
		wantStatus:     200,
	},
	{
		name:       "success: empty cursor and limit",
		url:        "/products",
		mockCursor: 0,
		mockLimit:  10,
		mockReturn: []*domain.Product{
			{
				ID:           1,
				Name:         "Laptop HUAWEI D16 2024",
				Manufacturer: "HUAWEI",
				Price:        57499,
				Amount:       21,
				Status:       true,
				Category:     "PCs, laptops, peripherals",
			},
			{
				ID:           2,
				Name:         "Microphone Fifine AM8",
				Manufacturer: "Fifine",
				Price:        4499,
				Amount:       18,
				Status:       true,
				Category:     "PC accessories",
			},
			{
				ID:           3,
				Name:         "Apple iPhone 15 128 GB",
				Manufacturer: "Apple",
				Price:        56999,
				Amount:       23,
				Status:       true,
				Category:     "Smartphones and photographic equipment",
			},
			{
				ID:           4,
				Name:         "TV Samsung UE43U8000FUXRU",
				Manufacturer: "Samsung",
				Price:        30499,
				Amount:       5,
				Status:       true,
				Category:     "TV, consoles, and audio",
			},
			{
				ID:           5,
				Name:         "LFF Server HDD Toshiba MG09",
				Manufacturer: "Toshiba",
				Price:        68799,
				Amount:       12,
				Status:       true,
				Category:     "Network equipment",
			},
		},
		mockPagination: &domain.Pagination{NextCursor: 0, HasMore: false, Limit: 10},
		wantErr:        false,
		wantStatus:     200,
	},
	{
		name:        "error: invalid cursor",
		url:         "/products?cursor=a",
		mockReturn:  nil,
		skipService: true,
		wantErr:     true,
		wantStatus:  400,
		wantResp:    domain.ErrorResponse{Message: domain.ErrInvalidCursor.Error()},
	},
	{
		name:        "error: invalid limit",
		url:         "/products?limit=a",
		mockReturn:  nil,
		skipService: true,
		wantErr:     true,
		wantStatus:  400,
		wantResp:    domain.ErrorResponse{Message: domain.ErrInvalidLimit.Error()},
	},
	{
		name:       "error: failed list query",
		url:        "/products",
		mockCursor: 0,
		mockLimit:  10,
		mockReturn: nil,
		mockError:  domain.ErrListQuery,
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrListQuery.Error()},
	},
	{
		name:       "error: internal server error",
		url:        "/products",
		mockCursor: 0,
		mockLimit:  10,
		mockReturn: nil,
		mockError:  errors.New("internal server error"),
		wantErr:    true,
		wantStatus: 500,
		wantResp:   domain.ErrorResponse{Message: "internal server error"},
	},
}

func TestProductHandler_List(t *testing.T) {
	for _, tt := range testsList {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			if !tt.skipService {
				mockSvc.On("ListProducts", mock.Anything, tt.mockCursor, tt.mockLimit).Return(tt.mockReturn, tt.mockPagination, tt.mockError)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", tt.url, nil)

			productHandler := NewProductHandler(mockSvc)

			productHandler.ListProducts(w, req)

			if tt.wantErr {
				var errResp domain.ErrorResponse
				err := json.NewDecoder(w.Body).Decode(&errResp)
				assert.NoError(t, err, "failed to decode w.Body")
				assert.Equal(t, tt.wantResp, errResp)
			} else {
				var listProductsResponse domain.ListProductsResponse
				err := json.NewDecoder(w.Body).Decode(&listProductsResponse)
				assert.NoError(t, err, "failed to decode w.Body")

				listProducts := listProductsResponse.Products
				assert.Equal(t, tt.mockReturn, listProducts)

				pagination := listProductsResponse.PaginationParams
				assert.Equal(t, tt.mockPagination, pagination)
				assert.Equal(t, tt.mockReturn, listProducts)
			}
			assert.Equal(t, tt.wantStatus, w.Code)
			mockSvc.AssertExpectations(t)
		})
	}
}

type TestDelete struct {
	name        string
	url         string
	idStr       string
	idInt       int
	mockError   error
	skipService bool
	wantErr     bool
	wantStatus  int
	wantResp    domain.ErrorResponse
}

var testsDelete = []TestDelete{
	{
		name:       "success",
		url:        "/products/{id}",
		idStr:      "1",
		idInt:      1,
		wantErr:    false,
		wantStatus: 204,
	},
	{
		name:       "error: no product",
		url:        "/products/{id}",
		idStr:      "2",
		idInt:      2,
		mockError:  domain.ErrProductsNotFound,
		wantErr:    true,
		wantStatus: 404,
		wantResp:   domain.ErrorResponse{Message: domain.ErrProductsNotFound.Error()},
	},
	{
		name:        "error: no id provided",
		url:         "/products/{id}",
		idStr:       "",
		skipService: true,
		wantErr:     true,
		wantStatus:  400,
		wantResp:    domain.ErrorResponse{Message: domain.ErrIDRequired.Error()},
	},
	{
		name:        "error: invalid id",
		url:         "/products/{id}",
		idStr:       "a",
		skipService: true,
		wantErr:     true,
		wantStatus:  400,
		wantResp:    domain.ErrorResponse{Message: domain.ErrInvalidID.Error()},
	},
	{
		name:       "error: internal server error",
		url:        "/products/{id}",
		idStr:      "1",
		idInt:      1,
		mockError:  errors.New("internal server error"),
		wantErr:    true,
		wantStatus: 500,
		wantResp:   domain.ErrorResponse{Message: "internal server error"},
	},
}

func TestProductHandler_Delete(t *testing.T) {
	for _, tt := range testsDelete {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", tt.url, nil)
			req = setChiURLParam(req, "id", tt.idStr)

			mockSvc := new(MockService)
			if !tt.skipService {
				mockSvc.On("DeleteProduct", req.Context(), tt.idInt).Return(tt.mockError)
			}

			productHandler := NewProductHandler(mockSvc)
			productHandler.DeleteProduct(w, req)

			if tt.wantErr {
				var errResp domain.ErrorResponse
				err := json.NewDecoder(w.Body).Decode(&errResp)
				assert.NoError(t, err, "failed to decode w.Body")
				assert.Equal(t, tt.wantResp, errResp)
				assert.Equal(t, tt.wantStatus, w.Code)
			}

			assert.Equal(t, tt.wantStatus, w.Code)
			mockSvc.AssertExpectations(t)
		})
	}
}

func setChiURLParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}
