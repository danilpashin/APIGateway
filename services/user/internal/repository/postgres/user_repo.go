package postgres

import (
	"apigateway/services/user/internal/domain"
	"context"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, insertData map[string]interface{}) (*domain.User, error) {
	return nil, nil
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
