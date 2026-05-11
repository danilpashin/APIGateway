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
	"github.com/google/go-cmp/cmp/cmpopts"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)

	m, err := migrate.New("file://../../../migrations", connStr)
	require.NoError(t, err)

	m.Up()

	t.Cleanup(func() {
		db.Exec("TRUNCATE products CASCADE")
		db.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
		db.Close()
	})

	return NewProductRepository(db)
}

var insertData = map[string]any{
	"name":         "test-product",
	"manufacturer": "test-manufacturer",
	"price":        10000,
	"amount":       10,
	"status":       true,
	"category":     "Household appliances",
}

func CreateTestProduct(repo *ProductRepository, t *testing.T) *domain.Product {
	product, err := repo.CreateProduct(context.Background(), insertData)
	require.NoError(t, err)
	return product
}

type TestCreate struct {
	name    string
	product *domain.Product
}

var testCreate = TestCreate{
	name: "success",
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

func TestProductRepository_Create(t *testing.T) {
	repo := setupTestDB(t)
	t.Run(testCreate.name, func(t *testing.T) {
		product, err := repo.CreateProduct(context.Background(), insertData)
		require.NoError(t, err)

		opts := cmpopts.IgnoreFields(domain.Product{}, "CreatedAt", "UpdatedAt")
		diff := cmp.Diff(testCreate.product, product, opts)
		assert.Empty(t, diff, "mismatch (-want +got):\n%s", diff)
	})
}

type TestUpdate struct {
	name       string
	updateData map[string]any
	newProduct *domain.Product
}

var testUpdate = TestUpdate{
	name: "success",
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

func TestProductRepository_Update(t *testing.T) {
	repo := setupTestDB(t)
	CreateTestProduct(repo, t)
	t.Run(testUpdate.name, func(t *testing.T) {
		newProduct, err := repo.UpdateProduct(context.Background(), 1, testUpdate.updateData)
		require.NoError(t, err)

		opts := cmpopts.IgnoreFields(domain.Product{}, "CreatedAt", "UpdatedAt")
		diff := cmp.Diff(testUpdate.newProduct, newProduct, opts)
		assert.Empty(t, diff, "mismatch (-want +got):\n%s", diff)
	})
}

type TestGet struct {
	name      string
	productID int
}

var testGet = TestGet{
	name:      "success",
	productID: 1,
}

func TestProductRepository_Get(t *testing.T) {
	repo := setupTestDB(t)
	t.Run(testGet.name, func(t *testing.T) {
		expected := CreateTestProduct(repo, t)
		product, err := repo.GetProduct(context.Background(), testGet.productID)
		require.NoError(t, err)

		opts := cmpopts.IgnoreFields(domain.Product{}, "CreatedAt", "UpdatedAt")
		diff := cmp.Diff(expected, product, opts)
		assert.Empty(t, diff, "mismatch (-want +got):\n%s", diff)
	})
}

type TestList struct {
	name    string
	cursor  int
	limit   uint64
	len     int
	hasMore bool
}

var testList = TestList{
	name:    "success",
	cursor:  1,
	limit:   2,
	len:     3,
	hasMore: true,
}

func TestProductRepository_List(t *testing.T) {
	repo := setupTestDB(t)

	var createdProducts []*domain.Product
	for range testList.len {
		p := CreateTestProduct(repo, t)
		createdProducts = append(createdProducts, p)
	}

	t.Run(testList.name, func(t *testing.T) {
		expected := createdProducts[:testList.limit]

		products, newCursor, hasMore, err := repo.ListProducts(context.Background(), testList.cursor, testList.limit)
		require.NoError(t, err)

		assert.Len(t, products, int(testList.limit))
		opts := cmpopts.IgnoreFields(domain.Product{}, "CreatedAt", "UpdatedAt")
		diff := cmp.Diff(expected, products, opts)
		assert.Empty(t, diff, "mismatch (-want +got):\n%s", diff)

		assert.Equal(t, createdProducts[testList.len-1].ID, newCursor)
		assert.True(t, hasMore)
	})
}

type TestDelete struct {
	name      string
	productID int
	result    error
}

var testDelete = TestDelete{
	name:      "success",
	productID: 1,
	result:    nil,
}

func TestProductRepository_Delete(t *testing.T) {
	repo := setupTestDB(t)
	CreateTestProduct(repo, t)
	t.Run(testDelete.name, func(t *testing.T) {
		err := repo.DeleteProduct(context.Background(), testDelete.productID)
		require.NoError(t, err)
	})
}
