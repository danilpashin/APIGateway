package handler

import (
	"apigateway/services/user/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockService) UpdateUser(ctx context.Context, id int, req *domain.UpdateUserRequest) (*domain.User, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockService) GetUser(ctx context.Context, id int) (*domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockService) ListUsers(ctx context.Context, cursor int, limit uint64) ([]*domain.User, *domain.Pagination, error) {
	args := m.Called(ctx, cursor, limit)
	return args.Get(0).([]*domain.User), args.Get(1).(*domain.Pagination), args.Error(2)
}

func (m *MockService) DeleteUser(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type TestCreate struct {
	name       string
	req        string
	mockInput  *domain.CreateUserRequest
	mockResp   *domain.User
	mockErr    error
	resp       *domain.CreateUserResponse
	wantErr    bool
	wantStatus int
	wantResp   error
}

var testsCreate = []TestCreate{
	{
		name:       "success",
		req:        `{"username":"test-user", "email":"test@gmail.com", "password":"test"}`,
		mockInput:  &domain.CreateUserRequest{Username: "test-user", Email: "test@gmail.com", Password: "test"},
		mockResp:   &domain.User{Username: "test-user", Email: "test@gmail.com"},
		resp:       &domain.CreateUserResponse{Username: "test-user", Email: "test@gmail.com"},
		wantErr:    false,
		wantStatus: 201,
	},
	{
		name:       "error: empty insert data",
		req:        `{"username":"", "email":"", "password":""}`,
		mockInput:  &domain.CreateUserRequest{Username: "", Email: "", Password: ""},
		mockResp:   nil,
		mockErr:    domain.ErrNoInsertData,
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrNoInsertData.Error()},
	},
	{
		name:       "error: invalid JSON request",
		req:        `!{s"d username"2:"", "email":""fj, "password":""(}`,
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrInvalidJSON.Error()},
	},
	{
		name:       "error: user already exist",
		req:        `{"username":"test-user", "email":"test@gmail.com", "password":"test"}`,
		mockInput:  &domain.CreateUserRequest{Username: "test-user", Email: "test@gmail.com", Password: "test"},
		mockErr:    domain.ErrUserExist,
		wantErr:    true,
		wantStatus: 409,
		wantResp:   domain.ErrorResponse{Message: domain.ErrUserExist.Error()},
	},
}

func TestUserHandler_Create(t *testing.T) {
	for _, tt := range testsCreate {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			mockSvc.On("CreateUser", mock.Anything, tt.mockInput).Return(tt.mockResp, tt.mockErr)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/users/register", strings.NewReader(tt.req))

			userHandler := NewUserHandler(mockSvc)

			userHandler.CreateUser(w, req)

			if tt.wantErr {
				var errResp domain.ErrorResponse
				err := json.NewDecoder(w.Body).Decode(&errResp)
				require.NoError(t, err, "failed to decode w.Body")

				assert.Equal(t, tt.wantResp, errResp)
			} else {
				var resp *domain.CreateUserResponse
				err := json.NewDecoder(w.Body).Decode(&resp)
				require.NoError(t, err, "failed to decode w.Body")

				assert.Equal(t, tt.resp, resp)
			}
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

type TestUpdate struct {
	name       string
	url        string
	idStr      string
	req        string
	resp       *domain.UpdateUserResponse
	mockID     int
	mockInput  *domain.UpdateUserRequest
	mockResp   *domain.User
	mockErr    error
	wantErr    bool
	wantStatus int
	wantResp   error
}

var testsUpdate = []TestUpdate{
	{
		name:   "success",
		url:    "/users/{id}",
		idStr:  "1",
		req:    `{"username":"UPD-test-user", "email":"upd-test@gmail.com", "oldPassword":"12345678", "newPassword":"new_password123"}`,
		resp:   &domain.UpdateUserResponse{Username: "UPD-test-user", Email: "upd-test@gmail.com"},
		mockID: 1,
		mockInput: &domain.UpdateUserRequest{
			Username:    "UPD-test-user",
			Email:       "upd-test@gmail.com",
			Password:    "12345678",
			NewPassword: "new_password123",
		},
		mockResp:   &domain.User{ID: 1, Username: "UPD-test-user", Email: "upd-test@gmail.com", PasswordHash: "new_passwordhash123"},
		wantErr:    false,
		wantStatus: 200,
	},
	{
		name:       "error: empty id",
		url:        "/users/{id}",
		idStr:      "",
		req:        `{}`,
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrIDRequired.Error()},
	},
	{
		name:       "error: invalid id",
		url:        "/users/{id}",
		idStr:      "1.5349",
		req:        `{}`,
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrInvalidID.Error()},
	},
	{
		name:       "error: invalid JSON request",
		url:        "/users/{id}",
		idStr:      "0",
		req:        `2{dfg{(}Ac2d:}`,
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrInvalidJSON.Error()},
	},
	{
		name:       "error: user not found",
		url:        "/users/{id}",
		idStr:      "1",
		mockID:     1,
		mockInput:  &domain.UpdateUserRequest{},
		mockErr:    domain.ErrUserNotFound,
		req:        `{}`,
		wantErr:    true,
		wantStatus: 404,
		wantResp:   domain.ErrorResponse{Message: domain.ErrUserNotFound.Error()},
	},
}

