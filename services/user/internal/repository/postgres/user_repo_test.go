package postgres

import (
	"apigateway/services/user/internal/database"
	"apigateway/services/user/internal/domain"
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *UserRepository {
	dbUser := os.Getenv("USER_POSTGRES_DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("USER_POSTGRES_DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "test"
	}
	dbHost := os.Getenv("USER_DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("USER_DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbName := os.Getenv("USER_DB_NAME")
	if dbName == "" {
		dbName = "users_pool"
	}

	var logger *slog.Logger

	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	pool, err := database.NewPgxPool(connStr, logger)
	require.NoError(t, err)

	m, err := migrate.New("file://../../../migrations", connStr)
	require.NoError(t, err)

	m.Up()

	if err = pool.Ping(context.Background()); err != nil {
		t.Fatal("failed to connect database: ", err)
	}

	t.Cleanup(func() {
		pool.Exec(context.Background(), "TRUNCATE users CASCADE")
		pool.Exec(context.Background(), "ALTER SEQUENCE users_id_seq RESTART WITH 1")
		pool.Close()
	})

	return NewUserRepository(pool)
}

var testUser = map[string]any{
	"username":      "test-user",
	"email":         "test-email@gmail.com",
	"password_hash": "fake-password-hash123",
	"role_id":       3,
}

func CreateTestUser(repo UserRepoInterface, t *testing.T) *domain.User {
	user, err := repo.CreateUser(context.Background(), testUser)
	require.NoError(t, err)
	return user
}

type TestCreate struct {
	name string
	want *domain.User
}

var testCreate = TestCreate{
	name: "success",
	want: &domain.User{
		ID:           1,
		Username:     "test-user",
		Email:        "test-email@gmail.com",
		PasswordHash: "fake-password-hash123",
		Role:         "user",
	},
}

func TestUserRepository_Create(t *testing.T) {
	repo := setupTestDB(t)
	t.Run(testCreate.name, func(t *testing.T) {
		user, err := repo.CreateUser(context.Background(), testUser)
		testCreate.want.CreatedAt = user.CreatedAt
		testCreate.want.UpdatedAt = user.UpdatedAt
		require.NoError(t, err)

		assert.Equal(t, testCreate.want, user)
	})
}

type TestUpdate struct {
	name       string
	userID     int
	updateData map[string]any
	want       *domain.User
}

var testUpdate = TestUpdate{
	name:   "success",
	userID: 1,
	updateData: map[string]any{
		"username":      "UPD-test-user",
		"email":         "UPD-test-email@gmail.com",
		"password_hash": "UPD-fake-password-hash123",
	},
	want: &domain.User{
		ID:           1,
		Username:     "UPD-test-user",
		Email:        "UPD-test-email@gmail.com",
		PasswordHash: "UPD-fake-password-hash123",
		Role:         "user",
	},
}

func TestUserRepository_Update(t *testing.T) {
	repo := setupTestDB(t)
	CreateTestUser(repo, t)
	t.Run(testUpdate.name, func(t *testing.T) {
		user, err := repo.UpdateUser(context.Background(), testUpdate.userID, testUpdate.updateData)
		testUpdate.want.CreatedAt = user.CreatedAt
		testUpdate.want.UpdatedAt = user.UpdatedAt
		require.NoError(t, err)

		assert.Equal(t, testUpdate.want, user)
	})
}

type TestGet struct {
	name string
}

var testGet = TestGet{
	name: "success",
}

func TestUserRepository_Get(t *testing.T) {
	repo := setupTestDB(t)
	expected := CreateTestUser(repo, t)
	t.Run(testGet.name, func(t *testing.T) {
		user, err := repo.GetUser(context.Background(), expected.ID)
		require.NoError(t, err)

		assert.Equal(t, expected, user)
	})
}

type TestList struct {
	name    string
	cursor  int
	limit   uint64
	len     int
	hasMore bool
}

var testList = TestList{
	name:    "success",
	cursor:  1,
	limit:   2,
	len:     3,
	hasMore: true,
}

func TestUserRepository_List(t *testing.T) {
	repo := setupTestDB(t)

	var createdUsers []*domain.User
	for range testList.len {
		u := CreateTestUser(repo, t)
		createdUsers = append(createdUsers, u)
	}

	t.Run(testList.name, func(t *testing.T) {
		users, nextCursor, hasMore, err := repo.ListUsers(context.Background(), testList.cursor, testList.limit)
		require.NoError(t, err)

		assert.Equal(t, createdUsers[:testList.limit], users)
		assert.Equal(t, createdUsers[testList.len-1].ID, nextCursor)
		assert.Equal(t, testList.hasMore, hasMore)
	})
}

type TestDelete struct {
	name   string
	userID int
	want   error
}

var testDelete = TestDelete{
	name:   "success",
	userID: 1,
	want:   nil,
}

func TestUserRepository_Delete(t *testing.T) {
	repo := setupTestDB(t)
	CreateTestUser(repo, t)

	t.Run(testDelete.name, func(t *testing.T) {
		err := repo.DeleteUser(context.Background(), testDelete.userID)
		require.NoError(t, err)
	})
}
