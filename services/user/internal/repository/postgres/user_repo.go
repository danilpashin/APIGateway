package postgres

import (
	"apigateway/services/user/internal/domain"
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, insertData map[string]any) (*domain.User, error) {
	var user domain.User
	builder := squirrel.Insert("users").
		SetMap(insertData).
		Suffix(`RETURNING id, username, email, password_hash, (SELECT name FROM roles WHERE id = 1) AS role_name, created_at, updated_at`).
		PlaceholderFormat(squirrel.Dollar)
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.New("failed insert query build")
	}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, errors.New("failed insert query")
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id int, updateData map[string]any) (*domain.User, error) {
	var user domain.User
	var role_id int
	builder := squirrel.Update("users").SetMap(updateData).Where(squirrel.Eq{"id": id}).Suffix(`RETURNING id, username, email, password_hash, role_id, created_at, updated_at`)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.New("failed update query build")
	}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &role_id, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	query = `SELECT name FROM roles WHERE id = $1`

	err = r.db.QueryRowContext(ctx, query, role_id).Scan(&user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUser(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User
	var role_id int
	query := `SELECT * FROM users WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &role_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	query = `SELECT name FROM roles WHERE id = $1`

	err = r.db.QueryRowContext(ctx, query, role_id).Scan(&user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) ListUsers(ctx context.Context, cursor int, limit uint64) ([]*domain.User, int, bool, error) {
	return nil, 0, false, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	return nil
}
