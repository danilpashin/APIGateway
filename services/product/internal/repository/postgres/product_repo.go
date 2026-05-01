package postgres

import (
	"apigateway/services/product/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
)

type ProductRepository struct {
	db *sql.DB
}

func NewPostgresProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, insertData map[string]any) (*domain.Product, error) {
	var product domain.Product
	builder := squirrel.Insert("products").
		SetMap(insertData).
		Suffix("RETURNING id, name, manufacturer, price, amount, status, category, created_at, updated_at").
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.New("invalid insert query")
	}

	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&product.ID, &product.Name, &product.Manufacturer, &product.Price, &product.Amount, &product.Status, &product.Category, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			return nil, domain.ErrProductExist
		}
		return nil, err
	}
	log.Printf("Product name is %s, created on %s\n", product.Name, product.CreatedAt)

	return &product, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, id int, updateData map[string]any) (*domain.Product, error) {
	var product domain.Product
	builder := squirrel.Update("products").
		SetMap(updateData).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id, name, manufacturer, price, amount, status, category, created_at, updated_at").
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, domain.ErrUpdateQuery
	}

	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&product.ID, &product.Name, &product.Manufacturer, &product.Price, &product.Amount, &product.Status, &product.Category, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrProductsNotFound
		}
		return nil, err
	}
	log.Printf("Product name is %s, updated on %s\n", product.Name, product.UpdatedAt)

	return &product, nil
}

func (r *ProductRepository) GetProduct(ctx context.Context, id int) (*domain.Product, error) {
	var product domain.Product
	builder := squirrel.Select("*").From("products").Where(squirrel.Eq{"id": id}).PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, domain.ErrGetQuery
	}

	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&product.ID, &product.Name, &product.Manufacturer, &product.Price, &product.Amount, &product.Status, &product.Category, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) ListProducts(ctx context.Context, cursor int, limit uint64) ([]*domain.Product, int, bool, error) {
	listProducts := make([]*domain.Product, 0)

	builder := squirrel.Select("*").From("products").
		Where(squirrel.GtOrEq{"id": cursor}).Limit(limit + 1).OrderBy("id ASC").
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, 0, false, domain.ErrListQuery
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, false, fmt.Errorf("query failed: %v", err)
	}

	for rows.Next() {
		product := new(domain.Product)
		err = rows.Scan(&product.ID, &product.Name, &product.Manufacturer, &product.Price, &product.Amount, &product.Status, &product.Category, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, 0, false, fmt.Errorf("scan failed: %w", err)
		}
		listProducts = append(listProducts, product)
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, 0, false, err
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
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id int) error {
	builder := squirrel.Delete("products").Where(squirrel.Eq{"id": id}).PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return domain.ErrDeleteQuery
	}

	err = r.db.QueryRowContext(ctx, query, args...).Scan()
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	return nil
}
