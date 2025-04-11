package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/service/admin"
	"os"

	"github.com/LearnShareApp/learn-share-backend/internal/config"
	"github.com/LearnShareApp/learn-share-backend/internal/repository"
	"github.com/LearnShareApp/learn-share-backend/internal/service/category"
	"github.com/LearnShareApp/learn-share-backend/internal/service/common"
	"github.com/LearnShareApp/learn-share-backend/internal/service/image"
	"github.com/LearnShareApp/learn-share-backend/internal/service/lesson"
	"github.com/LearnShareApp/learn-share-backend/internal/service/review"
	"github.com/LearnShareApp/learn-share-backend/internal/service/schedule"
	"github.com/LearnShareApp/learn-share-backend/internal/service/teacher"
	"github.com/LearnShareApp/learn-share-backend/internal/service/user"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest"
	"github.com/LearnShareApp/learn-share-backend/pkg/jwt"
	"github.com/LearnShareApp/learn-share-backend/pkg/livekit"
	"github.com/LearnShareApp/learn-share-backend/pkg/migrator"
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

type Services struct {
	jwt.JWTService
	user.UserService
	teacher.TeacherService
	schedule.ScheduleService
	review.ReviewService
	lesson.LessonService
	image.ImageService
	category.CategoryService
	admin.AdminService
	common.CommonService
}

func NewServices(
	jwtService *jwt.JWTService,
	userService *user.UserService,
	teacherService *teacher.TeacherService,
	scheduleService *schedule.ScheduleService,
	reviewService *review.ReviewService,
	lessonService *lesson.LessonService,
	imageService *image.ImageService,
	categoryService *category.CategoryService,
	adminService *admin.AdminService,
	commonService *common.CommonService,
) *Services {
	return &Services{
		JWTService:      *jwtService,
		UserService:     *userService,
		TeacherService:  *teacherService,
		ScheduleService: *scheduleService,
		ReviewService:   *reviewService,
		LessonService:   *lessonService,
		ImageService:    *imageService,
		CategoryService: *categoryService,
		AdminService:    *adminService,
		CommonService:   *commonService,
	}
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
			return nil, err
		}

		log.Info("up migrations successfully")
	}

	/*----------------------------------------------------------*/

	// services
	jwtService := jwt.NewService(config.JwtSecretKey, jwt.WithIssuer("learn-share-backend"))
	liveKitService := livekit.NewService(config.LiveKit)
	minioService := minio.NewService(minioClient, config.Minio.Bucket)

	userService := user.NewService(repo, minioService)
	teacherService := teacher.NewService(repo)
	scheduleService := schedule.NewService(repo)
	reviewService := review.NewService(repo)
	lessonService := lesson.NewService(repo, liveKitService)
	imageService := image.NewService(minioService)
	commonService := common.NewService(repo)
	categoryService := category.NewService(repo)
	adminService := admin.NewService(repo)

	services := NewServices(
		jwtService,
		userService,
		teacherService,
		scheduleService,
		reviewService,
		lessonService,
		imageService,
		categoryService,
		adminService,
		commonService,
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
		return err
	}

	if err := app.db.Close(); err != nil {
		app.log.Error("failed to close database", zap.Error(err))
	}

	return nil
}
