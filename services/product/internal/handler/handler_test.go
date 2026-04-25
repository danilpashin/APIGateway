package handler

import (
	"apigateway/services/product/internal/domain"
	"apigateway/services/product/internal/service"
	"context"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/go-cmp/cmp"
)

type MockProductRepo struct {
	createProduct func(ctx context.Context, insertData map[string]interface{}) (*domain.Product, error)
	updateProduct func(ctx context.Context, id int, updateData map[string]interface{}) (*domain.Product, error)
	getProduct    func(ctx context.Context, id int) (*domain.Product, error)
	listProducts  func(ctx context.Context, cursor int, limit uint64) ([]*domain.Product, int, bool, error)
	deleteProduct func(ctx context.Context, id int) error
}

func (m *MockProductRepo) CreateProduct(ctx context.Context, insertData map[string]interface{}) (*domain.Product, error) {
	if m.createProduct != nil {
		return m.createProduct(ctx, insertData)
	}
	return nil, nil
}

func (m *MockProductRepo) UpdateProduct(ctx context.Context, id int, updateData map[string]interface{}) (*domain.Product, error) {
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
	name       string
	product    *domain.Product
	req        string
	resp       string
	wantErr    bool
	wantStatus int
	wantResp   domain.ErrorResponse
}

var testsCreate = []TestCreate{
	{
		name:       "general",
		product:    &domain.Product{Name: "Laptop HUAWEI D16 2024", Manufacturer: "HUAWEI", Price: 57499, Amount: 21, Status: true, Category: "PCs, laptops, peripherals"},
		req:        `{"name": "Laptop HUAWEI D16 2024", "manufacturer": "HUAWEI", "price": 57499, "amount": 21, "status": true, "category": "PCs, laptops, peripherals"}`,
		resp:       `{"name": "Laptop HUAWEI D16 2024", "manufacturer": "HUAWEI", "price": 57499, "amount": 21, "status": true, "category": "PCs, laptops, peripherals"}`,
		wantErr:    false,
		wantStatus: 201,
	},
	{
		name:       "without status",
		product:    &domain.Product{Name: "Laptop HUAWEI D16 2024", Manufacturer: "HUAWEI", Price: 57499, Amount: 21, Status: true, Category: "PCs, laptops, peripherals"},
		req:        `{"name": "Laptop HUAWEI D16 2024", "manufacturer": "HUAWEI", "price": 57499, "amount": 21, "category": "PCs, laptops, peripherals"}`,
		resp:       `{"name": "Laptop HUAWEI D16 2024", "manufacturer": "HUAWEI", "price": 57499, "amount": 21, "status": true, "category": "PCs, laptops, peripherals"}`,
		wantErr:    false,
		wantStatus: 201,
	},
	{
		name:       "already created",
		product:    &domain.Product{Name: "Laptop HUAWEI D16 2024", Manufacturer: "HUAWEI", Price: 57499, Amount: 21, Status: true, Category: "PCs, laptops, peripherals"},
		req:        `{"name": "Laptop HUAWEI D16 2024", "manufacturer": "HUAWEI", "price": 57499, "amount": 21, "category": "PCs, laptops, peripherals"}`,
		wantErr:    true,
		wantStatus: 409,
		wantResp:   domain.ErrorResponse{Message: domain.ErrProductExist.Error()},
	},
	{
		name:       "no values",
		product:    nil,
		req:        `{"name": "", "manufacturer": "", "price": 0, "amount": 0, "category": ""}`,
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: "validation error", Details: map[string]string{"Amount": "this field is required", "Category": "this field is required", "Manufacturer": "this field is required", "Name": "this field is required", "Price": "this field is required"}},
	},
	{
		name:       "invalid status value",
		product:    nil,
		req:        `{"name": "Laptop HUAWEI D16 2024", "manufacturer": "HUAWEI", "price": 57499, "amount": 21, "status": 5123, "category": "PCs, laptops, peripherals"}`,
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrInvalidJSON.Error()},
	},
	{
		name:       "internal server error",
		product:    nil,
		req:        `{"name": "Laptop HUAWEI D16 2024", "manufacturer": "HUAWEI", "price": 57499, "amount": 21, "status": true, "category": "PCs, laptops, peripherals"}`,
		wantErr:    true,
		wantStatus: 500,
		wantResp:   domain.ErrorResponse{Message: "internal server error"},
	},
}

