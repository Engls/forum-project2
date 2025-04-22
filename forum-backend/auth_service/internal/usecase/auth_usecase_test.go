package usecase

import (
	utils "github.com/Engls/EnglsJwt"
	"github.com/Engls/forum-project2/auth_service/internal/entity"
	"github.com/Engls/forum-project2/auth_service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestAuthUsecase_Register(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание мока для AuthRepository
	mockAuthRepo := new(mocks.AuthRepository)
	mockAuthRepo.On("Register", mock.AnythingOfType("entity.User")).Return(nil)

	// Инициализация usecase
	jwtUtil := utils.NewJWTUtil("secret")
	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	// Вызов метода Register
	err := authUsecase.Register("testuser", "password", "user")

	// Проверка результата
	assert.NoError(t, err)

	// Проверка вызовов мока
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUsecase_Login(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание мока для AuthRepository
	mockAuthRepo := new(mocks.AuthRepository)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	user := entity.User{ID: 1, Username: "testuser", Password: string(hashedPassword), Role: "user"}
	mockAuthRepo.On("GetUserByUsername", "testuser").Return(user, nil)
	mockAuthRepo.On("SaveToken", 1, mock.AnythingOfType("string")).Return(nil)

	// Инициализация usecase
	jwtUtil := utils.NewJWTUtil("secret")
	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	// Вызов метода Login
	token, err := authUsecase.Login("testuser", "password")

	// Проверка результата
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Проверка вызовов мока
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUsecase_GetUserRole(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание мока для AuthRepository
	mockAuthRepo := new(mocks.AuthRepository)
	user := entity.User{ID: 1, Username: "testuser", Password: "hashedpassword", Role: "user"}
	mockAuthRepo.On("GetUserByUsername", "testuser").Return(user, nil)

	// Инициализация usecase
	jwtUtil := utils.NewJWTUtil("secret")
	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	// Вызов метода GetUserRole
	role, err := authUsecase.GetUserRole("testuser")

	// Проверка результата
	assert.NoError(t, err)
	assert.Equal(t, "user", role)

	// Проверка вызовов мока
	mockAuthRepo.AssertExpectations(t)
}
