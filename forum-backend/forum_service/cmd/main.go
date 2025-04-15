package main

import (
	"database/sql"
	"forum-project/forum_service/internal/config"
	"forum-project/forum_service/internal/controllers/http"
	"forum-project/forum_service/internal/repository"
	"forum-project/forum_service/internal/usecase"
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

	postRepo := repository.NewPostRepository(db)
	postUsecase := usecase.NewPostUsecases(postRepo)
	postHandler := http.NewPostHandler(postUsecase, postRepo)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.POST("/posts", postHandler.CreatePost)
	r.GET("/posts", postHandler.GetPosts)

	if err := r.Run(":8081"); err != nil {
		logger.Fatal("Failed to run server", zap.Error(err))
	}
}

// performMigrations выполняет миграции программно
func performMigrations(db *sql.DB) error {
	// SQL для создания таблицы posts
	createPostTable := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		author_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		FOREIGN KEY (author_id) REFERENCES users(id)
	);`

	// Выполнение SQL запросов
	_, err := db.Exec(createPostTable)
	if err != nil {
		return err
	}

	return nil
}
