package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/Engls/forum-project2/auth_service/internal/entity"
	"github.com/Engls/forum-project2/auth_service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestAuthRepository_Register_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание мока для DB
	mockDB := new(mocks.DB)

	// Тестовые данные
	user := entity.User{Username: "testuser", Password: "hashedpassword", Role: "user"}

	// Настройка моков
	mockDB.On("Exec", "INSERT INTO users (username, password, role) VALUES (?, ?, ?)", user.Username, user.Password, user.Role).Return(sql.Result(nil), nil)

	// Инициализация репозитория
	authRepo := NewAuthRepository(mockDB, logger)

	// Вызов метода Register
	err := authRepo.Register(user)

	// Проверка результата
	assert.NoError(t, err)

	// Проверка вызовов мока
	mockDB.AssertExpectations(t)
}

func TestAuthRepository_Register_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание мока для DB
	mockDB := new(mocks.DB)

	// Тестовые данные
	user := entity.User{Username: "testuser", Password: "hashedpassword", Role: "user"}

	// Настройка моков
	mockDB.On("Exec", "INSERT INTO users (username, password, role) VALUES (?, ?, ?)", user.Username, user.Password, user.Role).Return(nil, errors.New("failed to register user"))

	// Инициализация репозитория
	authRepo := NewAuthRepository(mockDB, logger)

	// Вызов метода Register
	err := authRepo.Register(user)

	// Проверка результата
	assert.Error(t, err)

	// Проверка вызовов мока
	mockDB.AssertExpectations(t)
}

func TestAuthRepository_GetUserByUsername_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание мока для DB
	mockDB := new(mocks.DB)

	// Тестовые данные
	user := entity.User{ID: 1, Username: "testuser", Password: "hashedpassword", Role: "user"}

	// Настройка моков
	mockDB.On("Get", mock.Anything, "SELECT id, username, password, role FROM users WHERE username=?", user.Username).Run(func(args mock.Arguments) {
		dest := args.Get(0).(*entity.User)
		*dest = user
	}).Return(nil)

	// Инициализация репозитория
	authRepo := NewAuthRepository(mockDB, logger)

	// Вызов метода GetUserByUsername
	result, err := authRepo.GetUserByUsername(user.Username)

	// Проверка результата
	assert.NoError(t, err)
	assert.Equal(t, user, result)

	// Проверка вызовов мока
	mockDB.AssertExpectations(t)
}

func TestAuthRepository_GetUserByUsername_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание мока для DB
	mockDB := new(mocks.DB)

	// Тестовые данные
	username := "testuser"

	// Настройка моков
	mockDB.On("Get", mock.Anything, "SELECT id, username, password, role FROM users WHERE username=?", username).Return(errors.New("failed to get user by username"))

	// Инициализация репозитория
	authRepo := NewAuthRepository(mockDB, logger)

	// Вызов метода GetUserByUsername
	result, err := authRepo.GetUserByUsername(username)

	// Проверка результата
	assert.Error(t, err)
	assert.Equal(t, entity.User{}, result)

	// Проверка вызовов мока
	mockDB.AssertExpectations(t)
}

func TestAuthRepository_SaveToken_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание мока для DB
	mockDB := new(mocks.DB)

	// Тестовые данные
	userID := 1
	token := "valid.jwt.token"

	// Настройка моков
	mockDB.On("Exec", "INSERT INTO tokens (user_id, token) VALUES (?, ?)", userID, token).Return(sql.Result(nil), nil)

	// Инициализация репозитория
	authRepo := NewAuthRepository(mockDB, logger)

	// Вызов метода SaveToken
	err := authRepo.SaveToken(userID, token)

	// Проверка результата
	assert.NoError(t, err)

	// Проверка вызовов мока
	mockDB.AssertExpectations(t)
}

func TestAuthRepository_SaveToken_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание мока для DB
	mockDB := new(mocks.DB)

	// Тестовые данные
	userID := 1
	token := "valid.jwt.token"

	// Настройка моков
	mockDB.On("Exec", "INSERT INTO tokens (user_id, token) VALUES (?, ?)", userID, token).Return(nil, errors.New("failed to save token"))

	// Инициализация репозитория
	authRepo := NewAuthRepository(mockDB, logger)

	// Вызов метода SaveToken
	err := authRepo.SaveToken(userID, token)

	// Проверка результата
	assert.Error(t, err)

	// Проверка вызовов мока
	mockDB.AssertExpectations(t)
}
