package service

import (
	"apigateway/services/user/internal/domain"
	"apigateway/services/user/internal/repository/postgres"
	"apigateway/services/user/internal/validator"
	"context"
	"errors"
	"pkg/env"
	"pkg/response"

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

func (s *UserService) CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	if err := validator.New(req); err != nil {
		errorData := response.FormatValidationError(err)
		return nil, domain.ErrorResponse{Message: "validation error", Details: errorData}
	}

	insertData := make(map[string]any)

	if req.Username != "" {
		insertData["username"] = req.Username
	}
	if req.Email != "" {
		insertData["email"] = req.Email
	}
	if req.Password != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptCost)
		if err != nil {
			return nil, errors.New("failed to generate password hash")
		}
		insertData["password_hash"] = string(passwordHash)
	}

	return s.repo.CreateUser(ctx, insertData)
}

func (s *UserService) UpdateUser(ctx context.Context, id int, req *domain.UpdateUserRequest) (*domain.User, error) {
	if err := validator.New(req); err != nil {
		errorData := response.FormatValidationError(err)
		return nil, domain.ErrorResponse{Message: "validation error", Details: errorData}
	}

	currentUser, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	updateData := make(map[string]any)

	if req.Username != nil {
		if *req.Username != "" {
			updateData["username"] = *req.Username
		}
	}

	if req.Password != nil {
		if err = bcrypt.CompareHashAndPassword([]byte(currentUser.PasswordHash), []byte(*req.Password)); err != nil {
			return nil, domain.ErrorResponse{Message: "incorrect password"}
		} else {
			if req.Email != nil {
				if *req.Email != "" {
					updateData["email"] = *req.Email
				}
			}
			if req.NewPassword != nil {
				if *req.NewPassword != "" {
					passwordHash, err := bcrypt.GenerateFromPassword([]byte(*req.NewPassword), bcryptCost)
					if err != nil {
						return nil, errors.New("failed to generate password hash")
					}
					updateData["password_hash"] = passwordHash
				}
			}
		}
	}

	if len(updateData) == 0 {
		return nil, domain.ErrNoUpdateData
	}

	return s.repo.UpdateUser(ctx, id, updateData)
}

func (s *UserService) GetUser(ctx context.Context, id int) (*domain.User, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidID
	}

	return s.repo.GetUser(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context, cursor int, limit uint64) ([]*domain.User, *domain.Pagination, error) {
	if cursor < 0 {
		cursor = 0
	}
	if limit <= 0 {
		limit = 10
	} else if limit > 50 {
		limit = 50
	}

	listUsers, nextCursor, hasMore, err := s.repo.ListUsers(ctx, cursor, limit)
	if err != nil {
		return nil, nil, err
	}

	return listUsers, &domain.Pagination{NextCursor: nextCursor, HasMore: hasMore, Limit: limit}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	if id <= 0 {
		return domain.ErrInvalidID
	}

	return s.repo.DeleteUser(ctx, id)
}