func TestCreateProductHandler(t *testing.T) {
	for _, test := range testsCreate {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := MockProductRepo{
				createProduct: func(ctx context.Context, insertData map[string]interface{}) (*domain.Product, error) {
					if test.name == "internal server error" {
						return nil, domain.ErrQuery
					}

					if test.wantStatus == 409 {
						return nil, domain.ErrProductExist
					}

					return test.product, nil
				},
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/products", strings.NewReader(test.req))
			req.Header.Set("Content-Type", "application/json")

			productService := service.NewProductService(&mockRepo)
			productHandler := NewProductHandler(*productService)

			productHandler.CreateProductHandler(w, req)

			if test.wantErr {
				var errResp domain.ErrorResponse
				err := json.NewDecoder(w.Body).Decode(&errResp)
				if err != nil {
					t.Fatal("failed to decode w.Body: ", err)
				}

				if !reflect.DeepEqual(errResp, test.wantResp) {
					t.Fatalf("expected %s, got: %s", test.wantResp, errResp)
				}

				if w.Code != test.wantStatus {
					t.Fatalf("expected %d, got: %d", test.wantStatus, w.Code)
				}
			} else {
				var product domain.Product
				err := json.Unmarshal(w.Body.Bytes(), &product)
				if err != nil {
					t.Fatal("failed to unmarshal w.Body: ", err)
				}

				opts := cmp.FilterPath(func(p cmp.Path) bool {
					return p.String() == "CreatedAt" || p.String() == "UpdatedAt"
				}, cmp.Ignore())

				if diff := cmp.Diff(&product, test.product, opts); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}

				if w.Code != test.wantStatus {
					t.Fatalf("expected %d, got: %d", test.wantStatus, w.Code)
				}
			}
		})
	}
}

type TestUpdate struct {
	name       string
	product    *domain.Product
	productID  string
	req        string
	resp       *domain.Product
	wantErr    bool
	wantStatus int
	wantResp   domain.ErrorResponse
}

var testsUpdate = []TestUpdate{
	{
		name:       "general",
		product:    &domain.Product{ID: 1, Name: "Laptop HUAWEI D16 2024", Manufacturer: "HUAWEI", Price: 60499, Amount: 21, Status: true, Category: "PCs, laptops, peripherals"},
		productID:  "1",
		req:        `{"name": "Laptop HUAWEI D16 2025", "manufacturer": "HUAWEI", "price": 65999, "amount": 21, "category": "PCs, laptops, peripherals"}`,
		resp:       &domain.Product{ID: 1, Name: "Laptop HUAWEI D16 2025", Manufacturer: "HUAWEI", Price: 65999, Amount: 21, Status: true, Category: "PCs, laptops, peripherals"},
		wantErr:    false,
		wantStatus: 200,
	},
	{
		name:       "same values",
		product:    &domain.Product{ID: 1, Name: "Laptop HUAWEI D16 2024", Manufacturer: "HUAWEI", Price: 57499, Amount: 21, Status: true, Category: "PCs, laptops, peripherals"},
		productID:  "1",
		req:        `{"name": "Laptop HUAWEI D16 2024", "manufacturer": "HUAWEI", "price": 57499, "amount": 21, "status": true, "category": "PCs, laptops, peripherals"}`,
		resp:       &domain.Product{ID: 1, Name: "Laptop HUAWEI D16 2024", Manufacturer: "HUAWEI", Price: 57499, Amount: 21, Status: true, Category: "PCs, laptops, peripherals"},
		wantErr:    false,
		wantStatus: 200,
	},
	{
		name:       "product does not exist",
		product:    &domain.Product{ID: 1},
		productID:  "2",
		req:        `{"name": "Laptop HUAWEI D16 2025", "manufacturer": "HUAWEI", "price": 65999, "amount": 21, "category": "PCs, laptops, peripherals"}`,
		wantErr:    true,
		wantStatus: 404,
		wantResp:   domain.ErrorResponse{Message: domain.ErrProductsNotFound.Error()},
	},
	{
		name:       "no update data",
		product:    nil,
		productID:  "1",
		req:        `{}`,
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrNoUpdateData.Error()},
	},
	{
		name:       "invalid request",
		product:    nil,
		productID:  "1",
		req:        `1{-&{(}}`,
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrInvalidJSON.Error()},
	},
	{
		name:       "invalid id",
		product:    nil,
		productID:  "-1.5",
		req:        `{}`,
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrInvalidID.Error()},
	},
	{
		name:       "no id provided",
		product:    nil,
		productID:  "",
		req:        `{}`,
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrIDRequired.Error()},
	},
	{
		name:       "internal server error",
		product:    nil,
		productID:  "1",
		req:        `{"name": "Laptop HUAWEI D16 2025", "manufacturer": "HUAWEI", "price": 65999, "amount": 21, "category": "PCs, laptops, peripherals"}`,
		wantErr:    true,
		wantStatus: 500,
		wantResp:   domain.ErrorResponse{Message: "internal server error"},
	},
}

