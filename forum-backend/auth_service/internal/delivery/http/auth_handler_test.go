package http

import (
	"bytes"
	"encoding/json"
	"errors"
	utils "github.com/Engls/EnglsJwt"
	"github.com/Engls/forum-project2/auth_service/internal/entity"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/Engls/forum-project2/auth_service/internal/usecase"
	"github.com/Engls/forum-project2/auth_service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestAuthHandler_Register_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание мока для AuthUsecase
	mockAuthUsecase := new(mocks.AuthUsecase)
	mockAuthUsecase.On("Register", "testuser", "password", "user").Return(nil)

	// Инициализация обработчиков
	authHandler := NewAuthHandler(mockAuthUsecase, nil, logger)

	// Создание тестового маршрутизатора
	router := gin.Default()
	router.POST("/register", authHandler.Register)

	// Создание тестового запроса
	reqBody := map[string]string{
		"username": "testuser",
		"password": "password",
		"role":     "user",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверка результата
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User registered successfully")

	// Проверка вызовов мока
	mockAuthUsecase.AssertExpectations(t)
}

func TestAuthHandler_Register_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockAuthUsecase := new(mocks.AuthUsecase)
	jwtUtil := utils.NewJWTUtil("secret")

	// Тестовые данные
	req := entity.RegisterRequest{
		Username: "testuser",
		Password: "password",
		Role:     "user",
	}

	// Настройка моков
	mockAuthUsecase.On("Register", req.Username, req.Password, req.Role).Return(errors.New("failed to register user"))

	// Инициализация хендлера
	authHandler := NewAuthHandler(mockAuthUsecase, jwtUtil, logger)

	// Создание Gin контекста
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{"username":"testuser","password":"password","role":"user"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	// Вызов метода Register
	authHandler.Register(c)

	// Проверка результата
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to register user")

	// Проверка вызовов моков
	mockAuthUsecase.AssertExpectations(t)
}

func TestAuthHandler_Register_BadRequest(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockAuthUsecase := new(mocks.AuthUsecase)
	jwtUtil := utils.NewJWTUtil("secret")

	// Инициализация хендлера
	authHandler := NewAuthHandler(mockAuthUsecase, jwtUtil, logger)

	// Создание Gin контекста
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{"{"username":"testuser","password":"","role":""}`))
	c.Request.Header.Set("Content-Type", "application/json")

	// Вызов метода Register
	authHandler.Register(c)

	// Проверка результата
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")

	// Проверка вызовов моков
	mockAuthUsecase.AssertExpectations(t)
}

func TestAuthHandler_Login(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание мока для AuthUsecase
	mockAuthUsecase := new(mocks.AuthUsecase)
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo3LCJyb2xlIjoidXNlciIsImV4cCI6MTc0NTUxNzExOCwiaWF0IjoxNzQ1MjU3OTE4fQ.QmyaHsq-ruAyciGKkgCEgj0xsQZD1J5ER6CLjXhfgQc"
	mockAuthUsecase.On("Login", "testuser", "password").Return(token, nil)
	mockAuthUsecase.On("GetUserRole", "testuser").Return("user", nil)

	// Инициализация обработчиков
	jwtUtil := utils.NewJWTUtil("your-secret-key")
	authHandler := NewAuthHandler(mockAuthUsecase, jwtUtil, logger)

	// Создание тестового маршрутизатора
	router := gin.Default()
	router.POST("/login", authHandler.Login)

	// Создание тестового запроса
	reqBody := map[string]string{
		"username": "testuser",
		"password": "password",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверка результата
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")

	// Проверка вызовов мока
	mockAuthUsecase.AssertExpectations(t)
}

func TestAuthHandler_Login_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockAuthUsecase := new(mocks.AuthUsecase)
	jwtUtil := utils.NewJWTUtil("secret")

	// Тестовые данные
	req := entity.LoginRequest{
		Username: "testuser",
		Password: "password",
	}

	// Настройка моков
	mockAuthUsecase.On("Login", req.Username, req.Password).Return("", errors.New("invalid credentials"))

	// Инициализация хендлера
	authHandler := NewAuthHandler(mockAuthUsecase, jwtUtil, logger)

	// Создание Gin контекста
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`{"username":"testuser","password":"password"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	// Вызов метода Login
	authHandler.Login(c)

	// Проверка результата
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid credentials")

	// Проверка вызовов моков
	mockAuthUsecase.AssertExpectations(t)
}

func TestAuthHandler_Login_GetUserIDFromTokenFailure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockAuthUsecase := new(mocks.AuthUsecase)
	jwtUtil := utils.NewJWTUtil("secret")

	// Тестовые данные
	req := entity.LoginRequest{
		Username: "testuser",
		Password: "password",
	}
	token := "invalid.jwt.token"
	role := "user"

	// Настройка моков
	mockAuthUsecase.On("Login", req.Username, req.Password).Return(token, nil)
	mockAuthUsecase.On("GetUserRole", req.Username).Return(role, nil)

	// Инициализация хендлера
	authHandler := NewAuthHandler(mockAuthUsecase, jwtUtil, logger)

	// Создание Gin контекста
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`{"username":"testuser","password":"password"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	// Вызов метода Login
	authHandler.Login(c)

	// Проверка результата
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "invalid character")

	// Проверка вызовов моков
	mockAuthUsecase.AssertExpectations(t)
}
