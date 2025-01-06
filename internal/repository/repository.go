package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	dbPool *pgxpool.Pool
}

func New(pgxPool *pgxpool.Pool) *Repository {
	return &Repository{
		dbPool: pgxPool,
	}
}
