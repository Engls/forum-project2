package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"github.com/Engls/forum-project2/forum_service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestChatRepo_StoreMessage_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockDB := new(mocks.DB)

	// Инициализация репозитория
	chatRepo := NewChatRepository(mockDB, logger)

	// Тестовые данные
	msg := entity.ChatMessage{
		UserID:    1,
		Username:  "testuser",
		Content:   "This is a test message",
		Timestamp: time.Date(2025, time.April, 22, 23, 51, 38, 843016900, time.Local),
	}

	// Настройка моков
	mockDB.On("ExecContext", mock.Anything, mock.Anything, msg.UserID, msg.Username, msg.Content, msg.Timestamp.Format(time.RFC3339)).Return(sql.Result(nil), nil)

	// Вызов метода
	err := chatRepo.StoreMessage(context.Background(), msg)

	// Проверка результата
	assert.NoError(t, err)

	// Проверка вызовов моков
	mockDB.AssertExpectations(t)
}

func TestChatRepo_StoreMessage_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockDB := new(mocks.DB)

	// Инициализация репозитория
	chatRepo := NewChatRepository(mockDB, logger)

	// Тестовые данные
	msg := entity.ChatMessage{
		UserID:    1,
		Username:  "testuser",
		Content:   "This is a test message",
		Timestamp: time.Now(),
	}

	// Настройка моков
	mockDB.On("ExecContext", mock.Anything, mock.Anything, msg.UserID, msg.Username, msg.Content, msg.Timestamp.Format(time.RFC3339)).Return(nil, errors.New("failed to store message"))

	// Вызов метода
	err := chatRepo.StoreMessage(context.Background(), msg)

	// Проверка результата
	assert.Error(t, err)

	// Проверка вызовов моков
	mockDB.AssertExpectations(t)
}

func TestChatRepo_GetRecentMessages_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockDB := new(mocks.DB)

	// Инициализация репозитория
	chatRepo := NewChatRepository(mockDB, logger)

	// Тестовые данные
	limit := 10
	messages := []entity.ChatMessage{
		{ID: 1, UserID: 1, Username: "user1", Content: "Message 1", Timestamp: time.Now()},
		{ID: 2, UserID: 2, Username: "user2", Content: "Message 2", Timestamp: time.Now()},
	}

	// Настройка моков
	mockDB.On("SelectContext", mock.Anything, mock.Anything, mock.Anything, limit).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(1).(*[]entity.ChatMessage)
		*dest = messages
	})

	// Вызов метода
	result, err := chatRepo.GetRecentMessages(context.Background(), limit)

	// Проверка результата
	assert.NoError(t, err)
	assert.Equal(t, messages, result)

	// Проверка вызовов моков
	mockDB.AssertExpectations(t)
}

func TestChatRepo_GetRecentMessages_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockDB := new(mocks.DB)

	// Инициализация репозитория
	chatRepo := NewChatRepository(mockDB, logger)

	// Тестовые данные
	limit := 10

	// Настройка моков
	mockDB.On("SelectContext", mock.Anything, mock.Anything, mock.Anything, limit).Return(errors.New("failed to get recent messages"))

	// Вызов метода
	result, err := chatRepo.GetRecentMessages(context.Background(), limit)

	// Проверка результата
	assert.Error(t, err)
	assert.Nil(t, result)

	// Проверка вызовов моков
	mockDB.AssertExpectations(t)
}
