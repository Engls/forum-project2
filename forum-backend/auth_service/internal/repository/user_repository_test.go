package repository

import (
	"database/sql"
	"testing"

	"github.com/Engls/forum-project2/auth_service/internal/entity"
	"github.com/Engls/forum-project2/auth_service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestAuthRepository_Register(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание мока для DB
	mockDB := new(mocks.DB)
	mockDB.On("Exec", "INSERT INTO users (username, password, role) VALUES (?, ?, ?)", "testuser", "hashedpassword", "user").Return(sql.Result(nil), nil)

	// Инициализация репозитория
	authRepo := NewAuthRepository(mockDB, logger)

	// Вызов метода Register
	user := entity.User{Username: "testuser", Password: "hashedpassword", Role: "user"}
	err := authRepo.Register(user)

	// Проверка результата
	assert.NoError(t, err)

	// Проверка вызовов мока
	mockDB.AssertExpectations(t)
}

func TestAuthRepository_GetUserByUsername(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание мока для DB
	mockDB := new(mocks.DB)
	user := entity.User{ID: 1, Username: "testuser", Password: "hashedpassword", Role: "user"}
	mockDB.On("Get", mock.Anything, "SELECT id, username, password, role FROM users WHERE username=?", "testuser").Run(func(args mock.Arguments) {
		dest := args.Get(0).(*entity.User)
		*dest = user
	}).Return(nil)

	// Инициализация репозитория
	authRepo := NewAuthRepository(mockDB, logger)

	// Вызов метода GetUserByUsername
	result, err := authRepo.GetUserByUsername("testuser")

	// Проверка результата
	assert.NoError(t, err)
	assert.Equal(t, user, result)

	// Проверка вызовов мока
	mockDB.AssertExpectations(t)
}

func TestAuthRepository_SaveToken(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание мока для DB
	mockDB := new(mocks.DB)
	mockDB.On("Exec", "INSERT INTO tokens (user_id, token) VALUES (?, ?)", 1, "valid.jwt.token").Return(sql.Result(nil), nil)

	// Инициализация репозитория
	authRepo := NewAuthRepository(mockDB, logger)

	// Вызов метода SaveToken
	err := authRepo.SaveToken(1, "valid.jwt.token")

	// Проверка результата
	assert.NoError(t, err)

	// Проверка вызовов мока
	mockDB.AssertExpectations(t)
}
