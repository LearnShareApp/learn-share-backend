package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
)

type DbConfig struct {
	UserName string // `env:"POSTGRES_USER" env-default:"root"`
	Password string // `env:"POSTGRES_PASSWORD" env-default:"123"`
	Host     string // `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     int    // `env:"POSTGRES_PORT" env-default:"5432"`
	DbName   string // `env:"POSTGRES_DB" env-default:"yandex"`
}

type DB struct {
	ConnPool      *pgxpool.Pool
	closeConnOnce sync.Once
}

func New(ctx context.Context, config *DbConfig) (*DB, error) {
	var dbURL string = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		config.UserName,
		config.Password,
		config.Host,
		config.Port,
		config.DbName,
	)

	conn, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &DB{
		ConnPool:      conn,
		closeConnOnce: sync.Once{},
	}, nil
}

func (db *DB) Ping(ctx context.Context) error {
	return db.ConnPool.Ping(ctx)
}

func (db *DB) ClosePoolConn() {
	db.closeConnOnce.Do(db.ConnPool.Close)
}
