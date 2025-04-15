package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBPath         string
	JWTSecret      string
	MigrationsPath string // Добавляем путь к миграциям
}

func LoadConfig() (*Config, error) {
	godotenv.Load()

	cfg := &Config{
		Port:           getEnv("AUTH_SERVICE_PORT", ":8080"),
		DBPath:         getEnv("AUTH_SERVICE_DB_PATH", "auth.db"),
		JWTSecret:      getEnv("JWT_SECRET", "secret"),
		MigrationsPath: getEnv("AUTH_SERVICE_MIGRATIONS_PATH", "migrations"), // Получаем путь из переменной окружения
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
