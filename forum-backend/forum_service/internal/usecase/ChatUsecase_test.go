package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"github.com/Engls/forum-project2/forum_service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestChatUsecase_HandleMessage_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockChatRepo := new(mocks.ChatRepository)

	// Инициализация usecase
	chatUsecase := NewChatUsecase(mockChatRepo, logger)

	// Тестовые данные
	userID := 1
	username := "testuser"
	content := "This is a test message"
	message := entity.ChatMessage{
		UserID:    userID,
		Username:  username,
		Content:   content,
		Timestamp: time.Now(),
	}

	// Настройка моков
	mockChatRepo.On("StoreMessage", mock.Anything, message).Return(nil)

	// Вызов метода
	err := chatUsecase.HandleMessage(context.Background(), userID, username, content)

	// Проверка результата
	assert.NoError(t, err)

	// Проверка вызовов моков
	mockChatRepo.AssertExpectations(t)
}

func TestChatUsecase_HandleMessage_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockChatRepo := new(mocks.ChatRepository)

	// Инициализация usecase
	chatUsecase := NewChatUsecase(mockChatRepo, logger)

	// Тестовые данные
	userID := 1
	username := "testuser"
	content := "This is a test message"
	message := entity.ChatMessage{
		UserID:    userID,
		Username:  username,
		Content:   content,
		Timestamp: time.Now(),
	}

	// Настройка моков
	mockChatRepo.On("StoreMessage", mock.Anything, message).Return(errors.New("failed to store message"))

	// Вызов метода
	err := chatUsecase.HandleMessage(context.Background(), userID, username, content)

	// Проверка результата
	assert.Error(t, err)

	// Проверка вызовов моков
	mockChatRepo.AssertExpectations(t)
}

func TestChatUsecase_GetRecentMessages_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockChatRepo := new(mocks.ChatRepository)

	// Инициализация usecase
	chatUsecase := NewChatUsecase(mockChatRepo, logger)

	// Тестовые данные
	limit := 10
	messages := []entity.ChatMessage{
		{UserID: 1, Username: "user1", Content: "Message 1", Timestamp: time.Now()},
		{UserID: 2, Username: "user2", Content: "Message 2", Timestamp: time.Now()},
	}

	// Настройка моков
	mockChatRepo.On("GetRecentMessages", mock.Anything, limit).Return(messages, nil)

	// Вызов метода
	result, err := chatUsecase.GetRecentMessages(context.Background(), limit)

	// Проверка результата
	assert.NoError(t, err)
	assert.Equal(t, messages, result)

	// Проверка вызовов моков
	mockChatRepo.AssertExpectations(t)
}

func TestChatUsecase_GetRecentMessages_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockChatRepo := new(mocks.ChatRepository)

	// Инициализация usecase
	chatUsecase := NewChatUsecase(mockChatRepo, logger)

	// Тестовые данные
	limit := 10

	// Настройка моков
	mockChatRepo.On("GetRecentMessages", mock.Anything, limit).Return(nil, errors.New("failed to get recent messages"))

	// Вызов метода
	result, err := chatUsecase.GetRecentMessages(context.Background(), limit)

	// Проверка результата
	assert.Error(t, err)
	assert.Nil(t, result)

	// Проверка вызовов моков
	mockChatRepo.AssertExpectations(t)
}
