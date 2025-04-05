package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest"
	"github.com/LearnShareApp/learn-share-backend/pkg/livekit"
	"github.com/LearnShareApp/learn-share-backend/pkg/migrator"
	"github.com/LearnShareApp/learn-share-backend/pkg/storage/db/postgres"
	"github.com/LearnShareApp/learn-share-backend/pkg/storage/object/minio"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	maxPort      = 1<<16 - 1
	maskedString = "********"
)

type Config struct {
	DB           postgres.Config
	Migrator     migrator.Config
	Server       rest.Config
	LiveKit      livekit.Config
	Minio        minio.Config
	IsInitDb     bool   `env:"IS_INIT_DB" env-required:"true"`
	JwtSecretKey string `env:"SECRET_KEY" env-required:"true"`
}

func LoadConfig(paths []string) (*Config, error) {
	// Looking for .env file in different directories
	var envPath string

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			envPath = path

			break
		}
	}

	if envPath == "" {
		return nil, fmt.Errorf(".env file not found in any of the search paths: %v", paths) //nolint:err113
	}

	var config Config

	err := cleanenv.ReadConfig(envPath, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// Validate config validation.
func (c *Config) Validate() error {
	if !checkPortValidation(c.Server.Port) {
		return fmt.Errorf("invalid server port: %d", c.Server.Port) //nolint:err113
	}

	if !checkPortValidation(c.DB.Port) {
		return fmt.Errorf("invalid database port: %d", c.DB.Port) //nolint:err113
	}

	if !checkPortValidation(c.Migrator.Port) {
		return fmt.Errorf("invalid database (migrator) port: %d", c.DB.Port) //nolint:err113
	}

	if !checkPortValidation(c.Minio.Port) {
		return fmt.Errorf("invalid minio port: %d", c.Minio.Port) //nolint:err113
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

	if logConfig.Migrator.Password != "" {
		logConfig.Migrator.Password = maskedString
	}

	if logConfig.JwtSecretKey != "" {
		logConfig.JwtSecretKey = maskedString
	}

	if logConfig.LiveKit.APIKey != "" {
		logConfig.LiveKit.APIKey = maskedString
	}

	if logConfig.LiveKit.APISecret != "" {
		logConfig.LiveKit.APISecret = maskedString
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

func checkPortValidation(port int) bool {
	return port >= 1 && port <= maxPort
}
