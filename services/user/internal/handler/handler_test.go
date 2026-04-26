package handler

import (
	"apigateway/services/user/internal/domain"
	"context"
)

type MockUserRepo struct {
	createUser func(ctx context.Context, insertData map[string]interface{}) (*domain.User, error)
	updateUser func(ctx context.Context, id int, updateData map[string]interface{}) (*domain.User, error)
	getUser    func(ctx context.Context, id int) (*domain.User, error)
	listUsers  func(ctx context.Context, cursor int, limit uint64) ([]*domain.User, int, bool, error)
	deleteUser func(ctx context.Context, id int) error
}

func (m *MockUserRepo) CreateUser(ctx context.Context, insertData map[string]interface{}) (*domain.User, error) {
	if m.createUser != nil {
		return m.createUser(ctx, insertData)
	}
	return nil, nil
}

func (m *MockUserRepo) UpdateUser(ctx context.Context, id int, updateData map[string]interface{}) (*domain.User, error) {
	if m.updateUser != nil {
		return m.updateUser(ctx, id, updateData)
	}
	return nil, nil
}

func (m *MockUserRepo) GetUser(ctx context.Context, id int) (*domain.User, error) {
	if m.getUser != nil {
		return m.getUser(ctx, id)
	}
	return nil, nil
}

func (m *MockUserRepo) ListUsers(ctx context.Context, cursor int, limit uint64) ([]*domain.User, int, bool, error) {
	if m.listUsers != nil {
		return m.listUsers(ctx, cursor, limit)
	}
	return nil, 0, false, nil
}

func (m *MockUserRepo) DeleteUser(ctx context.Context, id int) error {
	if m.deleteUser != nil {
		return m.deleteUser(ctx, id)
	}
	return nil
}
