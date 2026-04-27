package service

import (
	"apigateway/services/user/internal/domain"
	"apigateway/services/user/internal/repository/postgres"
	"context"
	"errors"
	"pkg/env"

	"golang.org/x/crypto/bcrypt"
)

// Using bcrypt with cost 13 because it provides balance between security and performance.
// Lower costs (10-12) are too weak against modern GPU attacks.
var bcryptCost = env.GetEnvAsInt("BCRYPT_COST", 13)

type UserService struct {
	repo postgres.UserRepoInterface
}

func NewUserService(repo postgres.UserRepoInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error) {
	insertData := make(map[string]any)

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
		return nil, domain.ErrNoInsertData
	}

	return s.repo.CreateUser(ctx, insertData)
}

func (s *UserService) UpdateUser(ctx context.Context, id int, req domain.UpdateUserRequest) (*domain.User, error) {
	currentUser, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	updateData := make(map[string]any)

	if req.Username != "" {
		updateData["username"] = req.Username
	}

	if err = bcrypt.CompareHashAndPassword([]byte(currentUser.PasswordHash), []byte(req.Password)); err != nil {
		return nil, domain.ErrWrongPassword
	} else {
		if req.Email != "" {
			updateData["email"] = req.Email
		}
		if req.NewPassword != "" {
			if len(req.NewPassword) < 8 {
				return nil, domain.ErrInvalidPassword
			}
			passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
			if err != nil {
				return nil, errors.New("failed to generate password hash")
			}
			updateData["password_hash"] = passwordHash
		}
	}

	return s.repo.UpdateUser(ctx, id, updateData)
}

func (s *UserService) GetUser(ctx context.Context, id int) (*domain.User, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidID
	}

	return s.repo.GetUser(ctx, id)
}
