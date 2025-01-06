package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
)

type DbConfig struct {
	Host     string
	Port     int
	DbName   string
	User     string
	Password string
}

type DB struct {
	ConnPool      *pgxpool.Pool
	closeConnOnce sync.Once
}

func New(ctx context.Context, config *DbConfig) (*DB, error) {
	var dbURL string = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		config.User,
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
