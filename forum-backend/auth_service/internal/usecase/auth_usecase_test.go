package usecase

import (
	"errors"
	"testing"

	"github.com/Engls/EnglsJwt"
	"github.com/Engls/forum-project2/auth_service/internal/entity"
	"github.com/Engls/forum-project2/auth_service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthUsecase_Register_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockAuthRepo := new(mocks.AuthRepository)
	jwtUtil := EnglsJwt.NewJWTUtil("secret")

	// Тестовые данные
	username := "testuser"
	password := "password"
	role := "user"

	// Настройка моков
	mockAuthRepo.On("Register", mock.AnythingOfType("entity.User")).Return(nil)

	// Инициализация usecase
	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	// Вызов метода Register
	err := authUsecase.Register(username, password, role)

	// Проверка результата
	assert.NoError(t, err)

	// Проверка вызовов моков
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUsecase_Register_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockAuthRepo := new(mocks.AuthRepository)
	jwtUtil := EnglsJwt.NewJWTUtil("secret")

	// Тестовые данные
	username := "testuser"
	password := "password"
	role := "user"

	// Настройка моков
	mockAuthRepo.On("Register", mock.AnythingOfType("entity.User")).Return(errors.New("failed to register user"))

	// Инициализация usecase
	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	// Вызов метода Register
	err := authUsecase.Register(username, password, role)

	// Проверка результата
	assert.Error(t, err)

	// Проверка вызовов моков
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUsecase_Login_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockAuthRepo := new(mocks.AuthRepository)
	jwtUtil := EnglsJwt.NewJWTUtil("secret")

	// Тестовые данные
	username := "testuser"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := entity.User{ID: 1, Username: username, Password: string(hashedPassword), Role: "user"}

	// Настройка моков
	mockAuthRepo.On("GetUserByUsername", username).Return(user, nil)
	mockAuthRepo.On("SaveToken", user.ID, mock.Anything).Return(nil)

	// Инициализация usecase
	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	// Вызов метода Login
	resultToken, err := authUsecase.Login(username, password)

	// Проверка результата
	assert.NoError(t, err)
	assert.NotEmpty(t, resultToken)

	// Проверка вызовов моков
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUsecase_Login_Failure_InvalidCredentials(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockAuthRepo := new(mocks.AuthRepository)
	jwtUtil := EnglsJwt.NewJWTUtil("secret")

	// Тестовые данные
	username := "testuser"
	password := "password"

	// Настройка моков
	mockAuthRepo.On("GetUserByUsername", username).Return(entity.User{}, errors.New("user not found"))

	// Инициализация usecase
	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	// Вызов метода Login
	resultToken, err := authUsecase.Login(username, password)

	// Проверка результата
	assert.Error(t, err)
	assert.Equal(t, "", resultToken)

	// Проверка вызовов моков
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUsecase_Login_Failure_InvalidPassword(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockAuthRepo := new(mocks.AuthRepository)
	jwtUtil := EnglsJwt.NewJWTUtil("secret")

	// Тестовые данные
	username := "testuser"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("wrongpassword"), bcrypt.DefaultCost)
	user := entity.User{ID: 1, Username: username, Password: string(hashedPassword), Role: "user"}

	// Настройка моков
	mockAuthRepo.On("GetUserByUsername", username).Return(user, nil)

	// Инициализация usecase
	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	// Вызов метода Login
	resultToken, err := authUsecase.Login(username, password)

	// Проверка результата
	assert.Error(t, err)
	assert.Equal(t, "", resultToken)

	// Проверка вызовов моков
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUsecase_GetUserRole_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockAuthRepo := new(mocks.AuthRepository)
	jwtUtil := EnglsJwt.NewJWTUtil("secret")

	// Тестовые данные
	username := "testuser"
	user := entity.User{ID: 1, Username: username, Password: "hashedpassword", Role: "user"}

	// Настройка моков
	mockAuthRepo.On("GetUserByUsername", username).Return(user, nil)

	// Инициализация usecase
	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	// Вызов метода GetUserRole
	role, err := authUsecase.GetUserRole(username)

	// Проверка результата
	assert.NoError(t, err)
	assert.Equal(t, user.Role, role)

	// Проверка вызовов моков
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUsecase_GetUserRole_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockAuthRepo := new(mocks.AuthRepository)
	jwtUtil := EnglsJwt.NewJWTUtil("secret")

	// Тестовые данные
	username := "testuser"

	// Настройка моков
	mockAuthRepo.On("GetUserByUsername", username).Return(entity.User{}, errors.New("user not found"))

	// Инициализация usecase
	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	// Вызов метода GetUserRole
	role, err := authUsecase.GetUserRole(username)

	// Проверка результата
	assert.Error(t, err)
	assert.Equal(t, "", role)

	// Проверка вызовов моков
	mockAuthRepo.AssertExpectations(t)
}
