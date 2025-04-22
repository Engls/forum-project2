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

func TestCommentsUsecases_CreateComment_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockCommentRepo := new(mocks.CommentsRepository)

	// Инициализация usecase
	commentsUsecases := NewCommentsUsecases(mockCommentRepo, logger)

	// Тестовые данные
	comment := entity.Comment{
		PostId:   1,
		AuthorId: 1,
		Content:  "This is a test comment",
	}
	createdComment := comment
	createdComment.ID = 1

	// Настройка моков
	mockCommentRepo.On("CreateComment", mock.Anything, comment).Return(createdComment, nil)

	// Вызов метода
	result, err := commentsUsecases.CreateComment(context.Background(), comment)

	// Проверка результата
	assert.NoError(t, err)
	assert.Equal(t, createdComment, result)

	// Проверка вызовов моков
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentsUsecases_CreateComment_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockCommentRepo := new(mocks.CommentsRepository)

	// Инициализация usecase
	commentsUsecases := NewCommentsUsecases(mockCommentRepo, logger)

	// Тестовые данные
	comment := entity.Comment{
		PostId:   1,
		AuthorId: 1,
		Content:  "This is a test comment",
	}

	// Настройка моков
	mockCommentRepo.On("CreateComment", mock.Anything, comment).Return(entity.Comment{}, errors.New("failed to create comment"))

	// Вызов метода
	result, err := commentsUsecases.CreateComment(context.Background(), comment)

	// Проверка результата
	assert.Error(t, err)
	assert.Equal(t, entity.Comment{}, result)

	// Проверка вызовов моков
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentsUsecases_GetCommentByPostID_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockCommentRepo := new(mocks.CommentsRepository)

	// Инициализация usecase
	commentsUsecases := NewCommentsUsecases(mockCommentRepo, logger)

	// Тестовые данные
	comments := []entity.Comment{
		{ID: 1, PostId: 1, AuthorId: 1, Content: "Comment 1"},
		{ID: 2, PostId: 1, AuthorId: 2, Content: "Comment 2"},
	}

	// Настройка моков
	mockCommentRepo.On("GetCommentsByPostID", mock.Anything, 1).Return(comments, nil)

	// Вызов метода
	result, err := commentsUsecases.GetCommentByPostID(context.Background(), 1)

	// Проверка результата
	assert.NoError(t, err)
	assert.Equal(t, comments, result)

	// Проверка вызовов моков
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentsUsecases_GetCommentByPostID_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockCommentRepo := new(mocks.CommentsRepository)

	// Инициализация usecase
	commentsUsecases := NewCommentsUsecases(mockCommentRepo, logger)

	// Настройка моков
	mockCommentRepo.On("GetCommentsByPostID", mock.Anything, 1).Return(nil, errors.New("failed to get comments"))

	// Вызов метода
	result, err := commentsUsecases.GetCommentByPostID(context.Background(), 1)

	// Проверка результата
	assert.Error(t, err)
	assert.Nil(t, result)

	// Проверка вызовов моков
	mockCommentRepo.AssertExpectations(t)
}
