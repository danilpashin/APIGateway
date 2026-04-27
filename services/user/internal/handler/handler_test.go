package handler

import (
	"apigateway/services/user/internal/domain"
	"apigateway/services/user/internal/service"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

type MockUserRepo struct {
	createUser func(ctx context.Context, insertData map[string]any) (*domain.User, error)
	updateUser func(ctx context.Context, id int, updateData map[string]any) (*domain.User, error)
	getUser    func(ctx context.Context, id int) (*domain.User, error)
	listUsers  func(ctx context.Context, cursor int, limit uint64) ([]*domain.User, int, bool, error)
	deleteUser func(ctx context.Context, id int) error
}

func (m *MockUserRepo) CreateUser(ctx context.Context, insertData map[string]any) (*domain.User, error) {
	if m.createUser != nil {
		return m.createUser(ctx, insertData)
	}
	return nil, nil
}

func (m *MockUserRepo) UpdateUser(ctx context.Context, id int, updateData map[string]any) (*domain.User, error) {
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

type TestCreate struct {
	name       string
	user       *domain.User
	req        string
	resp       *domain.CreateUserResponse
	wantErr    bool
	wantStatus int
	wantResp   string
}

var testsCreate = []TestCreate{
	{
		name:       "general",
		user:       &domain.User{ID: 1, Username: "Danil132", Email: "rvn243@gmail.com", PasswordHash: "", Role: "user"},
		req:        `{"username":"Danil132", "email":"rvn243@gmail.com", "password":"test"}`,
		resp:       &domain.CreateUserResponse{Username: "Danil132", Email: "rvn243@gmail.com"},
		wantErr:    false,
		wantStatus: 201,
	},
	{
		name:       "empty insert data",
		user:       nil,
		req:        `{"username":"", "email":"", "password":""}`,
		wantErr:    true,
		wantStatus: 400,
		wantResp:   "insert data is empty or not enough\n",
	},
	{
		name:       "invalid JSON request",
		user:       nil,
		req:        `!{s"d username"2:"", "email":""fj, "password":""(}`,
		wantErr:    true,
		wantStatus: 400,
		wantResp:   "invalid JSON\n",
	},
}

func TestCreateUser(t *testing.T) {
	for _, test := range testsCreate {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := MockUserRepo{
				createUser: func(ctx context.Context, insertData map[string]any) (*domain.User, error) {
					return test.user, nil
				},
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/users/register", strings.NewReader(test.req))

			userService := service.NewUserService(&mockRepo)
			userHandler := NewUserHandler(*userService)

			userHandler.CreateUser(w, req)

			if test.wantErr {
				bodyBytes, err := io.ReadAll(w.Body)
				if err != nil {
					t.Fatal("failed to read w.Body: ", err)
				}

				got := string(bodyBytes)
				if !strings.EqualFold(got, test.wantResp) {
					t.Fatalf("expected %v, got: %v", test.wantResp, got)
				}

				if w.Code != test.wantStatus {
					t.Fatalf("expected %d, got: %d", test.wantStatus, w.Code)
				}
			} else {
				var resp *domain.CreateUserResponse
				err := json.NewDecoder(w.Body).Decode(&resp)
				if err != nil {
					t.Fatal("failed to decode w.Body: ", err)
				}

				if !reflect.DeepEqual(resp, test.resp) {
					t.Fatalf("expected %v, got: %v", test.resp, resp)
				}

				if w.Code != test.wantStatus {
					t.Fatalf("expected %d, got: %d", test.wantStatus, w.Code)
				}
			}
		})
	}
}

type TestUpdate struct {
	name       string
	user       *domain.User
	userID     string
	req        string
	resp       *domain.UpdateUserResponse
	wantErr    bool
	wantStatus int
	wantResp   string
}

var testsUpdate = []TestUpdate{
	{
		name:       "general",
		user:       &domain.User{ID: 1, Username: "Danil132", Email: "rvn243@gmail.com", PasswordHash: "$2a$10$nqFp/wGlchdjqATC22vgguUXzY.lXUoyizMsYeD8GjpG48bBk5tpe", Role: "user"},
		userID:     "1",
		req:        `{"username":"sad31fd", "email":"another_email@mail.ru", "oldPassword":"12345678", "newPassword":"new_password123"}`,
		resp:       &domain.UpdateUserResponse{Username: "sad31fd", Email: "another_email@mail.ru"},
		wantErr:    false,
		wantStatus: 200,
	},
	{
		name:       "invalid JSON",
		user:       nil,
		userID:     "0",
		req:        `2{dfg{(}Ac2d:}`,
		wantErr:    true,
		wantStatus: 400,
		wantResp:   "invalid JSON\n",
	},
}

func TestUpdateUser(t *testing.T) {
	for _, test := range testsUpdate {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := MockUserRepo{
				updateUser: func(ctx context.Context, id int, updateData map[string]any) (*domain.User, error) {
					username := updateData["username"].(string)
					email := updateData["email"].(string)
					passwordHash := updateData["password_hash"].([]byte)
					test.user.Username = username
					test.user.Email = email
					test.user.PasswordHash = string(passwordHash)

					return test.user, nil
				},
				getUser: func(ctx context.Context, id int) (*domain.User, error) {
					return test.user, nil
				},
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/users/{id}", strings.NewReader(test.req))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", test.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			userService := service.NewUserService(&mockRepo)
			userHandler := NewUserHandler(*userService)

			userHandler.UpdateUser(w, req)

			if test.wantErr {
				bodyBytes, err := io.ReadAll(w.Body)
				if err != nil {
					t.Fatal("failed to read w.Body: ", err)
				}

				got := string(bodyBytes)
				if !strings.EqualFold(got, test.wantResp) {
					t.Fatalf("expected %s, got: %s", test.wantResp, got)
				}

				if w.Code != test.wantStatus {
					t.Fatalf("expected %d, got: %d", test.wantStatus, w.Code)
				}
			} else {
				var resp domain.UpdateUserResponse
				err := json.NewDecoder(w.Body).Decode(&resp)
				if err != nil {
					t.Fatal("failed to decode w.Body: ", err)
				}

				if !reflect.DeepEqual(&resp, test.resp) {
					t.Fatalf("expected %v, got: %v", test.resp, &resp)
				}

				if w.Code != test.wantStatus {
					t.Fatalf("expected %d, got: %d", test.wantStatus, w.Code)
				}
			}
		})
	}
}
