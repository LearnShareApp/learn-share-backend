package repository

import (
	"context"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/entities"
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

	if err = createUsersTable(ctx, tx); err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}

	if err = createTeachersTable(ctx, tx); err != nil {
		return fmt.Errorf("error creating teachers table: %w", err)
	}

	if err = createCategoriesTable(ctx, tx); err != nil {
		return fmt.Errorf("error creating categories table: %w", err)
	}

	categories := []entities.Category{
		{Name: "Cooking", MinAge: 7},
		{Name: "Programming", MinAge: 7},
		{Name: "Drawing", MinAge: 0},
		{Name: "Dancing", MinAge: 0},
	}

	if err = seedCategories(ctx, tx, categories); err != nil {
		return fmt.Errorf("error seeding categories: %w", err)
	}

	if err = createSkillsTable(ctx, tx); err != nil {
		return fmt.Errorf("error creating skills table: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func createUsersTable(ctx context.Context, tx *sqlx.Tx) error {
	const query = `
    CREATE TABLE IF NOT EXISTS public.users(
        user_id SERIAL PRIMARY KEY NOT NULL,
        email TEXT UNIQUE NOT NULL,
        name TEXT NOT NULL,
        surname TEXT NOT NULL,
        password TEXT NOT NULL,
        registration_date TIMESTAMPTZ DEFAULT NOW(),
        birthdate DATE NOT NULL,
        avatar TEXT NOT NULL DEFAULT ''
    );
    `

	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to execute users table creation: %w", err)
	}
	return nil
}

func createTeachersTable(ctx context.Context, tx *sqlx.Tx) error {
	const query = `
    CREATE TABLE IF NOT EXISTS public.teachers(
        teacher_id SERIAL PRIMARY KEY NOT NULL,
        user_id INTEGER UNIQUE NOT NULL REFERENCES users(user_id) ON DELETE CASCADE
    );
    `

	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to execute users table creation: %w", err)
	}
	return nil
}

func createCategoriesTable(ctx context.Context, tx *sqlx.Tx) error {
	const query = `
    CREATE TABLE IF NOT EXISTS public.categories(
        category_id SERIAL PRIMARY KEY NOT NULL,
        name TEXT UNIQUE NOT NULL,
        min_age INTEGER NOT NULL
    );
    `

	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to execute categories table creation: %w", err)
	}
	return nil
}

func createSkillsTable(ctx context.Context, tx *sqlx.Tx) error {
	const query = `
	CREATE TABLE IF NOT EXISTS public.skills (
		skill_id SERIAL PRIMARY KEY,
		teacher_id INTEGER NOT NULL REFERENCES teachers(teacher_id) ON DELETE CASCADE,
		category_id INTEGER NOT NULL REFERENCES categories(category_id) ON DELETE CASCADE, 
		video_card_link TEXT,
		about TEXT,
		rate SMALLINT NOT NULL DEFAULT 0,
		is_active BOOLEAN NOT NULL DEFAULT TRUE, -- по хорошему FALSE но это если делать механизм подтверждения
		CONSTRAINT unique_teacher_category UNIQUE (teacher_id, category_id) -- Уникальность teacher_id и category_id
	);

	`

	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to execute skills table creation: %w", err)
	}

	return nil
}

func seedCategories(ctx context.Context, tx *sqlx.Tx, categories []entities.Category) error {
	const query = `
        INSERT INTO public.categories (name, min_age)
        VALUES ($1, $2)
        ON CONFLICT (name) DO NOTHING
    `

	for _, category := range categories {
		_, err := tx.ExecContext(ctx, query, category.Name, category.MinAge)
		if err != nil {
			return fmt.Errorf("failed to insert category %s: %w", category.Name, err)
		}
	}

	return nil
}
