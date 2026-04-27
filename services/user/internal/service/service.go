package service

import (
	"apigateway/services/user/internal/domain"
	"apigateway/services/user/internal/repository/postgres"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo postgres.UserRepoInterface
}

func NewUserService(repo postgres.UserRepoInterface) *UserService {
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

	if len(insertData) < 3 {
		return nil, errors.New("insert data is empty or not enough")
	}

	return s.repo.CreateUser(ctx, insertData)
}

func (s *UserService) UpdateUser(ctx context.Context, id int, req domain.UpdateUserRequest) (*domain.User, error) {
	currentUser, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	updateData := make(map[string]interface{})

	if req.Username != "" {
		updateData["username"] = req.Username
	}
	if req.Email != "" {
		updateData["email"] = req.Email
	}
	if req.NewPassword != "" {
		if err = bcrypt.CompareHashAndPassword([]byte(currentUser.PasswordHash), []byte(req.Password)); err != nil {
			return nil, errors.New("wrong password")
		}
		if len(req.NewPassword) < 8 {
			return nil, errors.New("password must be at least 8 characters")
		}
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
		if err != nil {
			return nil, errors.New("failed to generate password hash")
		}
		updateData["password_hash"] = passwordHash
	}

	return s.repo.UpdateUser(ctx, id, updateData)
}
