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

func TestCommentHandler_CreateComment_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	// Инициализация обработчика
	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger)

	// Генерация токена для пользователя
	token, err := jwtUtil.GenerateToken(1, "user")
	assert.NoError(t, err)

	// Тестовые данные
	comment := entity.Comment{
		Content: "This is a test comment",
	}
	commentJSON, _ := json.Marshal(comment)

	// Настройка моков
	mockCommentUsecase.On("CreateComment", mock.Anything, mock.Anything).Return(comment, nil)

	// Создание тестового запроса
	req, _ := http.NewRequest("POST", "/posts/1/comments", bytes.NewBuffer(commentJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Выполнение запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Вызов метода
	commentHandler.CreateComment(c)

	// Проверка результата
	assert.Equal(t, http.StatusCreated, w.Code)
	var responseComment entity.Comment
	err = json.Unmarshal(w.Body.Bytes(), &responseComment)
	assert.NoError(t, err)
	assert.Equal(t, comment, responseComment)

	// Проверка вызовов моков
	mockCommentUsecase.AssertExpectations(t)
}

func TestCommentHandler_CreateComment_InvalidPostID(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	// Инициализация обработчика
	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger)

	// Генерация токена для пользователя
	token, err := jwtUtil.GenerateToken(1, "user")
	assert.NoError(t, err)

	// Тестовые данные
	comment := entity.Comment{
		Content: "This is a test comment",
	}
	commentJSON, _ := json.Marshal(comment)

	// Создание тестового запроса
	req, _ := http.NewRequest("POST", "/posts/invalid/comments", bytes.NewBuffer(commentJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Выполнение запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid"}}

	// Вызов метода
	commentHandler.CreateComment(c)

	// Проверка результата
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid post ID")
}

func TestCommentHandler_CreateComment_MissingAuthorizationHeader(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	// Инициализация обработчика
	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger)

	// Тестовые данные
	comment := entity.Comment{
		Content: "This is a test comment",
	}
	commentJSON, _ := json.Marshal(comment)

	// Создание тестового запроса
	req, _ := http.NewRequest("POST", "/posts/1/comments", bytes.NewBuffer(commentJSON))
	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Вызов метода
	commentHandler.CreateComment(c)

	// Проверка результата
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header required")
}

func TestCommentHandler_CreateComment_InvalidAuthorizationHeaderFormat(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	// Инициализация обработчика
	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger)

	// Тестовые данные
	comment := entity.Comment{
		Content: "This is a test comment",
	}
	commentJSON, _ := json.Marshal(comment)

	// Создание тестового запроса
	req, _ := http.NewRequest("POST", "/posts/1/comments", bytes.NewBuffer(commentJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "InvalidFormat")

	// Выполнение запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Вызов метода
	commentHandler.CreateComment(c)

	// Проверка результата
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid Authorization header format")
}

func TestCommentHandler_CreateComment_InvalidToken(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	// Инициализация обработчика
	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger)

	// Тестовые данные
	comment := entity.Comment{
		Content: "This is a test comment",
	}
	commentJSON, _ := json.Marshal(comment)

	// Создание тестового запроса
	req, _ := http.NewRequest("POST", "/posts/1/comments", bytes.NewBuffer(commentJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid.jwt.token")

	// Выполнение запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Вызов метода
	commentHandler.CreateComment(c)

	// Проверка результата
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token or user ID")
}

func TestCommentHandler_CreateComment_FailedToCreateComment(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	// Инициализация обработчика
	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger)

	// Генерация токена для пользователя
	token, err := jwtUtil.GenerateToken(1, "user")
	assert.NoError(t, err)

	// Тестовые данные
	comment := entity.Comment{
		Content: "This is a test comment",
	}
	commentJSON, _ := json.Marshal(comment)

	// Настройка моков
	mockCommentUsecase.On("CreateComment", mock.Anything, mock.Anything).Return(entity.Comment{}, errors.New("failed to create comment"))

	// Создание тестового запроса
	req, _ := http.NewRequest("POST", "/posts/1/comments", bytes.NewBuffer(commentJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Выполнение запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Вызов метода
	commentHandler.CreateComment(c)

	// Проверка результата
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to create comment")

	// Проверка вызовов моков
	mockCommentUsecase.AssertExpectations(t)
}

func TestCommentHandler_GetCommentsByPostID_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	// Инициализация обработчика
	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger)

	// Тестовые данные
	comments := []entity.Comment{
		{ID: 1, PostId: 1, Content: "Comment 1"},
		{ID: 2, PostId: 1, Content: "Comment 2"},
	}

	// Настройка моков
	mockCommentUsecase.On("GetCommentByPostID", mock.Anything, 1).Return(comments, nil)

	// Создание тестового запроса
	req, _ := http.NewRequest("GET", "/posts/1/comments", nil)
	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Вызов метода
	commentHandler.GetCommentsByPostID(c)

	// Проверка результата
	assert.Equal(t, http.StatusOK, w.Code)
	var responseComments []entity.Comment
	err := json.Unmarshal(w.Body.Bytes(), &responseComments)
	assert.NoError(t, err)
	assert.Equal(t, comments, responseComments)

	// Проверка вызовов моков
	mockCommentUsecase.AssertExpectations(t)
}

func TestCommentHandler_GetCommentsByPostID_InvalidPostID(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	// Инициализация обработчика
	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger)

	// Создание тестового запроса
	req, _ := http.NewRequest("GET", "/posts/invalid/comments", nil)
	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid"}}

	// Вызов метода
	commentHandler.GetCommentsByPostID(c)

	// Проверка результата
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid post ID")
}

func TestCommentHandler_GetCommentsByPostID_FailedToGetComments(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	// Инициализация обработчика
	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger)

	// Настройка моков
	mockCommentUsecase.On("GetCommentByPostID", mock.Anything, 1).Return(nil, errors.New("failed to get comments"))

	// Создание тестового запроса
	req, _ := http.NewRequest("GET", "/posts/1/comments", nil)
	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	// Вызов метода
	commentHandler.GetCommentsByPostID(c)

	// Проверка результата
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to get comments")

	// Проверка вызовов моков
	mockCommentUsecase.AssertExpectations(t)
}
