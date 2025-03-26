package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
	"github.com/LearnShareApp/learn-share-backend/pkg/livekit"
	"github.com/LearnShareApp/learn-share-backend/pkg/migrator"
	"os"

	"github.com/LearnShareApp/learn-share-backend/internal/config"
	"github.com/LearnShareApp/learn-share-backend/internal/repository"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/categories/get_categories"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/image/get_image"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/lessons/approve_lesson"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/lessons/book_lesson"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/lessons/cancel_lesson"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/lessons/finish_lesson"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/lessons/get_lesson"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/lessons/get_lesson_shortdata"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/lessons/get_student_lessons"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/lessons/get_teacher_lessons"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/lessons/join_lesson"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/lessons/start_lesson"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/reviews/add_review"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/reviews/get_reviews"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/schedules/add_time"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/schedules/get_times"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/teachers/add_skill"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/teachers/become_teacher"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/teachers/get_teacher"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/teachers/get_teachers"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/users/auth/login"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/users/auth/registration"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/users/edit_user"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/users/get_user"
	"github.com/LearnShareApp/learn-share-backend/pkg/storage/db/postgres"
	"github.com/LearnShareApp/learn-share-backend/pkg/storage/object/minio"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Application struct {
	db     *sqlx.DB
	server *rest.Server
	log    *zap.Logger
}

func New(ctx context.Context, config *config.Config, log *zap.Logger) (*Application, error) {
	// database connection
	database, err := postgres.New(ctx, &config.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	log.Info("connected to database successfully")

	minioClient, err := minio.NewClient(&config.Minio)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to minio: %w", err)
	}

	if err = minio.CreateBucket(ctx, minioClient, config.Minio.Bucket); err != nil {
		return nil, fmt.Errorf("failed to create minio bucket: %w", err)
	}

	log.Info("connected to minio successfully")

	repo := repository.New(database)

	if config.IsInitDb {
		err = migrator.RunMigrations(&config.Migrator)
		if err != nil {
			return nil, err //nolint:wrapcheck
		}

		log.Info("up migrations successfully")
	}

	jwtService := jwt.NewService(config.JwtSecretKey, jwt.WithIssuer("learn-share-backend"))
	leveKitService := livekit.NewService(config.LiveKit)
	minioService := minio.NewService(minioClient, config.Minio.Bucket)

	var (
		registrationSrv         = registration.NewService(repo, jwtService, minioService)
		loginSrv                = login.NewService(repo, jwtService)
		getCategoriesSrv        = get_categories.NewService(repo)
		getProfileSrv           = get_user.NewService(repo)
		editUserSrv             = edit_user.NewService(repo, minioService)
		becomeTeacherSrv        = become_teacher.NewService(repo)
		addSkillSrv             = add_skill.NewService(repo)
		getTeacherSrv           = get_teacher.NewService(repo)
		getTeachersSrv          = get_teachers.NewService(repo)
		addScheduleTimeSrv      = add_time.NewService(repo)
		getScheduleTimesSrv     = get_times.NewService(repo)
		bookLessonSrv           = book_lesson.NewService(repo)
		getLessonsForTeacherSrv = get_teacher_lessons.NewService(repo)
		getLessonsForStudentSrv = get_student_lessons.NewService(repo)
		getLessonSrv            = get_lesson.NewService(repo)
		getLessonShortDataSrv   = get_lesson_shortdata.NewService(repo)
		cancelLessonSrv         = cancel_lesson.NewService(repo)
		approveLessonSrv        = approve_lesson.NewService(repo)
		startLessonSrv          = start_lesson.NewService(repo, leveKitService)
		finishLessonSrv         = finish_lesson.NewService(repo)
		joinLessonSrv           = join_lesson.NewService(repo, leveKitService)
		getImagSrv              = get_image.NewService(minioService)
		addReviewSrv            = add_review.NewService(repo)
		getReviewsSrv           = get_reviews.NewService(repo)
	)

	services := rest.NewServices(jwtService,
		registrationSrv,
		loginSrv,
		getCategoriesSrv,
		getProfileSrv,
		editUserSrv,
		becomeTeacherSrv,
		addSkillSrv,
		getTeacherSrv,
		getTeachersSrv,
		addScheduleTimeSrv,
		getScheduleTimesSrv,
		bookLessonSrv,
		getLessonsForTeacherSrv,
		getLessonsForStudentSrv,
		getLessonSrv,
		getLessonShortDataSrv,
		cancelLessonSrv,
		approveLessonSrv,
		startLessonSrv,
		finishLessonSrv,
		joinLessonSrv,
		getImagSrv,
		addReviewSrv,
		getReviewsSrv,
	)

	restServer := rest.NewServer(services, config.Server, log)

	return &Application{
		db:     database,
		server: restServer,
		log:    log,
	}, nil
}

// Run start application.
func (app *Application) Run() error {
	err := app.server.Start()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Shutdown gracefully stop application.
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
		return err // nolint:wrapcheck
	}

	if err := app.db.Close(); err != nil {
		app.log.Error("failed to close database", zap.Error(err))
	}

	return nil
}
