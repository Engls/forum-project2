package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBPath         string
	MigrationsPath string // Добавляем путь к миграциям
}

func LoadConfig() (*Config, error) {
	godotenv.Load()

	cfg := &Config{
		Port:           getEnv("FORUM_SERVICE_PORT", ":8081"),
		DBPath:         getEnv("FORUM_SERVICE_DB_PATH", "forum.db"),
		MigrationsPath: getEnv("FORUM_SERVICE_MIGRATIONS_PATH", "migrations"), // Получаем путь из переменной окружения
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
