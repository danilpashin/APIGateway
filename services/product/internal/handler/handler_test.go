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

	"github.com/google/go-cmp/cmp"
)

type MockProductRepo struct {
	createProduct func(ctx context.Context, insertData map[string]interface{}) (*domain.Product, error)
	updateProduct func(ctx context.Context, id int, updateData map[string]interface{}) (*domain.Product, error)
	getProduct    func(ctx context.Context, id int) (*domain.Product, error)
	listProducts  func(ctx context.Context, cursor int, limit uint64) ([]*domain.Product, error)
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

func (m *MockProductRepo) ListProducts(ctx context.Context, cursor int, limit uint64) ([]*domain.Product, error) {
	if m.listProducts != nil {
		return m.listProducts(ctx, cursor, limit)
	}
	return nil, nil
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
		wantResp:   domain.ErrorResponse{Error: "product already exists"},
	},
	{
		name:       "no values",
		product:    nil,
		req:        `{"name": "", "manufacturer": "", "price": 0, "amount": 0, "category": ""}`,
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Error: "validation error", Details: map[string]string{"Amount": "this field is required", "Category": "this field is required", "Manufacturer": "this field is required", "Name": "this field is required", "Price": "this field is required"}},
	},
}

func TestCreateProductHandler(t *testing.T) {
	for _, test := range testsCreate {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := MockProductRepo{
				createProduct: func(ctx context.Context, insertData map[string]interface{}) (*domain.Product, error) {
					if test.wantStatus == 409 {
						return nil, domain.ErrProductExist
					}
					return test.product, nil
				},
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/product", strings.NewReader(test.req))
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

				if test.wantStatus != w.Code {
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

				if diff := cmp.Diff(test.product, &product, opts); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}

				if test.wantStatus != w.Code {
					t.Fatalf("expected %d, got: %d", test.wantStatus, w.Code)
				}
			}
		})
	}
}
