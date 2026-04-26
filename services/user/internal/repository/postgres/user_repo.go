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

func (r *UserRepository) CreateUser(ctx context.Context, insertData map[string]interface{}) (*domain.User, error) {
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

func (r *UserRepository) UpdateUser(ctx context.Context, id int, updateData map[string]interface{}) (*domain.User, error) {
	return nil, nil
}

func (r *UserRepository) GetUser(ctx context.Context, id int) (*domain.User, error) {
	return nil, nil
}

func (r *UserRepository) ListUsers(ctx context.Context, cursor int, limit uint64) ([]*domain.User, error) {
	return nil, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	return nil
}
