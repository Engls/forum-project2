package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"log"
	"os"
	"time"

	"github.com/Engls/forum-project2/auth_service/internal/config"
	"github.com/Engls/forum-project2/auth_service/internal/delivery/http"
	"github.com/Engls/forum-project2/auth_service/internal/repository"
	"github.com/Engls/forum-project2/auth_service/internal/usecase"
	"github.com/Engls/forum-project2/auth_service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("AUTH_SERVICE_PORT: %s\n", cfg.Port)
	fmt.Printf("AUTH_SERVICE_DB_PATH: %s\n", cfg.DBPath)
	fmt.Printf("AUTH_SERVICE_MIGRATIONS_PATH: %s\n", cfg.MigrationsPath)
	fmt.Printf("JWT_SECRET: %s\n", cfg.JWTSecret)

	db, err := sqlx.Open("sqlite3", cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("Failed to create migrate driver: %v", err)
	}

	fmt.Println(os.Getwd())
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

	userRepo := repository.NewAuthRepository(db)
	jwtUtil := utils.NewJWTUtil(cfg.JWTSecret)
	userUsecase := usecase.NewAuthUsecase(userRepo, jwtUtil)

	authHandler := http.NewAuthHandler(userUsecase)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders:     []string{"Content-type", "Origin", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	if err := router.Run(cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
