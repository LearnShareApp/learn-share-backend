package migrator

import (
	"errors"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang-migrate/migrate/v4"
)

type Config struct {
	MigrationsPath string `env:"MIGRATIONS_PATH"   env-default:"migrations"`
	Host           string `env:"POSTGRES_HOST"     env-required:"true"`
	Port           int    `env:"POSTGRES_PORT"     env-required:"true"`
	DBName         string `env:"POSTGRES_DB"       env-required:"true"`
	UserName       string `env:"POSTGRES_USER"     env-required:"true"`
	Password       string `env:"POSTGRES_PASSWORD" env-required:"true"`
	SSLMode        string `env:"POSTGRES_SSLMODE"  env-default:"disable"`
}

func RunMigrations(cfg *Config) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	migrations, err := migrate.New(
		fmt.Sprintf("file://%s", cfg.MigrationsPath), //nolint:perfsprint    // path to migrations files
		dsn, // connection string to DB
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	defer migrations.Close()

	if err := migrations.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}

		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
