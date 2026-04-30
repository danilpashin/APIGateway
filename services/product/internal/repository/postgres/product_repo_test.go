//go:integration
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
		db.Close()
	})

	return NewPostgresProductRepository(db)
}

func TestCreateProduct(t *testing.T) {
	repo := setupTestDB(t)

	insertData := make(map[string]any)
	insertData["name"] = "test-product"
	insertData["manufacturer"] = "test-manufacturer"
	insertData["price"] = 10000
	insertData["amount"] = 10
	insertData["status"] = true
	insertData["category"] = "Household appliances"

	product, err := repo.CreateProduct(context.Background(), insertData)
	if err != nil {
		t.Fatal("failed to create product: ", err)
	}
	log.Print(product)
}
