package service

import (
	"apigateway/services/user/internal/domain"
	"context"
)

type UserServiceInterface interface {
	CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error)
	UpdateUser(ctx context.Context, id int, req *domain.UpdateUserRequest) (*domain.User, error)
	GetUser(ctx context.Context, id int) (*domain.User, error)
	ListUsers(ctx context.Context, cursor int, limit uint64) ([]*domain.User, *domain.Pagination, error)
	DeleteUser(ctx context.Context, id int) error
}
