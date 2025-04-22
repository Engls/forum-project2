package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"github.com/Engls/forum-project2/forum_service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestPostUsecase_CreatePost_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockPostRepo := new(mocks.PostRepository)

	// Инициализация usecase
	postUsecase := NewPostUsecase(mockPostRepo, logger)

	// Тестовые данные
	post := entity.Post{
		AuthorId: 1,
		Title:    "Test Post",
		Content:  "This is a test post",
	}
	createdPost := post
	createdPost.ID = 1

	// Настройка моков
	mockPostRepo.On("CreatePost", mock.Anything, post).Return(&createdPost, nil)

	// Вызов метода
	result, err := postUsecase.CreatePost(context.Background(), post)

	// Проверка результата
	assert.NoError(t, err)
	assert.Equal(t, &createdPost, result)

	// Проверка вызовов моков
	mockPostRepo.AssertExpectations(t)
}

func TestPostUsecase_CreatePost_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockPostRepo := new(mocks.PostRepository)

	// Инициализация usecase
	postUsecase := NewPostUsecase(mockPostRepo, logger)

	// Тестовые данные
	post := entity.Post{
		AuthorId: 1,
		Title:    "Test Post",
		Content:  "This is a test post",
	}

	// Настройка моков
	mockPostRepo.On("CreatePost", mock.Anything, post).Return(nil, errors.New("failed to create post"))

	// Вызов метода
	result, err := postUsecase.CreatePost(context.Background(), post)

	// Проверка результата
	assert.Error(t, err)
	assert.Nil(t, result)

	// Проверка вызовов моков
	mockPostRepo.AssertExpectations(t)
}

func TestPostUsecase_GetPosts_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockPostRepo := new(mocks.PostRepository)

	// Инициализация usecase
	postUsecase := NewPostUsecase(mockPostRepo, logger)

	// Тестовые данные
	posts := []entity.Post{
		{ID: 1, AuthorId: 1, Title: "Post 1", Content: "Content 1"},
		{ID: 2, AuthorId: 2, Title: "Post 2", Content: "Content 2"},
	}

	// Настройка моков
	mockPostRepo.On("GetPosts", mock.Anything).Return(posts, nil)

	// Вызов метода
	result, err := postUsecase.GetPosts(context.Background())

	// Проверка результата
	assert.NoError(t, err)
	assert.Equal(t, posts, result)

	// Проверка вызовов моков
	mockPostRepo.AssertExpectations(t)
}

func TestPostUsecase_GetPosts_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockPostRepo := new(mocks.PostRepository)

	// Инициализация usecase
	postUsecase := NewPostUsecase(mockPostRepo, logger)

	// Настройка моков
	mockPostRepo.On("GetPosts", mock.Anything).Return(nil, errors.New("failed to get posts"))

	// Вызов метода
	result, err := postUsecase.GetPosts(context.Background())

	// Проверка результата
	assert.Error(t, err)
	assert.Nil(t, result)

	// Проверка вызовов моков
	mockPostRepo.AssertExpectations(t)
}

func TestPostUsecase_GetPostByID_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockPostRepo := new(mocks.PostRepository)

	// Инициализация usecase
	postUsecase := NewPostUsecase(mockPostRepo, logger)

	// Тестовые данные
	post := entity.Post{
		ID:       1,
		AuthorId: 1,
		Title:    "Test Post",
		Content:  "This is a test post",
	}

	// Настройка моков
	mockPostRepo.On("GetPostByID", mock.Anything, 1).Return(&post, nil)

	// Вызов метода
	result, err := postUsecase.GetPostByID(context.Background(), 1)

	// Проверка результата
	assert.NoError(t, err)
	assert.Equal(t, &post, result)

	// Проверка вызовов моков
	mockPostRepo.AssertExpectations(t)
}

func TestPostUsecase_GetPostByID_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockPostRepo := new(mocks.PostRepository)

	// Инициализация usecase
	postUsecase := NewPostUsecase(mockPostRepo, logger)

	// Настройка моков
	mockPostRepo.On("GetPostByID", mock.Anything, 1).Return(nil, errors.New("failed to get post"))

	// Вызов метода
	result, err := postUsecase.GetPostByID(context.Background(), 1)

	// Проверка результата
	assert.Error(t, err)
	assert.Nil(t, result)

	// Проверка вызовов моков
	mockPostRepo.AssertExpectations(t)
}

func TestPostUsecase_UpdatePost_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockPostRepo := new(mocks.PostRepository)

	// Инициализация usecase
	postUsecase := NewPostUsecase(mockPostRepo, logger)

	// Тестовые данные
	post := entity.Post{
		ID:       1,
		AuthorId: 1,
		Title:    "Updated Post",
		Content:  "This is an updated post",
	}

	// Настройка моков
	mockPostRepo.On("UpdatePost", mock.Anything, post).Return(&post, nil)

	// Вызов метода
	result, err := postUsecase.UpdatePost(context.Background(), post)

	// Проверка результата
	assert.NoError(t, err)
	assert.Equal(t, &post, result)

	// Проверка вызовов моков
	mockPostRepo.AssertExpectations(t)
}

func TestPostUsecase_UpdatePost_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockPostRepo := new(mocks.PostRepository)

	// Инициализация usecase
	postUsecase := NewPostUsecase(mockPostRepo, logger)

	// Тестовые данные
	post := entity.Post{
		ID:       1,
		AuthorId: 1,
		Title:    "Updated Post",
		Content:  "This is an updated post",
	}

	// Настройка моков
	mockPostRepo.On("UpdatePost", mock.Anything, post).Return(nil, errors.New("failed to update post"))

	// Вызов метода
	result, err := postUsecase.UpdatePost(context.Background(), post)

	// Проверка результата
	assert.Error(t, err)
	assert.Nil(t, result)

	// Проверка вызовов моков
	mockPostRepo.AssertExpectations(t)
}

func TestPostUsecase_DeletePost_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockPostRepo := new(mocks.PostRepository)

	// Инициализация usecase
	postUsecase := NewPostUsecase(mockPostRepo, logger)

	// Настройка моков
	mockPostRepo.On("DeletePost", mock.Anything, 1).Return(nil)

	// Вызов метода
	err := postUsecase.DeletePost(context.Background(), 1)

	// Проверка результата
	assert.NoError(t, err)

	// Проверка вызовов моков
	mockPostRepo.AssertExpectations(t)
}

func TestPostUsecase_DeletePost_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockPostRepo := new(mocks.PostRepository)

	// Инициализация usecase
	postUsecase := NewPostUsecase(mockPostRepo, logger)

	// Настройка моков
	mockPostRepo.On("DeletePost", mock.Anything, 1).Return(errors.New("failed to delete post"))

	// Вызов метода
	err := postUsecase.DeletePost(context.Background(), 1)

	// Проверка результата
	assert.Error(t, err)

	// Проверка вызовов моков
	mockPostRepo.AssertExpectations(t)
}
