package config

import (
	"encoding/json"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/pkg/object_storage/minio"
	"os"
	"strconv"

	"github.com/LearnShareApp/learn-share-backend/internal/service/livekit"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest"
	"github.com/LearnShareApp/learn-share-backend/pkg/db/postgres"
	"github.com/joho/godotenv"
)

type Config struct {
	Db           postgres.DbConfig
	Server       rest.ServerConfig
	LiveKit      livekit.ApiConfig
	Minio        minio.MinioConfig
	IsInitDb     bool
	JwtSecretKey string
}

func LoadConfig() (*Config, error) {
	// Looking for .env file in different directories
	envPaths := []string{
		".env",
		"./config/.env",
		"./internal/config/.env",
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

	// LiveKit config
	config.LiveKit.ApiKey = os.Getenv("LIVEKIT_API_KEY")
	config.LiveKit.ApiSecret = os.Getenv("LIVEKIT_API_SECRET")

	// Minio config
	config.Minio.Port, err = getEnvAsInt("MINIO_PORT")
	if err != nil {
		return nil, fmt.Errorf("invalid MINIO_PORT: %w", err)
	}
	config.Minio.Host = os.Getenv("MINIO_HOST")
	config.Minio.AccessKey = os.Getenv("MINIO_ACCESS_KEY")
	config.Minio.SecretKey = os.Getenv("MINIO_SECRET_KEY")
	config.Minio.Bucket = os.Getenv("MINIO_BUCKET")
	var value = os.Getenv("IS_MINIO_SSL")
	config.Minio.IsSSL, err = strconv.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf("invalid IS_MINIO_SSL: %w", err)
	}

	// Should Init Db
	value = os.Getenv("IS_INIT_DB")
	config.IsInitDb, err = strconv.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf("invalid IS_INIT_DB: %w", err)
	}

	// JWT secret key
	config.JwtSecretKey = os.Getenv("SECRET_KEY")

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

func getEnvAsInt(key string) (int, error) {
	valueStr := os.Getenv(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s must be an integer, got %s", key, valueStr)
	}
	return value, nil
}

// Validate config validation
func (c *Config) Validate() error {

	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

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

	if c.Minio.Port < 1 || c.Minio.Port > 65535 {
		return fmt.Errorf("invalid minio port: %d", c.Minio.Port)
	}
	if c.Minio.Host == "" {
		return fmt.Errorf("minio host cannot be empty")
	}
	if c.Minio.AccessKey == "" {
		return fmt.Errorf("minio access key cannot be empty")
	}
	if c.Minio.SecretKey == "" {
		return fmt.Errorf("minio secret key cannot be empty")
	}
	if c.Minio.Bucket == "" {
		return fmt.Errorf("minio bucket cannot be empty")
	}

	return nil
}

// LogConfig logs configuration with sensitive data masking
func (c *Config) LogConfig() string {
	// Create a copy of config for logging
	logConfig := *c

	// Mask password
	if logConfig.Db.Password != "" {
		logConfig.Db.Password = "********"
	}

	if logConfig.JwtSecretKey != "" {
		logConfig.JwtSecretKey = "********"
	}

	if logConfig.LiveKit.ApiKey != "" {
		logConfig.LiveKit.ApiKey = "********"
	}
	if logConfig.LiveKit.ApiSecret != "" {
		logConfig.LiveKit.ApiSecret = "********"
	}

	if logConfig.Minio.AccessKey != "" {
		logConfig.Minio.AccessKey = "********"
	}
	if logConfig.Minio.SecretKey != "" {
		logConfig.Minio.SecretKey = "********"
	}

	// Convert to JSON with indents for readability
	jsonBytes, err := json.MarshalIndent(logConfig, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error marshaling config: %v", err)
	}

	return fmt.Sprintf("Application Configuration:\n%s", string(jsonBytes))
}
