package application

import (
	"context"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/config"
	"github.com/LearnShareApp/learn-share-backend/pkg/db/postgres"
)

type Application struct {
	context context.Context
}

func New(ctx context.Context, config *config.Config) (*Application, error) {

	db, err := postgres.New(ctx, &config.DbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %v", err)
	}
	if err = db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping to postgres: %v", err)
	}

	//repo := repository.New(db.ConnPool)

	// TODO: repo, services, rest-server
	return &Application{context: ctx}, nil
}

// Run запускает приложение
func (app *Application) Run() error {
	// TODO: run rest server
	return nil
}

// Shutdown gracefully останавливает приложение
func (app *Application) Shutdown(ctx context.Context) error {
	// TODO stop rest server, close connection with db
	return nil
}
