package main

import (
	utils "github.com/Engls/EnglsJwt"
	_ "github.com/Engls/forum-project2/forum_service/docs"
	"github.com/Engls/forum-project2/forum_service/internal/config"
	"github.com/Engls/forum-project2/forum_service/internal/controllers/chat"
	"github.com/Engls/forum-project2/forum_service/internal/controllers/http"
	"github.com/Engls/forum-project2/forum_service/internal/repository"
	"github.com/Engls/forum-project2/forum_service/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"time"
)

// @title Forum Service API
// @version 1.0
// @description This is the API documentation for the Auth Service.
// @host localhost:8081
// @BasePath /
func main() {
	// Инициализация логгера
	utils.InitLogger()
	logger := utils.GetLogger()

	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	logger.Info("Configuration loaded",
		zap.String("DB_PATH", cfg.DBPath),
		zap.String("PORT", cfg.Port),
		zap.String("JWT_SECRET", cfg.JWTSecret),
	)

	// Подключение к базе данных
	db, err := sqlx.Open("sqlite3", cfg.DBPath)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Инициализация репозиториев и usecase
	postRepo := repository.NewPostRepository(db, logger)
	commentRepo := repository.NewCommentsRepository(db, logger)
	chatRepo := repository.NewChatRepository(db, logger)
	postUsecase := usecase.NewPostUsecase(postRepo, logger)
	commentUsecase := usecase.NewCommentsUsecases(commentRepo, logger)
	hub := chat.NewHub()
	chatUsecase := usecase.NewChatUsecase(chatRepo, logger)
	jwtUtil := utils.NewJWTUtil(cfg.JWTSecret)

	// Инициализация обработчиков
	postHandler := http.NewPostHandler(postUsecase, postRepo, jwtUtil, logger)
	commentHandler := http.NewCommentHandler(commentUsecase, jwtUtil, logger)
	chatHandler := http.NewChatHandler(hub, chatUsecase, jwtUtil, logger)

	// Запуск хаба
	go hub.Run()

	// Настройка маршрутизатора
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Content-type", "Origin", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/ws", chatHandler.ServeWS)
	router.POST("/posts", postHandler.CreatePost)
	router.GET("/posts", postHandler.GetPosts)
	router.DELETE("/posts/:id", postHandler.DeletePost)
	router.POST("/posts/:id/comments", commentHandler.CreateComment)
	router.GET("/posts/:id/comments", commentHandler.GetCommentsByPostID)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Запуск сервера
	if err := router.Run(cfg.Port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