func TestUpdateProductHandler(t *testing.T) {
	for _, test := range testsUpdate {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := MockProductRepo{
				updateProduct: func(ctx context.Context, id int, updateData map[string]interface{}) (*domain.Product, error) {
					if test.name == "internal server error" {
						return nil, domain.ErrQuery
					}

					if id == test.product.ID {
						return test.product, nil
					}

					return nil, domain.ErrProductsNotFound
				},
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/products/{id}", strings.NewReader(test.req))
			req.Header.Set("Content-Type", "application/json")

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", test.productID)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			productService := service.NewProductService(&mockRepo)
			productHandler := NewProductHandler(*productService)

			productHandler.UpdateProductHandler(w, req)

			if test.wantErr {
				var errResp domain.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errResp)
				if err != nil {
					t.Fatal("failed to unmarshal w.Body: ", err)
				}

				if !reflect.DeepEqual(errResp, test.wantResp) {
					t.Fatalf("expected %v, got: %v", test.wantResp, errResp)
				}

				if w.Code != test.wantStatus {
					t.Fatalf("expected %d, got: %d", test.wantStatus, w.Code)
				}
			} else {
				var product domain.Product
				err := json.Unmarshal(w.Body.Bytes(), &product)
				if err != nil {
					t.Fatal("failed to unmarshal w.Body: ", err)
				}

				opts := cmp.FilterPath(func(p cmp.Path) bool {
					return p.String() == "CreatedAt" || p.String() == "UpdatedAt"
				}, cmp.Ignore())

				if diff := cmp.Diff(&product, test.product, opts); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}

				if w.Code != test.wantStatus {
					t.Fatalf("expected %d, got: %d", test.wantStatus, w.Code)
				}
			}
		})
	}
}

type TestGet struct {
	name       string
	product    *domain.Product
	productID  string
	url        string
	wantErr    bool
	wantStatus int
	wantResp   domain.ErrorResponse
}

var testsGet = []TestGet{
	{
		name:       "general",
		product:    &domain.Product{ID: 1, Name: "Laptop HUAWEI D16 2024", Manufacturer: "HUAWEI", Price: 57499, Amount: 21, Status: true, Category: "PCs, laptops, peripherals"},
		productID:  "1",
		url:        "/products/{id}",
		wantErr:    false,
		wantStatus: 200,
	},
	{
		name:       "product not found",
		product:    &domain.Product{ID: 1},
		productID:  "2",
		url:        "/products/{id}",
		wantErr:    true,
		wantStatus: 404,
		wantResp:   domain.ErrorResponse{Message: domain.ErrProductsNotFound.Error()},
	},
	{
		name:       "no id provided",
		product:    nil,
		productID:  "",
		url:        "/products/{id}",
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrIDRequired.Error()},
	},
	{
		name:       "invalid id",
		product:    nil,
		productID:  "-1.5",
		url:        "/products/{id}",
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrInvalidID.Error()},
	},
	{
		name:       "internal server error",
		product:    nil,
		productID:  "1",
		url:        "/products/{id}",
		wantErr:    true,
		wantStatus: 500,
		wantResp:   domain.ErrorResponse{Message: "internal server error"},
	},
}

