// auth_service/internal/config/config.go
package config

import (
	"os"

	"github.com/joho/godotenv" // Если используете .env
)

type Config struct {
	Port           string
	DBPath         string // <--- Измените здесь
	MigrationsPath string
	JWTSecret      string
}

func LoadConfig() (Config, error) {
	// ...
	// Загрузка переменных окружения (если используете .env)
	err := godotenv.Load()
	if err != nil {
		// Если .env не найден, продолжаем без него
		// return Config{}, fmt.Errorf("error loading .env file: %w", err) // Или обработайте иначе
	}

	cfg := Config{
		Port:           getEnv("AUTH_SERVICE_PORT", ":8080"),
		DBPath:         getEnv("AUTH_SERVICE_DB_PATH", "./db/forum.db"), // <--- Измените здесь
		MigrationsPath: getEnv("AUTH_SERVICE_MIGRATIONS_PATH", "./migrations"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key"),
	}
	// ...
	return cfg, nil
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
