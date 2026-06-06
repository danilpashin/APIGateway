package service

import (
	"apigateway/services/user/internal/domain"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) CreateUser(ctx context.Context, insertData map[string]any) (*domain.User, error) {
	args := m.Called(ctx, insertData)
	p, _ := args[0].(*domain.User)
	return p, args.Error(1)
}

func (m *MockRepo) UpdateUser(ctx context.Context, id int, updateData map[string]any) (*domain.User, error) {
	args := m.Called(ctx, id, updateData)
	p, _ := args[0].(*domain.User)
	return p, args.Error(1)
}

func (m *MockRepo) GetUser(ctx context.Context, id int) (*domain.User, error) {
	args := m.Called(ctx, id)
	p, _ := args[0].(*domain.User)
	return p, args.Error(1)
}

func (m *MockRepo) ListUsers(ctx context.Context, cursor int, limit uint64) ([]*domain.User, int, bool, error) {
	args := m.Called(ctx, cursor, limit)
	p, _ := args[0].([]*domain.User)
	return p, args.Int(1), args.Bool(2), args.Error(3)
}

func (m *MockRepo) DeleteUser(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type TestCreate struct {
	name      string
	input     domain.CreateUserRequest
	mockInput map[string]any
	mockResp  *domain.User
	mockError error
	want      *domain.User
	wantErr   bool
	wantResp  error
}

var testsCreate = []TestCreate{
	{
		name: "success: all values",
		input: domain.CreateUserRequest{
			Username: "test-user",
			Email:    "test@gmail.com",
			Password: "test-password-123",
		},
		mockInput: map[string]any{
			"username":      "test-user",
			"email":         "test@gmail.com",
			"password_hash": "test-password-hash-123",
		},
		mockResp: &domain.User{
			ID:           1,
			Username:     "test-user",
			Email:        "test@gmail.com",
			PasswordHash: "test-password-hash-123",
			Role:         "user",
		},
		want: &domain.User{
			ID:           1,
			Username:     "test-user",
			Email:        "test@gmail.com",
			PasswordHash: "test-password-hash-123",
			Role:         "user",
		},
		wantErr: false,
	},
	{
		name: "error: missing required field: username",
		input: domain.CreateUserRequest{
			Email:    "test@gmail.com",
			Password: "test-password-123",
		},
		wantErr:  true,
		wantResp: domain.ErrorResponse{Message: "validation error", Details: map[string]string{"Username": "this field is required"}},
	},
	{
		name: "error: invalid name format",
		input: domain.CreateUserRequest{
			Username: "t",
			Email:    "test@gmail.com",
			Password: "test-password-123",
		},
		wantErr:  true,
		wantResp: domain.ErrorResponse{Message: "validation error", Details: map[string]string{"Username": "must be at least 3 characters"}},
	},
	{
		name: "error: missing required field: email",
		input: domain.CreateUserRequest{
			Username: "test-user",
			Password: "test-password-123",
		},
		wantErr:  true,
		wantResp: domain.ErrorResponse{Message: "validation error", Details: map[string]string{"Email": "this field is required"}},
	},
	{
		name: "error: invalid email format",
		input: domain.CreateUserRequest{
			Username: "test-user",
			Email:    "1",
			Password: "test-password-123",
		},
		wantErr:  true,
		wantResp: domain.ErrorResponse{Message: "validation error", Details: map[string]string{"Email": "invalid email format"}},
	},
	{
		name: "error: missing required field: password",
		input: domain.CreateUserRequest{
			Username: "test-user",
			Email:    "test@gmail.com",
		},
		wantErr:  true,
		wantResp: domain.ErrorResponse{Message: "validation error", Details: map[string]string{"Password": "this field is required"}},
	},
	{
		name: "error: invalid password format",
		input: domain.CreateUserRequest{
			Username: "test-user",
			Email:    "test@gmail.com",
			Password: "123",
		},
		wantErr:  true,
		wantResp: domain.ErrorResponse{Message: "validation error", Details: map[string]string{"Password": "must be at least 8 characters"}},
	},
}

func TestUserService_Create(t *testing.T) {
	for _, tt := range testsCreate {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			if !tt.wantErr {
				mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(tt.mockResp, tt.wantResp)
			}

			userService := NewUserService(mockRepo)

			user, err := userService.CreateUser(context.Background(), &tt.input)

			if tt.wantErr {
				assert.Equal(t, tt.want, user)
				assert.Equal(t, tt.wantResp, err)
			} else {
				assert.Equal(t, tt.want, user)
				assert.Equal(t, nil, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

type TestUpdate struct {
	name           string
	input          domain.UpdateUserRequest
	mockID         int
	mockGetResp    *domain.User
	mockGetErr     error
	mockUpdateResp *domain.User
	mockUpdateErr  error
	want           *domain.User
	wantErr        bool
	wantResp       error
}

var testsUpdate = []TestUpdate{
	{
		name: "success",
		input: domain.UpdateUserRequest{
			Username:    stringPtr("UPD-test-user"),
			Email:       stringPtr("upd-test@gmail.com"),
			Password:    stringPtr("test-password-123"),
			NewPassword: stringPtr("UPD-test-password-123"),
		},
		mockID: 1,
		mockGetResp: &domain.User{
			ID:           1,
			Username:     "UPD-test-user",
			Email:        "upd-test@gmail.com",
			PasswordHash: "$2a$13$io.K4Ps.bMCNo/D2SfXlvejKnMjBmQkhQzJspzS0BeMNyrwkTfN0q",
		},
		mockUpdateResp: &domain.User{
			ID:           1,
			Username:     "UPD-test-user",
			Email:        "upd-test@gmail.com",
			PasswordHash: "$2a$13$io.K4Ps.bMCNo/D2SfXlvejKnMjBmQkhQzJspzS0BeMNyrwkTfN0q",
		},
		want: &domain.User{
			ID:           1,
			Username:     "UPD-test-user",
			Email:        "upd-test@gmail.com",
			PasswordHash: "$2a$13$io.K4Ps.bMCNo/D2SfXlvejKnMjBmQkhQzJspzS0BeMNyrwkTfN0q",
		},
		wantErr: false,
	},
	{
		name: "error: wrong password",
		input: domain.UpdateUserRequest{
			Username: stringPtr("UPD-test-user"),
			Email:    stringPtr("upd-test@gmail.com"),
			Password: stringPtr("test-password-123"),
		},
		mockGetResp: &domain.User{
			ID:           1,
			Username:     "UPD-test-user",
			Email:        "upd-test@gmail.com",
			PasswordHash: "wrong-hash",
		},
		wantErr:  true,
		wantResp: domain.ErrorResponse{Message: "incorrect password"},
	},
	{
		name: "error: missing all update values",
		input: domain.UpdateUserRequest{
			Username:    nil,
			Email:       nil,
			Password:    nil,
			NewPassword: nil,
		},
		mockID:   1,
		wantErr:  true,
		wantResp: domain.ErrNoUpdateData,
	},
	{
		name: "error: invalid username format",
		input: domain.UpdateUserRequest{
			Username: stringPtr("t"),
		},
		mockID:   1,
		wantErr:  true,
		wantResp: domain.ErrorResponse{Message: "validation error", Details: map[string]string{"Username": "must be at least 3 characters"}},
	},
	{
		name: "error: invalid email format",
		input: domain.UpdateUserRequest{
			Username: stringPtr("UPD-test-user"),
			Email:    stringPtr("1"),
		},
		mockID:   1,
		wantErr:  true,
		wantResp: domain.ErrorResponse{Message: "validation error", Details: map[string]string{"Email": "invalid email format"}},
	},
	{
		name: "error: invalid password format",
		input: domain.UpdateUserRequest{
			Username:    stringPtr("UPD-Test-user"),
			Email:       stringPtr("upd-test@gmail.com"),
			Password:    stringPtr("test-password-123"),
			NewPassword: stringPtr("123"),
		},
		mockID: 1,
		mockGetResp: &domain.User{
			ID:           1,
			Username:     "UPD-test-user",
			Email:        "upd-test@gmail.com",
			PasswordHash: "$2a$13$io.K4Ps.bMCNo/D2SfXlvejKnMjBmQkhQzJspzS0BeMNyrwkTfN0q",
		},
		wantErr:  true,
		wantResp: domain.ErrorResponse{Message: "validation error", Details: map[string]string{"NewPassword": "must be at least 8 characters"}},
	},
	{
		name: "error: user not found",
		input: domain.UpdateUserRequest{
			Username:    stringPtr("UPD-Test-user"),
			Email:       stringPtr("upd-test@gmail.com"),
			Password:    stringPtr("test-password-123"),
			NewPassword: stringPtr("UPD-test-password-123"),
		},
		mockID:     1,
		mockGetErr: domain.ErrUserNotFound,
		wantErr:    true,
		wantResp:   domain.ErrUserNotFound,
	},
}

func TestUserService_Update(t *testing.T) {
	for _, tt := range testsUpdate {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			mockRepo.On("GetUser", mock.Anything, mock.Anything).Return(tt.mockGetResp, tt.mockGetErr)
			mockRepo.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(tt.mockUpdateResp, tt.mockUpdateErr)

			userService := NewUserService(mockRepo)

			user, err := userService.UpdateUser(context.Background(), tt.mockID, &tt.input)

			if tt.wantErr {
				assert.Equal(t, tt.want, user)
				assert.Equal(t, tt.wantResp, err)
			} else {
				assert.Equal(t, tt.want, user)
				assert.Equal(t, nil, err)
			}
		})
	}
}

func stringPtr(s string) *string { return &s }

type TestGet struct {
	name     string
	mockID   int
	mockResp *domain.User
	want     *domain.User
	wantErr  bool
	wantResp error
}

var testsGet = []TestGet{
	{
		name:   "success",
		mockID: 1,
		mockResp: &domain.User{
			ID:       1,
			Username: "Test-user",
			Email:    "test@gmail.com",
		},
		want: &domain.User{
			ID:       1,
			Username: "Test-user",
			Email:    "test@gmail.com",
		},
		wantErr: false,
	},
	{
		name:     "error: invalid ID format",
		mockID:   -1,
		wantErr:  true,
		wantResp: domain.ErrInvalidID,
	},
}

func TestUserService_Get(t *testing.T) {
	for _, tt := range testsGet {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			if !tt.wantErr {
				mockRepo.On("GetUser", mock.Anything, tt.mockID).Return(tt.mockResp, tt.wantResp)
			}

			userService := NewUserService(mockRepo)

			user, err := userService.GetUser(context.Background(), tt.mockID)

			if tt.wantErr {
				assert.Equal(t, tt.want, user)
				assert.Equal(t, tt.wantResp, err)
			} else {
				assert.Equal(t, tt.want, user)
				assert.Equal(t, nil, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

type TestList struct {
	name           string
	cursor         int
	limit          uint64
	mockNextCursor int
	mockHasMore    bool
	mockCursor     int
	mockLimit      uint64
	want           []*domain.User
	wantPagination *domain.Pagination
	wantErr        bool
	wantResp       error
}

var testsList = []TestList{
	{
		name:           "success: first two users (cursor=1, limit=2)",
		cursor:         1,
		limit:          2,
		mockNextCursor: 3,
		mockHasMore:    true,
		mockCursor:     1,
		mockLimit:      2,
		want: []*domain.User{
			{
				ID:       1,
				Username: "test-user",
				Email:    "test@gmail.com",
			},
			{
				ID:       2,
				Username: "test-user",
				Email:    "test@gmail.com",
			},
			{
				ID:       3,
				Username: "test-user",
				Email:    "test@gmail.com",
			},
		},
		wantPagination: &domain.Pagination{
			NextCursor: 3,
			HasMore:    true,
			Limit:      2,
		},
		wantErr: false,
	},
	{
		name:       "success: end of list (no more)",
		cursor:     3,
		limit:      3,
		mockCursor: 3,
		mockLimit:  3,
		want: []*domain.User{
			{
				ID:       1,
				Username: "test-user",
				Email:    "test@gmail.com",
			},
			{
				ID:       2,
				Username: "test-user",
				Email:    "test@gmail.com",
			},
			{
				ID:       3,
				Username: "test-user",
				Email:    "test@gmail.com",
			},
		},
		wantPagination: &domain.Pagination{
			NextCursor: 0,
			HasMore:    false,
			Limit:      3,
		},
		wantErr: false,
	},
	{
		name:       "success: negative cursor and null limit",
		cursor:     -2,
		limit:      0,
		mockCursor: 0,
		mockLimit:  10,
		want: []*domain.User{
			{
				ID:       1,
				Username: "test-user",
				Email:    "test@gmail.com",
			},
			{
				ID:       2,
				Username: "test-user",
				Email:    "test@gmail.com",
			},
			{
				ID:       3,
				Username: "test-user",
				Email:    "test@gmail.com",
			},
		},
		wantPagination: &domain.Pagination{
			NextCursor: 0,
			HasMore:    false,
			Limit:      10,
		},
		wantErr: false,
	},
	{
		name:       "success: limit > 50 clamped to 50",
		cursor:     1,
		limit:      100,
		mockCursor: 1,
		mockLimit:  50,
		want: []*domain.User{
			{
				ID:       1,
				Username: "test-user",
				Email:    "test@gmail.com",
			},
			{
				ID:       2,
				Username: "test-user",
				Email:    "test@gmail.com",
			},
			{
				ID:       3,
				Username: "test-user",
				Email:    "test@gmail.com",
			},
		},
		wantPagination: &domain.Pagination{
			NextCursor: 0,
			HasMore:    false,
			Limit:      50,
		},
		wantErr: false,
	},
	{
		name:       "error: repository error",
		cursor:     1,
		limit:      2,
		mockCursor: 1,
		mockLimit:  2,
		wantErr:    true,
		wantResp:   errors.New("repo error"),
	},
}

func TestUserService_List(t *testing.T) {
	for _, tt := range testsList {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			mockRepo.On("ListUsers", mock.Anything, tt.mockCursor, tt.mockLimit).Return(tt.want, tt.mockNextCursor, tt.mockHasMore, tt.wantResp)

			userService := NewUserService(mockRepo)

			users, pagination, err := userService.ListUsers(context.Background(), tt.cursor, tt.limit)

			if tt.wantErr {
				assert.Equal(t, tt.want, users)
				assert.Equal(t, tt.wantResp, err)
			} else {
				assert.Equal(t, tt.want, users)
				assert.Equal(t, tt.wantPagination, pagination)
				assert.Equal(t, nil, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

type TestDelete struct {
	name     string
	mockID   int
	wantErr  bool
	wantResp error
}

var testsDelete = []TestDelete{
	{
		name:    "success",
		mockID:  1,
		wantErr: false,
	},
	{
		name:     "error: invalid ID",
		mockID:   -1,
		wantErr:  true,
		wantResp: domain.ErrInvalidID,
	},
}

func TestUserService_Delete(t *testing.T) {
	for _, tt := range testsDelete {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			if !tt.wantErr {
				mockRepo.On("DeleteUser", mock.Anything, tt.mockID).Return(tt.wantResp)
			}

			userService := NewUserService(mockRepo)

			err := userService.DeleteUser(context.Background(), tt.mockID)

			if tt.wantErr {
				assert.Equal(t, tt.wantResp, err)
			} else {
				assert.Equal(t, nil, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
