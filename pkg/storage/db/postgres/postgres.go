package postgres

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type DBConfig struct {
	Host     string
	Port     int
	DbName   string
	User     string
	Password string
}

func New(ctx context.Context, config *DBConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%d",
		config.User,
		config.Password,
		config.DbName,
		config.Host,
		config.Port,
	)

	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