func TestGetProductHandler(t *testing.T) {
	for _, test := range testsGet {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := MockProductRepo{
				getProduct: func(ctx context.Context, id int) (*domain.Product, error) {
					if test.name == "internal server error" {
						return nil, domain.ErrQuery
					}

					if test.product.ID != id {
						return nil, domain.ErrProductsNotFound
					}

					return test.product, nil
				},
			}

			productService := service.NewProductService(&mockRepo)
			productHandler := NewProductHandler(*productService)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", test.url, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", test.productID)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			productHandler.GetProductHandler(w, req)

			if test.wantErr {
				var errResp domain.ErrorResponse
				err := json.NewDecoder(w.Body).Decode(&errResp)
				if err != nil {
					t.Fatal("failed to decode w.Body: ", err)
				}

				if !reflect.DeepEqual(test.wantResp, errResp) {
					t.Fatalf("expected %v, got: %v", test.wantResp, errResp)
				}

				if w.Code != test.wantStatus {
					t.Fatalf("expected %v, got: %v", test.wantStatus, w.Code)
				}
			} else {
				var product domain.Product
				err := json.NewDecoder(w.Body).Decode(&product)
				if err != nil {
					t.Fatal("failed to decode w.Body: ", err)
				}

				if !reflect.DeepEqual(test.product, &product) {
					t.Fatalf("expected %v, got: %v", test.product, &product)
				}

				if w.Code != test.wantStatus {
					t.Fatalf("expected %d, got: %d", test.wantStatus, w.Code)
				}
			}
		})
	}
}

type TestList struct {
	name           string
	products       []*domain.Product
	url            string
	wantErr        bool
	wantStatus     int
	wantResp       domain.ErrorResponse
	wantPagination *domain.Pagination
}

