package main

import (
	"fmt"
	"github.com/Engls/forum-project2/forum_service/internal/config"
	"github.com/Engls/forum-project2/forum_service/internal/controllers/http"
	"github.com/Engls/forum-project2/forum_service/internal/repository"
	"github.com/Engls/forum-project2/forum_service/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Println(os.Getwd())
	db, err := sqlx.Open("sqlite3", cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	postRepo := repository.NewPostRepository(db)
	postUsecase := usecase.NewPostUsecase(postRepo)
	postHandler := http.NewPostHandler(postUsecase, postRepo)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders:     []string{"Content-type", "Origin", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.POST("/posts", postHandler.CreatePost)
	router.GET("/posts", postHandler.GetPosts)

	if err := router.Run(cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
