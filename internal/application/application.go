package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/config"
	"github.com/LearnShareApp/learn-share-backend/internal/repository"
	"github.com/LearnShareApp/learn-share-backend/internal/service/jwt"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/auth/login"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/auth/registration"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/categories/get_categories"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/lessons/book_lesson"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/schedules/add_time"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/schedules/get_times"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/teachers/add_skill"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/teachers/become_teacher"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/teachers/get_teacher"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/teachers/get_teachers"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/users/get_user"
	"github.com/LearnShareApp/learn-share-backend/pkg/db/postgres"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"os"
	"time"
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

	jwtService := jwt.NewJwtService(config.SecretKey, jwt.WithIssuer("learn-share-backend"), jwt.WithDuration(time.Hour*24*7))

	var (
		registrationSrv     = registration.NewService(repo, jwtService)
		loginSrv            = login.NewService(repo, jwtService)
		getCategoriesSrv    = get_categories.NewService(repo)
		getProfileSrv       = get_user.NewService(repo)
		becomeTeacherSrv    = become_teacher.NewService(repo)
		addSkillSrv         = add_skill.NewService(repo)
		getTeacherSrv       = get_teacher.NewService(repo)
		getTeachersSrv      = get_teachers.NewService(repo)
		addScheduleTimeSrv  = add_time.NewService(repo)
		getScheduleTimesSrv = get_times.NewService(repo)
		bookLessonSrv       = book_lesson.NewService(repo)
	)

	services := rest.NewServices(jwtService,
		registrationSrv,
		loginSrv,
		getCategoriesSrv,
		getProfileSrv,
		becomeTeacherSrv,
		addSkillSrv,
		getTeacherSrv,
		getTeachersSrv,
		addScheduleTimeSrv,
		getScheduleTimesSrv,
		bookLessonSrv)

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
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
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