var testsList = []TestList{
	{
		name: "general 4 products, 1 cursor, 3 limit",
		products: []*domain.Product{
			{ID: 1, Name: "Laptop HUAWEI D16 2024", Manufacturer: "HUAWEI", Price: 57499, Amount: 21, Status: true, Category: "PCs, laptops, peripherals"},
			{ID: 2, Name: "Microphone Fifine AM8", Manufacturer: "Fifine", Price: 4499, Amount: 18, Status: true, Category: "PC accessories"},
			{ID: 3, Name: "Apple iPhone 15 128 GB", Manufacturer: "Apple", Price: 56999, Amount: 23, Status: true, Category: "Smartphones and photographic equipment"},
			{ID: 4, Name: "TV Samsung UE43U8000FUXRU", Manufacturer: "Samsung", Price: 30499, Amount: 5, Status: true, Category: "TV, consoles, and audio"},
		},
		url:            "/products?cursor=1&limit=3",
		wantErr:        false,
		wantStatus:     200,
		wantPagination: &domain.Pagination{NextCursor: 4, HasMore: true, Limit: 3},
	},
	{
		name: "general 3 products, 2 limit",
		products: []*domain.Product{
			{ID: 1, Name: "Laptop HUAWEI D16 2024", Manufacturer: "HUAWEI", Price: 57499, Amount: 21, Status: true, Category: "PCs, laptops, peripherals"},
			{ID: 2, Name: "Microphone Fifine AM8", Manufacturer: "Fifine", Price: 4499, Amount: 18, Status: true, Category: "PC accessories"},
			{ID: 3, Name: "Apple iPhone 15 128 GB", Manufacturer: "Apple", Price: 56999, Amount: 23, Status: true, Category: "Smartphones and photographic equipment"},
		},
		url:            "/products?cursor=1&limit=2",
		wantErr:        false,
		wantStatus:     200,
		wantPagination: &domain.Pagination{NextCursor: 3, HasMore: true, Limit: 2},
	},
	{
		name: "general 1 product(ID=3), 3 cursor, 2 limit",
		products: []*domain.Product{
			{ID: 3, Name: "Apple iPhone 15 128 GB", Manufacturer: "Apple", Price: 56999, Amount: 23, Status: true, Category: "Smartphones and photographic equipment"},
		},
		url:            "/products?cursor=3&limit=2",
		wantErr:        false,
		wantStatus:     200,
		wantPagination: &domain.Pagination{NextCursor: 0, HasMore: false, Limit: 2},
	},
	{
		name:           "general no products",
		products:       []*domain.Product{},
		url:            "/products?cursor=1&limit=1",
		wantErr:        false,
		wantStatus:     200,
		wantPagination: &domain.Pagination{NextCursor: 0, HasMore: false, Limit: 1},
	},
	{
		name: "0 cursor, 1 product",
		products: []*domain.Product{
			{ID: 3, Name: "Apple iPhone 15 128 GB", Manufacturer: "Apple", Price: 56999, Amount: 23, Status: true, Category: "Smartphones and photographic equipment"},
		},
		url:            "/products?cursor=0&limit=1",
		wantErr:        false,
		wantStatus:     200,
		wantPagination: &domain.Pagination{NextCursor: 0, HasMore: false, Limit: 1},
	},
	{
		name: "-1 cursor, 1 product",
		products: []*domain.Product{
			{ID: 3, Name: "Apple iPhone 15 128 GB", Manufacturer: "Apple", Price: 56999, Amount: 23, Status: true, Category: "Smartphones and photographic equipment"},
		},
		url:            "/products?cursor=-1&limit=1",
		wantErr:        false,
		wantStatus:     200,
		wantPagination: &domain.Pagination{NextCursor: 0, HasMore: false, Limit: 1},
	},
	{
		name: "empty cursor and limit",
		products: []*domain.Product{
			{ID: 3, Name: "Apple iPhone 15 128 GB", Manufacturer: "Apple", Price: 56999, Amount: 23, Status: true, Category: "Smartphones and photographic equipment"},
		},
		url:            "/products",
		wantErr:        false,
		wantStatus:     200,
		wantPagination: &domain.Pagination{NextCursor: 0, HasMore: false, Limit: 10},
	},
	{
		name:       "invalid cursor",
		products:   []*domain.Product{},
		url:        "/products?cursor=a",
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrInvalidCursor.Error()},
	},
	{
		name:       "invalid limit",
		products:   []*domain.Product{},
		url:        "/products?limit=a",
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrInvalidLimit.Error()},
	},
	{
		name:       "error list query",
		products:   []*domain.Product{},
		url:        "/products",
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrListQuery.Error()},
	},
	{
		name:       "internal server error",
		products:   nil,
		url:        "/products",
		wantErr:    true,
		wantStatus: 500,
		wantResp:   domain.ErrorResponse{Message: "internal server error"},
	},
}

