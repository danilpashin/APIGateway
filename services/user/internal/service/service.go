package service

import "apigateway/services/user/internal/repository/postgres"

type UserService struct {
	repo postgres.UserRepository
}

func NewUserService(repo postgres.UserRepository) *UserService {
	return &UserService{repo: repo}
}