func TestUserHandler_Update(t *testing.T) {
	for _, tt := range testsUpdate {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			mockSvc.On("UpdateUser", mock.Anything, tt.mockID, tt.mockInput).Return(tt.mockResp, tt.mockErr)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", tt.url, strings.NewReader(tt.req))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.idStr)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			userHandler := NewUserHandler(mockSvc)

			userHandler.UpdateUser(w, req)

			if tt.wantErr {
				var errResp domain.ErrorResponse
				err := json.NewDecoder(w.Body).Decode(&errResp)
				require.NoError(t, err, "failed to decode w.Body")

				assert.Equal(t, tt.wantResp, errResp)
			} else {
				var resp domain.UpdateUserResponse
				err := json.NewDecoder(w.Body).Decode(&resp)
				require.NoError(t, err, "failed to decode w.Body")

				assert.Equal(t, tt.resp, &resp)
			}
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

type TestGet struct {
	name        string
	idStr       string
	url         string
	resp        *domain.GetUserResponse
	mockID      int
	mockResp    *domain.User
	mockErr     error
	skipService bool
	wantErr     bool
	wantStatus  int
	wantResp    error
}

var testsGet = []TestGet{
	{
		name:  "success",
		idStr: "1",
		url:   "/users/{id}",
		resp: &domain.GetUserResponse{
			Username: "test-user",
		},
		mockID: 1,
		mockResp: &domain.User{
			ID:       1,
			Username: "test-user",
		},
		wantErr:    false,
		wantStatus: 200,
	},
	{
		name:       "error: empty id",
		idStr:      "",
		mockErr:    domain.ErrIDRequired,
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrIDRequired.Error()},
	},
	{
		name:       "error: invalid id",
		idStr:      "1.5349",
		mockErr:    domain.ErrInvalidID,
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrInvalidID.Error()},
	},
	{
		name:       "error: user not found",
		idStr:      "1",
		mockID:     1,
		mockErr:    domain.ErrUserNotFound,
		wantErr:    true,
		wantStatus: 404,
		wantResp:   domain.ErrorResponse{Message: domain.ErrUserNotFound.Error()},
	},
}

