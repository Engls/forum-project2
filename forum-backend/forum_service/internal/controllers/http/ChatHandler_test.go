package http

import (
	utils "github.com/Engls/EnglsJwt"
	"net/http/httptest"
	"testing"

	"github.com/Engls/forum-project2/forum_service/internal/controllers/chat"
	"github.com/Engls/forum-project2/forum_service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestChatHandler_ServeWS_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockChatUsecase := new(mocks.ChatUsecase)
	jwtUtil := utils.NewJWTUtil("secret")
	hub := chat.NewHub()

	// Инициализация обработчика
	chatHandler := NewChatHandler(hub, mockChatUsecase, jwtUtil, logger)

	// Генерация токена для пользователя
	token, err := jwtUtil.GenerateToken(1, "user")
	assert.NoError(t, err)

	// Создание тестового сервера
	router := gin.Default()
	router.GET("/ws/chat", chatHandler.ServeWS)

	server := httptest.NewServer(router)
	defer server.Close()

	// Создание WebSocket клиента
	url := "ws" + server.URL[4:] + "/ws/chat?token=" + token + "&userID=1&username=user&auth=true"
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)
	defer ws.Close()

	// Проверка соединения
	assert.NotNil(t, ws)
}

func TestChatHandler_ServeWS_InvalidToken(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockChatUsecase := new(mocks.ChatUsecase)
	jwtUtil := utils.NewJWTUtil("secret")
	hub := chat.NewHub()

	// Инициализация обработчика
	chatHandler := NewChatHandler(hub, mockChatUsecase, jwtUtil, logger)

	// Создание тестового сервера
	router := gin.Default()
	router.GET("/ws/chat", chatHandler.ServeWS)

	server := httptest.NewServer(router)
	defer server.Close()

	// Создание WebSocket клиента
	url := "ws" + server.URL[4:] + "/ws/chat?token=invalid_token&userID=1&username=user&auth=true"
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)
	defer ws.Close()

	// Проверка соединения
	assert.NotNil(t, ws)
}

func TestChatHandler_ServeWS_InvalidUserID(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	mockChatUsecase := new(mocks.ChatUsecase)
	jwtUtil := utils.NewJWTUtil("secret")
	hub := chat.NewHub()

	// Инициализация обработчика
	chatHandler := NewChatHandler(hub, mockChatUsecase, jwtUtil, logger)

	// Генерация токена для пользователя
	token, err := jwtUtil.GenerateToken(1, "user")
	assert.NoError(t, err)

	// Создание тестового сервера
	router := gin.Default()
	router.GET("/ws/chat", chatHandler.ServeWS)

	server := httptest.NewServer(router)
	defer server.Close()

	// Создание WebSocket клиента
	url := "ws" + server.URL[4:] + "/ws/chat?token=" + token + "&userID=invalid&username=user&auth=true"
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)
	defer ws.Close()

	// Проверка соединения
	assert.NotNil(t, ws)
}
