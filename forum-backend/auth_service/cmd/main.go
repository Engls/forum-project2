package main

import (
	"database/sql"
	"forum/auth_service/internal/config"
	"forum/auth_service/internal/delivery/http"
	"forum/auth_service/internal/repository"
	"forum/auth_service/internal/usecase"
	"forum/auth_service/internal/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg := config.LoadConfig()

	db, err := sqlx.Connect("sqlite3", cfg.DBPath)
	if err != nil {
		logger.Fatal("Failed to connect to DataBase", zap.Error(err))
	}
	defer db.Close()

	// Выполнение миграций программно
	err = performMigrations(db.DB)
	if err != nil {
		logger.Fatal("Failed to perform migrations", zap.Error(err))
	}

	authRepo := repository.NewAuthRepository(db)
	jwtUtil := utils.NewJWTUtil(cfg.JWTSecret)
	authUsecase := usecase.NewAuthUsecase(authRepo, jwtUtil)
	authHandler := http.NewAuthHandler(authUsecase)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	if err := r.Run(":8080"); err != nil {
		logger.Fatal("Failed to run server", zap.Error(err))
	}
}

// performMigrations выполняет миграции программно
func performMigrations(db *sql.DB) error {
	// SQL для создания таблицы users
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user'
	);`

	// SQL для создания таблицы tokens
	createTokenTable := `
	CREATE TABLE IF NOT EXISTS tokens (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		token TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	// Выполнение SQL запросов
	_, err := db.Exec(createUserTable)
	if err != nil {
		return err
	}

	_, err = db.Exec(createTokenTable)
	if err != nil {
		return err
	}

	return nil
}
