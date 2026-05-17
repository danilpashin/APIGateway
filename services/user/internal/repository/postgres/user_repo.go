package postgres

import (
	"apigateway/services/user/internal/domain"
	"context"
	"database/sql"

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
		Suffix(`RETURNING id, username, email, password_hash, (SELECT name FROM roles WHERE id = role_id) AS role_name, created_at, updated_at`).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id int, updateData map[string]any) (*domain.User, error) {
	var user domain.User
	builder := squirrel.Update("users").
		SetMap(updateData).Where(squirrel.Eq{"id": id}).
		Suffix(`RETURNING id, username, email, password_hash, (SELECT name FROM roles WHERE id = role_id) AS role_name, created_at, updated_at`).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUser(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User
	query := `SELECT u.id, u.username, u.email, u.password_hash, u.created_at, u.updated_at, r.name AS role_name FROM users u LEFT JOIN roles r ON u.role_id = r.id WHERE u.id = $1`

	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
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
