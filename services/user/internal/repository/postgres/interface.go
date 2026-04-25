package postgres

import (
	"apigateway/services/user/internal/domain"
	"context"
)

type UserRepoInterface interface {
	CreateUser(ctx context.Context, insertData map[string]interface{}) (*domain.User, error)
	UpdateUser(ctx context.Context, id int, updateData map[string]interface{}) (*domain.User, error)
	GetUser(ctx context.Context, id int) (*domain.User, error)
	ListUsers(ctx context.Context, cursor int, limit uint64) ([]*domain.User, int, bool, error)
	DeleteUser(ctx context.Context, id int) error
}
