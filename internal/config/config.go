package config

import (
	"encoding/json"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/service/livekit"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest"
	"github.com/LearnShareApp/learn-share-backend/pkg/db/postgres"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	Db           postgres.DbConfig
	Server       rest.ServerConfig
	LiveKit      livekit.ApiConfig
	IsInitDb     bool
	JwtSecretKey string
}

func LoadConfig() (*Config, error) {
	// Ищем .env файл в разных директориях
	envPaths := []string{
		".env",
		"./config/.env",
	}

	envFound := false
	for _, path := range envPaths {
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Load(path); err != nil {
				return nil, fmt.Errorf("error loading %s file: %w", path, err)
			}
			envFound = true
			break
		}
	}

	if !envFound {
		return nil, fmt.Errorf(".env file not found in any of the search paths: %v", envPaths)
	}

	config := &Config{}

	// Server config
	var err error
	config.Server.Port, err = getEnvAsInt("SERVER_PORT")
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}

	// Database config
	config.Db.Port, err = getEnvAsInt("DB_PORT")
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	config.Db.Host = os.Getenv("DB_HOST")
	config.Db.DbName = os.Getenv("DB_NAME")
	config.Db.User = os.Getenv("DB_USER")
	config.Db.Password = os.Getenv("DB_PASSWORD")

	// Should Init Db
	value := os.Getenv("IS_INIT_DB")
	config.IsInitDb, err = strconv.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf("invalid IS_INIT_DB: %w", err)
	}

	// JWT secret key
	config.JwtSecretKey = os.Getenv("SECRET_KEY")

	// LiveKit config
	config.LiveKit.ApiKey = os.Getenv("LIVEKIT_API_KEY")
	config.LiveKit.ApiSecret = os.Getenv("LIVEKIT_API_SECRET")

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// getEnvAsInt теперь возвращает ошибку вместо значения по умолчанию
func getEnvAsInt(key string) (int, error) {
	valueStr := os.Getenv(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s must be an integer, got %s", key, valueStr)
	}
	return value, nil
}

// Validate валидация конфига
func (c *Config) Validate() error {

	if c.Db.Port < 1 || c.Db.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", c.Db.Port)
	}

	if c.Db.Host == "" {
		return fmt.Errorf("database host cannot be empty")
	}

	if c.Db.DbName == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	if c.Db.User == "" {
		return fmt.Errorf("database user cannot be empty")
	}

	if c.JwtSecretKey == "" {
		return fmt.Errorf("jwt secret key cannot be empty")
	}

	if c.LiveKit.ApiKey == "" {
		return fmt.Errorf("live kit api key cannot be empty")
	}

	if c.LiveKit.ApiSecret == "" {
		return fmt.Errorf("live kit api secret cannot be empty")
	}

	return nil
}

// LogConfig логирует конфигурацию с маскированием чувствительных данных
func (c *Config) LogConfig() string {
	// Создаем копию конфига для логирования
	logConfig := *c

	// Маскируем пароль
	if logConfig.Db.Password != "" {
		logConfig.Db.Password = "********"
	}

	if logConfig.JwtSecretKey != "" {
		logConfig.JwtSecretKey = "********"
	}

	if logConfig.LiveKit.ApiSecret != "" {
		logConfig.LiveKit.ApiSecret = "********"
	}

	// Преобразуем в JSON с отступами для читаемости
	jsonBytes, err := json.MarshalIndent(logConfig, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error marshaling config: %v", err)
	}

	return fmt.Sprintf("Application Configuration:\n%s", string(jsonBytes))
}
