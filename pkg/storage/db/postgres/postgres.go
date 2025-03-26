package postgres

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host     string `env:"POSTGRES_HOST"     env-required:"true"`
	Port     int    `env:"POSTGRES_PORT"     env-required:"true"`
	DBName   string `env:"POSTGRES_DB"       env-required:"true"`
	User     string `env:"POSTGRES_USER"     env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
}

func New(ctx context.Context, config *Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%d",
		config.User,
		config.Password,
		config.DBName,
		config.Host,
		config.Port,
	)

	databaseConn, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	if err := databaseConn.Ping(); err != nil {
		return nil, err //nolint:wrapcheck
	}

	return databaseConn, nil
}
