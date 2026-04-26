package service

import (
	"apigateway/services/user/internal/domain"
	"apigateway/services/user/internal/repository/postgres"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo postgres.UserRepository
}

func NewUserService(repo postgres.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error) {
	insertData := make(map[string]interface{})

	if req.Username != "" {
		insertData["username"] = req.Username
	}
	if req.Email != "" {
		insertData["email"] = req.Email
	}
	if req.Password != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
		if err != nil {
			return nil, errors.New("failed to generate password hash")
		}
		insertData["password_hash"] = passwordHash
	}

	return nil, nil
}
