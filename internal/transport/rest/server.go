package rest

import (
	"context"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/service/jwt"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/categories/get_categories"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/lessons/book_lesson"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/schedules/add_time"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/schedules/get_times"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/teachers/add_skill"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/teachers/become_teacher"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/teachers/get_lessons"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/teachers/get_teacher"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/teachers/get_teachers"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/users/get_user"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"path"
	"time"

	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/middlewares"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/auth/login"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/auth/registration"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	_ "github.com/LearnShareApp/learn-share-backend/docs"
)

const (
	defaultHTTPServerWriteTimeout = time.Second * 15
	defaultHTTPServerReadTimeout  = time.Second * 15

	authRoute     = "/auth"
	userRoute     = "/user"
	usersRoute    = "/users"
	teacherRoute  = "/teacher"
	teachersRoute = "/teachers"
	lessonRoute   = "/lesson"
	apiRoute      = "/api"
)

type ServerConfig struct {
	Port int
}

type Services struct {
	JwtSrv                  *jwt.Service
	RegSrv                  *registration.Service
	LoginSrv                *login.Service
	GetCategoriesSrv        *get_categories.Service
	GetProfileSrv           *get_user.Service
	BecomeTeacherSrv        *become_teacher.Service
	AddSkillSrv             *add_skill.Service
	GetTeacherSrv           *get_teacher.Service
	GetTeachersSrv          *get_teachers.Service
	AddScheduleTimeSrv      *add_time.Service
	GetScheduleTimesSrv     *get_times.Service
	BookLessonSrv           *book_lesson.Service
	GetLessonsForTeacherSrv *get_lessons.Service
}

type Server struct {
	server *http.Server
	logger *zap.Logger
}

func NewServices(jwtSrv *jwt.Service,
	reg *registration.Service,
	login *login.Service,
	getCategories *get_categories.Service,
	getProfile *get_user.Service,
	becomeTeacherSrv *become_teacher.Service,
	addSkillSrv *add_skill.Service,
	getTeacher *get_teacher.Service,
	getTeachers *get_teachers.Service,
	addScheduleTimeSrv *add_time.Service,
	getScheduleTimeSrv *get_times.Service,
	bookLessonSrv *book_lesson.Service,
	getLessonsForTeacherSrv *get_lessons.Service) *Services {
	return &Services{
		JwtSrv:                  jwtSrv,
		RegSrv:                  reg,
		LoginSrv:                login,
		GetCategoriesSrv:        getCategories,
		GetProfileSrv:           getProfile,
		BecomeTeacherSrv:        becomeTeacherSrv,
		AddSkillSrv:             addSkillSrv,
		GetTeacherSrv:           getTeacher,
		GetTeachersSrv:          getTeachers,
		AddScheduleTimeSrv:      addScheduleTimeSrv,
		GetScheduleTimesSrv:     getScheduleTimeSrv,
		BookLessonSrv:           bookLessonSrv,
		GetLessonsForTeacherSrv: getLessonsForTeacherSrv,
	}
}

func NewServer(services *Services, config ServerConfig, log *zap.Logger) *Server {
	router := chi.NewRouter()

	router.Use(middlewares.LoggerMiddleware(log.Named("log_middleware")))
	router.Use(middlewares.CorsMiddleware)

	// root router
	apiRouter := chi.NewRouter()

	// public rotes

	// auth routes
	authRouter := chi.NewRouter()
	authRouter.Post(registration.Route, registration.MakeHandler(services.RegSrv, log))
	authRouter.Post(login.Route, login.MakeHandler(services.LoginSrv, log))
	apiRouter.Mount(authRoute, authRouter)

	// categories route
	apiRouter.Get(get_categories.Route, get_categories.MakeHandler(services.GetCategoriesSrv, log))

	// users route
	usersRouter := chi.NewRouter()
	usersRouter.Get(get_user.PublicRoute, get_user.MakePublicHandler(services.GetProfileSrv, log))

	// teachers route
	teachersRouter := chi.NewRouter()
	teachersRouter.Get(get_teacher.PublicRoute, get_teacher.MakePublicHandler(services.GetTeacherSrv, log))
	teachersRouter.Get(get_times.PublicRoute, get_times.MakePublicHandler(services.GetScheduleTimesSrv, log))
	apiRouter.Mount(teachersRoute, teachersRouter)

	// protected routes
	apiRouter.Group(func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware(services.JwtSrv, log.Named("jwt_middleware")))

		// protected routes
		r.Get(path.Join(userRoute, get_user.ProtectedRoute), get_user.MakeProtectedHandler(services.GetProfileSrv, log))
		r.Get(path.Join(teacherRoute, get_teacher.ProtectedRoute), get_teacher.MakeProtectedHandler(services.GetTeacherSrv, log))
		r.Get(path.Join(teacherRoute, get_times.ProtectedRoute), get_times.MakeProtectedHandler(services.GetScheduleTimesSrv, log))
		r.Get(path.Join(teachersRoute, get_teachers.Route), get_teachers.MakeHandler(services.GetTeachersSrv, log))
		r.Get(path.Join(teacherRoute, get_lessons.Route), get_lessons.MakeHandler(services.GetLessonsForTeacherSrv, log))

		r.Post(path.Join(teacherRoute, become_teacher.Route), become_teacher.MakeHandler(services.BecomeTeacherSrv, log))
		r.Post(path.Join(teacherRoute, add_skill.Route), add_skill.MakeHandler(services.AddSkillSrv, log))
		r.Post(path.Join(teacherRoute, add_time.Route), add_time.MakeHandler(services.AddScheduleTimeSrv, log))
		r.Post(path.Join(lessonRoute, book_lesson.Route), book_lesson.MakeHandler(services.BookLessonSrv, log))
	})

	apiRouter.Mount(usersRoute, usersRouter)

	router.Mount(apiRoute, apiRouter)

	// Добавляем swagger endpoint
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("./swagger/doc.json"), // URL указывающий на JSON документацию
	))
	router.Get("/*", httpSwagger.Handler(
		httpSwagger.URL("./swagger/doc.json"), // URL указывающий на JSON документацию
	))

	return &Server{
		server: &http.Server{
			Handler:      router,
			Addr:         fmt.Sprintf(":%d", config.Port),
			WriteTimeout: defaultHTTPServerWriteTimeout,
			ReadTimeout:  defaultHTTPServerReadTimeout,
		},
		logger: log,
	}
}

func (s *Server) Start() error {
	//eg := errgroup.Group{}
	//
	//eg.Go(func() error {
	//	s.logger.Info("starting Rest server", zap.String("address", s.server.Addr))
	//	return s.server.ListenAndServe()
	//})
	s.logger.Info("starting Rest server", zap.String("address", s.server.Addr))
	return s.server.ListenAndServe()
}

// GracefulStop корректная остановка сервера
func (s *Server) GracefulStop(ctx context.Context) error {
	// Создаем контекст с таймаутом
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	s.logger.Info("shutting down Rest server", zap.String("address", s.server.Addr))

	// Остановка сервера
	if err := s.server.Shutdown(shutdownCtx); err != nil {
		err = fmt.Errorf("failed to shutdown rest server: %w", err)
		return err
	}

	s.logger.Info("rest server stopped")
	return nil
}