func TestUserHandler_Get(t *testing.T) {
	for _, tt := range testsGet {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &MockService{}
			mockSvc.On("GetUser", mock.Anything, tt.mockID).Return(tt.mockResp, tt.mockErr)

			userHandler := NewUserHandler(mockSvc)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/users/{id}", nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.idStr)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			userHandler.GetUser(w, req)

			if tt.wantErr {
				var errResp domain.ErrorResponse
				err := json.NewDecoder(w.Body).Decode(&errResp)
				require.NoError(t, err, "failed to decode w.Body")

				assert.Equal(t, tt.wantResp, errResp)
			} else {
				var resp domain.GetUserResponse
				err := json.NewDecoder(w.Body).Decode(&resp)
				require.NoError(t, err, "failed to decode w.Body")

				assert.Equal(t, tt.resp, &resp)
			}
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

type TestList struct {
	name           string
	url            string
	mockCursor     int
	mockLimit      uint64
	mockResp       []*domain.User
	mockPagination *domain.Pagination
	mockErr        error
	resp           domain.ListUsersResponse
	skipService    bool
	wantErr        bool
	wantStatus     int
	wantResp       error
}

var testsList = []TestList{
	{
		name:       "success",
		url:        "/users?cursor=1&limit=2",
		mockCursor: 1,
		mockLimit:  2,
		mockResp: []*domain.User{
			{
				ID:           1,
				Username:     "test-user1",
				Email:        "test1@gmail.com",
				PasswordHash: "test-passwordhash1",
				Role:         "user",
			},
			{
				ID:           2,
				Username:     "test-user2",
				Email:        "test2@gmail.com",
				PasswordHash: "test-passwordhash2",
				Role:         "admin",
			},
		},
		mockPagination: &domain.Pagination{
			NextCursor: 3,
			HasMore:    true,
			Limit:      2,
		},
		resp: domain.ListUsersResponse{
			Users: []*domain.User{
				{
					ID:           1,
					Username:     "test-user1",
					Email:        "test1@gmail.com",
					PasswordHash: "test-passwordhash1",
					Role:         "user",
				},
				{
					ID:           2,
					Username:     "test-user2",
					Email:        "test2@gmail.com",
					PasswordHash: "test-passwordhash2",
					Role:         "admin",
				},
			},
			Pagination: &domain.Pagination{
				NextCursor: 3,
				HasMore:    true,
				Limit:      2,
			},
		},
		wantErr:    false,
		wantStatus: 200,
	},
	{
		name: "error: users not found",
	},
	{
		name:       "error: invalid cursor",
		url:        "/users?cursor=1.5349&limit=1",
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrInvalidCursor.Error()},
	},
	{
		name:       "error: invalid limit",
		url:        "/users?cursor=1&limit=-1.256",
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrInvalidLimit.Error()},
	},
	{
		name:       "error: internal server error",
		url:        "/users?cursor=1&limit=1",
		mockCursor: 1,
		mockLimit:  1,
		mockErr:    errors.New("internal server error"),
		wantErr:    true,
		wantStatus: 500,
		wantResp:   domain.ErrorResponse{Message: "internal server error"},
	},
}

func TestUserHandler_List(t *testing.T) {
	for _, tt := range testsList {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			mockSvc.On("ListUsers", mock.Anything, tt.mockCursor, tt.mockLimit).Return(tt.mockResp, tt.mockPagination, tt.mockErr)

			userHandler := NewUserHandler(mockSvc)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", tt.url, nil)

			userHandler.ListUsers(w, req)

			if tt.wantErr {
				var errResp domain.ErrorResponse
				err := json.NewDecoder(w.Body).Decode(&errResp)
				require.NoError(t, err, "failed to decode w.Body")

				assert.Equal(t, tt.wantResp, errResp)
			} else {
				var resp domain.ListUsersResponse
				err := json.NewDecoder(w.Body).Decode(&resp)
				require.NoError(t, err, "failed to decode w.Body")

				assert.Equal(t, tt.resp, resp)
			}
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

type TestDelete struct {
	name        string
	url         string
	idStr       string
	mockID      int
	mockError   error
	skipService bool
	wantErr     bool
	wantStatus  int
	wantResp    error
}

var testsDelete = []TestDelete{
	{
		name:       "success",
		url:        "/users/{id}",
		idStr:      "1",
		mockID:     1,
		mockError:  nil,
		wantErr:    false,
		wantStatus: 204,
		wantResp:   nil,
	},
	{
		name:       "error: empty id",
		url:        "/users/{id}",
		idStr:      "",
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrIDRequired.Error()},
	},
	{
		name:       "error: invalid id",
		url:        "/users/{id}",
		idStr:      "1.5439",
		wantErr:    true,
		wantStatus: 400,
		wantResp:   domain.ErrorResponse{Message: domain.ErrInvalidID.Error()},
	},
	{
		name:       "error: user not found",
		url:        "/users/{id}",
		idStr:      "1",
		mockID:     1,
		mockError:  domain.ErrUserNotFound,
		wantErr:    true,
		wantStatus: 404,
		wantResp:   domain.ErrorResponse{Message: domain.ErrUserNotFound.Error()},
	},
}

func TestUserHandler_Delete(t *testing.T) {
	for _, tt := range testsDelete {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockService)
			mockSvc.On("DeleteUser", mock.Anything, tt.mockID).Return(tt.mockError)

			userHandler := NewUserHandler(mockSvc)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", tt.url, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.idStr)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			userHandler.DeleteUser(w, req)

			if tt.wantErr {
				var errResp domain.ErrorResponse
				err := json.NewDecoder(w.Body).Decode(&errResp)
				require.NoError(t, err)

				assert.Equal(t, tt.wantResp, errResp)

			}
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
