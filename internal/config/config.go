package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/pkg/storage/db/postgres"
	"github.com/LearnShareApp/learn-share-backend/pkg/storage/object/minio"
	"os"
	"strconv"

	"github.com/LearnShareApp/learn-share-backend/internal/service/livekit"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest"
	"github.com/joho/godotenv"
)

const (
	maxPort      = 1<<16 - 1
	maskedString = "********"
)

type Config struct {
	DB           postgres.DBConfig     `json:"db"`
	Server       rest.ServerConfig     `json:"server"`
	LiveKit      livekit.LiveKitConfig `json:"livekit"`
	Minio        minio.MinioConfig     `json:"minio"`
	IsInitDb     bool                  `json:"is_init_db"`
	JwtSecretKey string                `json:"jwt_secret_key"`
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

	var config Config

	// Server config
	var err error

	config.Server.Port, err = getEnvAsInt("SERVER_PORT")
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}

	// Database config
	config.DB.Port, err = getEnvAsInt("DB_PORT")
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	config.DB.Host = os.Getenv("DB_HOST")
	config.DB.DbName = os.Getenv("DB_NAME")
	config.DB.User = os.Getenv("DB_USER")
	config.DB.Password = os.Getenv("DB_PASSWORD")

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
	value := os.Getenv("IS_MINIO_SSL")

	config.Minio.IsSSL, err = strconv.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf("invalid IS_MINIO_SSL: %w", err)
	}

	// Should Init DB
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

	return &config, nil
}

func getEnvAsInt(key string) (int, error) {
	valueStr := os.Getenv(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s must be an integer, got %s", key, valueStr)
	}

	return value, nil
}

// Validate config validation.
func (c *Config) Validate() error {
	if !isValidPort(c.Server.Port) {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if !isValidPort(c.DB.Port) {
		return fmt.Errorf("invalid database port: %d", c.DB.Port)
	}

	if c.DB.Host == "" {
		return errors.New("database host cannot be empty")
	}

	if c.DB.DbName == "" {
		return errors.New("database name cannot be empty")
	}

	if c.DB.User == "" {
		return errors.New("database user cannot be empty")
	}

	if c.JwtSecretKey == "" {
		return errors.New("jwt secret key cannot be empty")
	}

	if c.LiveKit.ApiKey == "" {
		return errors.New("live kit api key cannot be empty")
	}

	if c.LiveKit.ApiSecret == "" {
		return errors.New("live kit api secret cannot be empty")
	}

	if !isValidPort(c.Minio.Port) {
		return fmt.Errorf("invalid minio port: %d", c.Minio.Port)
	}

	if c.Minio.Host == "" {
		return errors.New("minio host cannot be empty")
	}

	if c.Minio.AccessKey == "" {
		return errors.New("minio access key cannot be empty")
	}

	if c.Minio.SecretKey == "" {
		return errors.New("minio secret key cannot be empty")
	}

	if c.Minio.Bucket == "" {
		return errors.New("minio bucket cannot be empty")
	}

	return nil
}

// LogConfig logs configuration with sensitive data masking.
func (c *Config) LogConfig() (string, error) {
	// Create a copy of config for logging
	logConfig := *c

	// Mask passwords

	if logConfig.DB.Password != "" {
		logConfig.DB.Password = maskedString
	}

	if logConfig.JwtSecretKey != "" {
		logConfig.JwtSecretKey = maskedString
	}

	if logConfig.LiveKit.ApiKey != "" {
		logConfig.LiveKit.ApiKey = maskedString
	}

	if logConfig.LiveKit.ApiSecret != "" {
		logConfig.LiveKit.ApiSecret = maskedString
	}

	if logConfig.Minio.AccessKey != "" {
		logConfig.Minio.AccessKey = maskedString
	}

	if logConfig.Minio.SecretKey != "" {
		logConfig.Minio.SecretKey = maskedString
	}

	// Convert to JSON with indents for readability
	jsonBytes, err := json.MarshalIndent(logConfig, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling config: %w", err)
	}

	return "Application Configuration:\n" + string(jsonBytes), nil
}

func isValidPort(port int) bool {
	return port >= 1 && port <= maxPort
}
