package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateTables(ctx context.Context) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	defer tx.Rollback()

	if err := createUsersTable(ctx, tx); err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func createUsersTable(ctx context.Context, tx *sqlx.Tx) error {
	const query = `
    CREATE TABLE IF NOT EXISTS public.users(
        id SERIAL PRIMARY KEY NOT NULL,
        email TEXT UNIQUE NOT NULL,
        name TEXT NOT NULL,
        surname TEXT NOT NULL,
        password TEXT NOT NULL,
        registration_date TIMESTAMPTZ DEFAULT NOW(),
        birthdate DATE NOT NULL,
        avatar TEXT
    );
    `

	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to execute users table creation: %w", err)
	}
	return nil
}
