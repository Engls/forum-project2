package http

import (
	"bytes"
	"encoding/json"
	"errors"
	utils "github.com/Engls/EnglsJwt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"github.com/Engls/forum-project2/forum_service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestPostHandler_CreatePost_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	// Инициализация обработчика
	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger)

	// Тестовые данные
	post := &entity.Post{
		Title:   "Test Post",
		Content: "This is a test post",
	}
	postJSON, _ := json.Marshal(post)

	// Настройка моков
	mockPostRepo.On("GetUserIDByToken", mock.Anything, "valid.jwt.token").Return(1, nil)
	mockPostUsecase.On("CreatePost", mock.Anything, mock.Anything).Return(post, nil)

	// Создание тестового запроса
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(postJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")

	// Выполнение запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Вызов метода
	postHandler.CreatePost(c)

	// Проверка результата
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Test Post")

	// Проверка вызовов моков
	mockPostRepo.AssertExpectations(t)
	mockPostUsecase.AssertExpectations(t)
}

func TestPostHandler_CreatePost_MissingAuthorizationHeader(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	// Инициализация обработчика
	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger)

	// Тестовые данные
	post := entity.Post{
		Title:   "Test Post",
		Content: "This is a test post",
	}
	postJSON, _ := json.Marshal(post)

	// Создание тестового запроса
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(postJSON))
	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Вызов метода
	postHandler.CreatePost(c)

	// Проверка результата
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header required")
}

func TestPostHandler_CreatePost_InvalidAuthorizationHeaderFormat(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	// Инициализация обработчика
	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger)

	// Тестовые данные
	post := entity.Post{
		Title:   "Test Post",
		Content: "This is a test post",
	}
	postJSON, _ := json.Marshal(post)

	// Создание тестового запроса
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(postJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "InvalidFormat")

	// Выполнение запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Вызов метода
	postHandler.CreatePost(c)

	// Проверка результата
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid Authorization header format")
}

func TestPostHandler_CreatePost_InvalidToken(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	// Инициализация обработчика
	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger)

	// Тестовые данные
	post := entity.Post{
		Title:   "Test Post",
		Content: "This is a test post",
	}
	postJSON, _ := json.Marshal(post)

	// Настройка моков
	mockPostRepo.On("GetUserIDByToken", mock.Anything, "invalid.jwt.token\t").Return(0, errors.New("invalid token"))

	// Создание тестового запроса
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(postJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid.jwt.token	")

	// Выполнение запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Вызов метода
	postHandler.CreatePost(c)

	// Проверка результата
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")

	// Проверка вызовов моков
	mockPostRepo.AssertExpectations(t)
}

func TestPostHandler_CreatePost_FailedToCreatePost(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	// Инициализация обработчика
	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger)

	// Тестовые данные
	post := &entity.Post{
		Title:   "Test Post",
		Content: "This is a test post",
	}
	postJSON, _ := json.Marshal(post)

	// Настройка моков
	mockPostRepo.On("GetUserIDByToken", mock.Anything, "valid.jwt.token").Return(1, nil)
	mockPostUsecase.On("CreatePost", mock.Anything, mock.Anything).Return(post, errors.New("failed to create post"))

	// Создание тестового запроса
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(postJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")

	// Выполнение запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Вызов метода
	postHandler.CreatePost(c)

	// Проверка результата
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to create post")

	// Проверка вызовов моков
	mockPostRepo.AssertExpectations(t)
	mockPostUsecase.AssertExpectations(t)
}

func TestPostHandler_GetPosts_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	// Инициализация обработчика
	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger)

	// Тестовые данные
	posts := []entity.Post{
		{ID: 1, Title: "Post 1", Content: "Content 1"},
		{ID: 2, Title: "Post 2", Content: "Content 2"},
	}

	// Настройка моков
	mockPostRepo.On("GetPosts", mock.Anything).Return(posts, nil)

	// Создание тестового запроса
	req, _ := http.NewRequest("GET", "/posts", nil)
	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Вызов метода
	postHandler.GetPosts(c)

	// Проверка результата
	assert.Equal(t, http.StatusOK, w.Code)
	var responsePosts []entity.Post
	err := json.Unmarshal(w.Body.Bytes(), &responsePosts)
	assert.NoError(t, err)
	assert.Equal(t, posts, responsePosts)

	// Проверка вызовов моков
	mockPostRepo.AssertExpectations(t)
}

func TestPostHandler_GetPosts_FailedToGetPosts(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	// Инициализация обработчика
	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger)

	// Настройка моков
	mockPostRepo.On("GetPosts", mock.Anything).Return(nil, errors.New("failed to get posts"))

	// Создание тестового запроса
	req, _ := http.NewRequest("GET", "/posts", nil)
	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Вызов метода
	postHandler.GetPosts(c)

	// Проверка результата
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to get posts")

	// Проверка вызовов моков
	mockPostRepo.AssertExpectations(t)
}

func TestPostHandler_DeletePost_MissingAuthorizationHeader(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	// Инициализация обработчика
	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger)

	// Создание тестового запроса
	req, _ := http.NewRequest("DELETE", "/posts/1", nil)
	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Вызов метода
	postHandler.DeletePost(c)

	// Проверка результата
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header required")
}

func TestPostHandler_DeletePost_InvalidAuthorizationHeaderFormat(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	// Инициализация обработчика
	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger)

	// Создание тестового запроса
	req, _ := http.NewRequest("DELETE", "/posts/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "InvalidFormat")

	// Выполнение запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Вызов метода
	postHandler.DeletePost(c)

	// Проверка результата
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid Authorization header format")
}

func TestPostHandler_DeletePost_Success_Owner(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	// Инициализация обработчика
	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger)

	// Генерация токена для пользователя
	token, err := jwtUtil.GenerateToken(1, "user")
	assert.NoError(t, err)

	// Настройка моков
	mockPostRepo.On("GetPostByID", mock.Anything, 1).Return(&entity.Post{ID: 1, AuthorId: 1}, nil)
	mockPostUsecase.On("DeletePost", mock.Anything, 1).Return(nil)

	// Создание тестового запроса
	req, _ := http.NewRequest("DELETE", "/posts/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Выполнение запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Вызов метода
	postHandler.DeletePost(c)

	// Проверка результата
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверка вызовов моков
	mockPostRepo.AssertExpectations(t)
	mockPostUsecase.AssertExpectations(t)
}

func TestPostHandler_DeletePost_Success_Admin(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	// Инициализация обработчика
	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger)

	// Генерация токена для админа
	token, err := jwtUtil.GenerateToken(1, "admin")
	assert.NoError(t, err)

	// Настройка моков
	mockPostRepo.On("GetPostByID", mock.Anything, 1).Return(entity.Post{ID: 1, AuthorId: 2}, nil)
	mockPostUsecase.On("DeletePost", mock.Anything, 1).Return(nil)

	// Создание тестового запроса
	req, _ := http.NewRequest("DELETE", "/posts/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Выполнение запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Вызов метода
	postHandler.DeletePost(c)

	// Проверка результата
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверка вызовов моков
	mockPostUsecase.AssertExpectations(t)
}
