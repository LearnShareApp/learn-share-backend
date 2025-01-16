package application

import (
	"context"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/config"
	"github.com/LearnShareApp/learn-share-backend/internal/repository"
	"github.com/LearnShareApp/learn-share-backend/internal/service/jwt"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/auth/login"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/auth/registration"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/categories/get_categories"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/teachers/add_skill"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/teachers/become_teacher"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/users/get_profile"
	"github.com/LearnShareApp/learn-share-backend/pkg/db/postgres"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"os"
)

type Application struct {
	context context.Context
	db      *sqlx.DB
	server  *rest.Server
	log     *zap.Logger
}

func New(ctx context.Context, config config.Config, log *zap.Logger) (*Application, error) {

	// db connection
	db, err := postgres.New(ctx, &config.Db)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %v", err)
	}

	log.Info("connected to database successfully")

	// TODO: repo, services, rest-server

	repo := repository.New(db)
	if config.IsInitDb {
		err = repo.CreateTables(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create tables: %v", err)
		}
		log.Info("successfully created tables (if they not existed)")
	}

	jwtService := jwt.NewJwtService(config.SecretKey, jwt.WithIssuer("learn-share-backend"))

	registrationSrv := registration.NewService(repo, jwtService)
	loginSrv := login.NewService(repo, jwtService)

	getCategoriesSrv := get_categories.NewService(repo)
	getProfileSrv := get_profile.NewService(repo)
	becomeTeacherSrv := become_teacher.NewService(repo)
	addSkillSrv := add_skill.NewService(repo)

	services := rest.NewServices(jwtService,
		registrationSrv,
		loginSrv,
		getCategoriesSrv,
		getProfileSrv,
		becomeTeacherSrv,
		addSkillSrv)

	restServer := rest.NewServer(services, config.Server, log)

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

	app.log.Info("shutting down application...")

	go func() {
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			app.log.Error("graceful shutdown timed out... forcing exit")
			os.Exit(1)
		}
	}()

	if err := app.server.GracefulStop(ctx); err != nil {
		return err
	}

	if err := app.db.Close(); err != nil {
		app.log.Error("failed to close database", zap.Error(err))
	}

	return nil
}
