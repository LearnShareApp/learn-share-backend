package application

import (
	"context"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/config"
	"github.com/LearnShareApp/learn-share-backend/internal/repository"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/registration"
	"github.com/LearnShareApp/learn-share-backend/pkg/db/postgres"
	"go.uber.org/zap"
)

type Application struct {
	context context.Context
	db      *postgres.DB
	server  *rest.Server
	log     *zap.Logger
}

func New(ctx context.Context, config config.Config, log *zap.Logger) (*Application, error) {

	// db connection
	db, err := postgres.New(ctx, &config.Db)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %v", err)
	}
	if err = db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping to postgres: %v", err)
	}

	log.Info("connected to database successfully")

	// TODO: repo, services, rest-server

	repo := repository.New(db.ConnPool)

	//jwtService := jwt.NewJwtService(config.Jwt.SecretKey, jwt.WithIssuer("learn-share-backend"))

	registrationSrv := registration.NewService(repo)

	services := &rest.Services{
		RegSrv: registrationSrv,
	}

	restServer := rest.NewServer(services, config.Rest, log.Named("rest_server"))

	return &Application{
		context: ctx,
		db:      db,
		server:  restServer,
		log:     log,
	}, nil
}

// Run запускает приложение
func (app *Application) Run() error {
	return app.server.Start()
}

// Shutdown gracefully останавливает приложение
func (app *Application) Shutdown(ctx context.Context) error {
	app.db.ClosePoolConn()
	// TODO stop rest server, close connection with db
	return nil
}
