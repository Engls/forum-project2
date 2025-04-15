package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	DBPath    string
	JWTSecret string
}

func LoadConfig() Config {
	// Получаем текущую директорию
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Формируем абсолютный путь к общей базе данных
	dbPath := filepath.Join(dir, "..", "..", "forum.db")

	return Config{
		DBPath:    dbPath, // Абсолютный путь к общей базе данных
		JWTSecret: "your_jwt_secret",
	}
}
