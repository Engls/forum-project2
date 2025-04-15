package main

import (
	"fmt"
	"log"

	"github.com/Engls/forum-project2/auth_service/internal/config"
	"github.com/Engls/forum-project2/auth_service/internal/delivery/http"
	"github.com/Engls/forum-project2/auth_service/internal/repository"
	"github.com/Engls/forum-project2/auth_service/internal/usecase"
	"github.com/Engls/forum-project2/auth_service/internal/utils" // <--- Добавьте этот импорт

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx" // Импортируем sqlx
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Выводим значения переменных окружения для отладки
	fmt.Printf("AUTH_SERVICE_PORT: %s\n", cfg.Port)
	fmt.Printf("AUTH_SERVICE_DB_PATH: %s\n", cfg.DBPath)
	fmt.Printf("AUTH_SERVICE_MIGRATIONS_PATH: %s\n", cfg.MigrationsPath)
	fmt.Printf("JWT_SECRET: %s\n", cfg.JWTSecret)

	// Открываем соединение с базой данных с помощью sqlx
	db, err := sqlx.Open("sqlite3", cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// --- Миграции ---
	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{}) // Используем db.DB для миграций
	if err != nil {
		log.Fatalf("Failed to create migrate driver: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+cfg.MigrationsPath,
		"sqlite3", driver)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	fmt.Println("Миграции успешно применены!")
	// --- Конец миграций ---

	userRepo := repository.NewAuthRepository(db)
	jwtUtil := utils.NewJWTUtil(cfg.JWTSecret)               // <--- Создаем экземпляр JWTUtil
	userUsecase := usecase.NewAuthUsecase(userRepo, jwtUtil) // <--- Передаем jwtUtil

	authHandler := http.NewAuthHandler(userUsecase)

	router := gin.Default()
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	if err := router.Run(cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
