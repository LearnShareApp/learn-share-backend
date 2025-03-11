package postgres

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type DBConfig struct {
	Host     string `env:"DB_HOST"     env-required:"true"`
	Port     int    `env:"DB_PORT"     env-required:"true"`
	DBName   string `env:"DB_NAME"     env-required:"true"`
	User     string `env:"DB_USER"     env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
}

func New(ctx context.Context, config *DBConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%d",
		config.User,
		config.Password,
		config.DBName,
		config.Host,
		config.Port,
	)

	databaseConn, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := databaseConn.Ping(); err != nil {
		return nil, err
	}

	return databaseConn, nil
}
