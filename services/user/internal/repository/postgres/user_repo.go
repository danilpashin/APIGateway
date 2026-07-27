package postgres

import (
	"apigateway/services/user/internal/domain"
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
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

	err = r.pool.QueryRow(ctx, query, args...).
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

	err = r.pool.QueryRow(ctx, query, args...).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUser(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User
	query := `SELECT u.id, u.username, u.email, u.password_hash, u.created_at, u.updated_at, r.name AS role_name FROM users u LEFT JOIN roles r ON u.role_id = r.id WHERE u.id = $1`

	err := r.pool.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) ListUsers(ctx context.Context, cursor int, limit uint64) ([]*domain.User, int, bool, error) {
	extLimit := limit + 1
	query := `SELECT u.id, u.username, u.email, u.password_hash, u.created_at, u.updated_at, r.name AS role_name FROM users u LEFT JOIN roles r ON u.role_id = r.id WHERE u.id >= $1 LIMIT $2`

	rows, err := r.pool.Query(ctx, query, cursor, extLimit)
	if err != nil {
		return nil, 0, false, err
	}
	defer rows.Close()

	listUsers := make([]*domain.User, 0, int(extLimit))
	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.Role)
		if err != nil {
			return nil, 0, false, err
		}
		listUsers = append(listUsers, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, false, err
	}

	if len(listUsers) > int(limit) {
		nextCursor := listUsers[int(limit)].ID
		hasMore := true
		return listUsers[:limit], nextCursor, hasMore, nil
	}
	listUsersClip := slices.Clip(listUsers)

	return listUsersClip, 0, false, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	commandTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("pool exec error: %w", err)
	}

	rows := commandTag.RowsAffected()

	if rows == 0 {
		return domain.ErrUserNotFound
	}

	if rows != 1 {
		return fmt.Errorf("expected to affect 1 row, affected: %d", rows)
	}

	return nil
}