func TestListProductsHandler(t *testing.T) {
	for _, test := range testsList {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := MockProductRepo{
				listProducts: func(ctx context.Context, cursor int, limit uint64) ([]*domain.Product, int, bool, error) {
					if test.name == "internal server error" {
						return nil, 0, false, domain.ErrQuery
					}

					if test.name == "error list query" {
						return nil, 0, false, domain.ErrListQuery
					}

					listProducts := make([]*domain.Product, 0, limit+1)

					currSize := 0
					for _, product := range test.products {
						if product.ID >= cursor {
							listProducts = append(listProducts, product)
							currSize++
						}
						if currSize == int(limit+1) {
							break
						}
					}

					var products []*domain.Product
					var nextCursor int
					var hasMore bool

					if len(listProducts) > int(limit) {
						nextCursor = listProducts[len(listProducts)-1].ID
						hasMore = true
						products = listProducts[:limit]
					} else {
						nextCursor = 0
						hasMore = false
						products = listProducts
					}

					return products, nextCursor, hasMore, nil
				},
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", test.url, nil)

			productService := service.NewProductService(&mockRepo)
			productHandler := NewProductHandler(*productService)

			productHandler.ListProductsHandler(w, req)

			if test.wantErr {
				var errResp domain.ErrorResponse
				err := json.NewDecoder(w.Body).Decode(&errResp)
				if err != nil {
					t.Fatal("failed to decode w.Body: ", err)
				}

				if !reflect.DeepEqual(test.wantResp, errResp) {
					t.Fatalf("expected %v, got: %v", test.wantResp, errResp)
				}

				if w.Code != test.wantStatus {
					t.Fatalf("expected %d, got: %d", test.wantStatus, w.Code)
				}
			} else {
				var listProductsResponse domain.ListProductsResponse
				err := json.NewDecoder(w.Body).Decode(&listProductsResponse)
				if err != nil {
					t.Fatal("failed to decode w.Body: ", err)
				}

				listProducts := listProductsResponse.Products

				pagination := listProductsResponse.PaginationParams
				if !reflect.DeepEqual(test.wantPagination, pagination) {
					t.Fatalf("expected %v, got: %v", test.wantPagination, pagination)
				}

				limit := pagination.Limit

				if len(listProducts) > int(limit) {
					t.Fatalf("expected at most %d products, got: %d", limit, len(test.products))
				}

				expectedLen := int(limit)
				if len(listProducts) < expectedLen {
					expectedLen = len(listProducts)
				}

				if !reflect.DeepEqual(test.products[:expectedLen], listProducts) {
					t.Fatalf("expected %v, got: %v", test.products, listProducts)
				}

				if w.Code != test.wantStatus {
					t.Fatalf("expected %d, got: %d", test.wantStatus, w.Code)
				}
			}
		})
	}
}

type TestDelete struct {
	name       string
	url        string
	productID  string
	wantErr    bool
	wantStatus int
	wantResp   domain.ErrorResponse
}

var testsDelete = []TestDelete{
	{
		name:       "general",
		url:        "/products/{id}",
		productID:  "1",
		wantErr:    false,
		wantStatus: 204,
	},
	{
		name:       "no product",
		url:        "/products/{id}",
		productID:  "2",
		wantErr:    true,
		wantStatus: 404,
		wantResp:   domain.ErrorResponse{Message: domain.ErrProductsNotFound.Error()},
	},
	{
		name:       "no id provided",
		url:        "/products/{id}",
		productID:  "",
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrIDRequired.Error()},
	},
	{
		name:       "invalid id",
		url:        "/products/{id}",
		productID:  "a",
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrInvalidID.Error()},
	},
	{
		name:       "internal server error",
		url:        "/products/{id}",
		productID:  "1",
		wantErr:    true,
		wantStatus: 500,
		wantResp:   domain.ErrorResponse{Message: "internal server error"},
	},
}

func TestDeleteProductHandler(t *testing.T) {
	for _, test := range testsDelete {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := MockProductRepo{
				deleteProduct: func(ctx context.Context, id int) error {
					if test.name == "internal server error" {
						return domain.ErrQuery
					}

					if test.name == "no product" {
						return domain.ErrProductsNotFound
					}
					return nil
				},
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", test.url, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", test.productID)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			productService := service.NewProductService(&mockRepo)
			productHandler := NewProductHandler(*productService)

			productHandler.DeleteProductHandler(w, req)

			if test.wantErr {
				var errResp domain.ErrorResponse
				err := json.NewDecoder(w.Body).Decode(&errResp)
				if err != nil {
					t.Fatal("failed to decode w.Body: ", err)
				}

				if !reflect.DeepEqual(errResp, test.wantResp) {
					t.Fatalf("expected %v, got: %v", test.wantResp, errResp)
				}

				if w.Code != test.wantStatus {
					t.Fatalf("expected %d, got: %d", test.wantStatus, w.Code)
				}
			} else {
				if w.Code != test.wantStatus {
					t.Fatalf("expected %d, got: %d", test.wantStatus, w.Code)
				}
			}
		})
	}
}
