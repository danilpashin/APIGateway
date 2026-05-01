//go:integration
package postgres

import (
	"apigateway/services/product/internal/domain"
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/go-cmp/cmp"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupTestDB(t *testing.T) *ProductRepository {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5433"
	}

	connStr := fmt.Sprintf("postgresql://postgres:test@%s:%s/products?sslmode=disable", dbHost, dbPort)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatal("failed to open database: ", err)
	}

	m, err := migrate.New("file://../../../migrations", connStr)
	if err != nil {
		t.Fatal("failed to init migrations: ", err)
	}

	m.Up()

	t.Cleanup(func() {
		db.Exec("TRUNCATE products CASCADE")
		db.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
		db.Close()
	})

	return NewPostgresProductRepository(db)
}

var insertData = map[string]any{
	"name":         "test-product",
	"manufacturer": "test-manufacturer",
	"price":        10000,
	"amount":       10,
	"status":       true,
	"category":     "Household appliances",
}

func CreateTestProduct(repo *ProductRepository, t *testing.T) {
	_, err := repo.CreateProduct(context.Background(), insertData)
	if err != nil {
		t.Fatal("failed to create test product: ", err)
	}
}

type TestCreate struct {
	name    string
	product *domain.Product
}

var testCreate = TestCreate{
	name: "general",
	product: &domain.Product{
		ID:           1,
		Name:         "test-product",
		Manufacturer: "test-manufacturer",
		Price:        10000,
		Amount:       10,
		Status:       true,
		Category:     "Household appliances",
	},
}

func TestCreateProduct(t *testing.T) {
	test := testCreate
	repo := setupTestDB(t)
	t.Run(test.name, func(t *testing.T) {
		product, err := repo.CreateProduct(context.Background(), insertData)
		if err != nil {
			t.Fatal("failed to create test product: ", err)
		}

		opts := cmp.FilterPath(func(p cmp.Path) bool {
			return p.String() == "CreatedAt" || p.String() == "UpdatedAt"
		}, cmp.Ignore())

		if diff := cmp.Diff(test.product, product, opts); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

type TestUpdate struct {
	name       string
	updateData map[string]any
	newProduct *domain.Product
}

var testUpdate = TestUpdate{
	name: "general",
	updateData: map[string]any{
		"name":         "UPD-test-product",
		"manufacturer": "UPD-test-manufacturer",
		"price":        15000,
		"amount":       12,
		"status":       false,
		"category":     "PC accessories",
	},
	newProduct: &domain.Product{
		ID:           1,
		Name:         "UPD-test-product",
		Manufacturer: "UPD-test-manufacturer",
		Price:        15000,
		Amount:       12,
		Status:       false,
		Category:     "PC accessories",
	},
}

func TestUpdateProduct(t *testing.T) {
	repo := setupTestDB(t)
	CreateTestProduct(repo, t)
	t.Run(testUpdate.name, func(t *testing.T) {
		newProduct, err := repo.UpdateProduct(context.Background(), 1, testUpdate.updateData)
		if err != nil {
			t.Fatal("failed to update test product: ", err)
		}

		opts := cmp.FilterPath(func(p cmp.Path) bool {
			return p.String() == "CreatedAt" || p.String() == "UpdatedAt"
		}, cmp.Ignore())

		if diff := cmp.Diff(testUpdate.newProduct, newProduct, opts); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	})
}

type TestGet struct {
	name      string
	productID int
	product   *domain.Product
}

var testGet = TestGet{
	name:      "general",
	productID: 1,
	product: &domain.Product{
		ID:           1,
		Name:         "test-product",
		Manufacturer: "test-manufacturer",
		Price:        10000,
		Amount:       10,
		Status:       true,
		Category:     "Household appliances",
	},
}

func TestGetProduct(t *testing.T) {
	repo := setupTestDB(t)
	CreateTestProduct(repo, t)
	t.Run(testGet.name, func(t *testing.T) {
		product, err := repo.GetProduct(context.Background(), testGet.productID)
		if err != nil {
			t.Fatal("failed to get test product: ", err)
		}

		opts := cmp.FilterPath(func(p cmp.Path) bool {
			return p.String() == "CreatedAt" || p.String() == "UpdatedAt"
		}, cmp.Ignore())

		if diff := cmp.Diff(testGet.product, product, opts); diff != "" {
			t.Fatalf("mismatсh (-want +got):\n%s", diff)
		}
	})
}

type TestList struct {
	name         string
	cursor       int
	limit        uint64
	listProducts []*domain.Product
	newCursor    int
	hasMore      bool
}

var testList = TestList{
	name:   "general",
	cursor: 1,
	limit:  2,
	listProducts: []*domain.Product{
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
	newCursor: 3,
	hasMore:   true,
}

func TestListProducts(t *testing.T) {
	repo := setupTestDB(t)
	CreateTestProduct(repo, t)
	CreateTestProduct(repo, t)
	CreateTestProduct(repo, t)
	t.Run(testList.name, func(t *testing.T) {
		products, newCursor, hasMore, err := repo.ListProducts(context.Background(), testList.cursor, testList.limit)
		if err != nil {
			t.Fatal("failed to get list of test products: ", err)
		}

		opts := cmp.FilterPath(func(p cmp.Path) bool {
			return p.String() == "CreatedAt" || p.String() == "UpdatedAt"
		}, cmp.Ignore())

		if diff := cmp.Diff(testList.listProducts[:testList.limit], products, opts); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}

		if newCursor != testList.newCursor {
			t.Errorf("expected new cursor %d, got: %d", testList.newCursor, newCursor)
		}

		if hasMore != testList.hasMore {
			t.Fatalf("expected has more %v, got: %v", testList.hasMore, hasMore)
		}
	})
}

type TestDelete struct {
	name      string
	productID int
	result    error
}

var testDelete = TestDelete{
	name:      "general",
	productID: 1,
	result:    nil,
}

func TestDeleteProduct(t *testing.T) {
	repo := setupTestDB(t)
	CreateTestProduct(repo, t)
	t.Run(testDelete.name, func(t *testing.T) {
		err := repo.DeleteProduct(context.Background(), testDelete.productID)
		if err != nil {
			t.Fatal("failed to delete product: ", err)
		}
	})
}
