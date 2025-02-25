package rest

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"time"

	_ "github.com/LearnShareApp/learn-share-backend/docs"
	"github.com/LearnShareApp/learn-share-backend/internal/service/jwt"
	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/middlewares"
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

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
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
	lessonsRoute  = "/lessons"
	apiRoute      = "/api"
	studentRoute  = "/student"
	reviewRoute   = "/review"
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
	EditProfileSrv          *edit_user.Service
	BecomeTeacherSrv        *become_teacher.Service
	AddSkillSrv             *add_skill.Service
	GetTeacherSrv           *get_teacher.Service
	GetTeachersSrv          *get_teachers.Service
	AddScheduleTimeSrv      *add_time.Service
	GetScheduleTimesSrv     *get_times.Service
	BookLessonSrv           *book_lesson.Service
	GetLessonsForTeacherSrv *get_teacher_lessons.Service
	GetLessonsForStudentSrv *get_student_lessons.Service
	GetLessonSrv            *get_lesson.Service
	GetLessonSortSrv        *get_lesson_shortdata.Service
	CancelLessonSrv         *cancel_lesson.Service
	ApproveLessonSrv        *approve_lesson.Service
	StartLessonSrv          *start_lesson.Service
	FinishLessonSrv         *finish_lesson.Service
	JoinLessonSrv           *join_lesson.Service
	GetImageSrv             *get_image.Service
	AddReviewSrv            *add_review.Service
	GetReviewsSrv           *get_reviews.Service
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
	editProfileSrv *edit_user.Service,
	becomeTeacherSrv *become_teacher.Service,
	addSkillSrv *add_skill.Service,
	getTeacher *get_teacher.Service,
	getTeachers *get_teachers.Service,
	addScheduleTimeSrv *add_time.Service,
	getScheduleTimeSrv *get_times.Service,
	bookLessonSrv *book_lesson.Service,
	getLessonsForTeacherSrv *get_teacher_lessons.Service,
	getLessonsForStudentSrv *get_student_lessons.Service,
	getLessonSrv *get_lesson.Service,
	getLessonSortSrv *get_lesson_shortdata.Service,
	cancelLessonSrv *cancel_lesson.Service,
	approveLessonSrv *approve_lesson.Service,
	startLessonSrv *start_lesson.Service,
	finishLessonSrv *finish_lesson.Service,
	joinLesson *join_lesson.Service,
	getImageSrv *get_image.Service,
	addReviewSrv *add_review.Service,
	getReviewsSrv *get_reviews.Service,
) *Services {
	return &Services{
		JwtSrv: jwtSrv,
		RegSrv: reg,

		LoginSrv:                login,
		GetCategoriesSrv:        getCategories,
		GetProfileSrv:           getProfile,
		EditProfileSrv:          editProfileSrv,
		BecomeTeacherSrv:        becomeTeacherSrv,
		AddSkillSrv:             addSkillSrv,
		GetTeacherSrv:           getTeacher,
		GetTeachersSrv:          getTeachers,
		AddScheduleTimeSrv:      addScheduleTimeSrv,
		GetScheduleTimesSrv:     getScheduleTimeSrv,
		BookLessonSrv:           bookLessonSrv,
		GetLessonsForTeacherSrv: getLessonsForTeacherSrv,
		GetLessonsForStudentSrv: getLessonsForStudentSrv,
		GetLessonSrv:            getLessonSrv,
		GetLessonSortSrv:        getLessonSortSrv,
		CancelLessonSrv:         cancelLessonSrv,
		ApproveLessonSrv:        approveLessonSrv,
		StartLessonSrv:          startLessonSrv,
		FinishLessonSrv:         finishLessonSrv,
		JoinLessonSrv:           joinLesson,
		GetImageSrv:             getImageSrv,
		AddReviewSrv:            addReviewSrv,
		GetReviewsSrv:           getReviewsSrv,
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

	// image route
	apiRouter.Get(get_image.Route, get_image.MakeHandler(services.GetImageSrv, log))

	// users route
	usersRouter := chi.NewRouter()
	usersRouter.Get(get_user.PublicRoute, get_user.MakePublicHandler(services.GetProfileSrv, log))
	apiRouter.Mount(usersRoute, usersRouter)

	// teachers route
	teachersRouter := chi.NewRouter()
	teachersRouter.Get(get_teacher.PublicRoute, get_teacher.MakePublicHandler(services.GetTeacherSrv, log))
	teachersRouter.Get(get_times.PublicRoute, get_times.MakePublicHandler(services.GetScheduleTimesSrv, log))
	teachersRouter.Get(get_reviews.Route, get_reviews.MakeHandler(services.GetReviewsSrv, log))
	apiRouter.Mount(teachersRoute, teachersRouter)

	// lessons route
	lessonsRouter := chi.NewRouter()
	lessonsRouter.Get(get_lesson.Route, get_lesson.MakeHandler(services.GetLessonSrv, log))
	lessonsRouter.Get(get_lesson_shortdata.Route, get_lesson_shortdata.MakeHandler(services.GetLessonSortSrv, log))
	apiRouter.Mount(lessonsRoute, lessonsRouter)

	// protected routes
	apiRouter.Group(func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware(services.JwtSrv, log.Named("jwt_middleware")))

		// protected routes
		r.Get(path.Join(userRoute, get_user.ProtectedRoute), get_user.MakeProtectedHandler(services.GetProfileSrv, log))
		r.Get(path.Join(teacherRoute, get_teacher.ProtectedRoute), get_teacher.MakeProtectedHandler(services.GetTeacherSrv, log))
		r.Get(path.Join(teacherRoute, get_times.ProtectedRoute), get_times.MakeProtectedHandler(services.GetScheduleTimesSrv, log))
		r.Get(path.Join(teachersRoute, get_teachers.Route), get_teachers.MakeHandler(services.GetTeachersSrv, log))
		r.Get(path.Join(teacherRoute, get_teacher_lessons.Route), get_teacher_lessons.MakeHandler(services.GetLessonsForTeacherSrv, log))
		r.Get(path.Join(studentRoute, get_student_lessons.Route), get_student_lessons.MakeHandler(services.GetLessonsForStudentSrv, log))
		r.Get(path.Join(lessonsRoute, join_lesson.Route), join_lesson.MakeHandler(services.JoinLessonSrv, log))

		r.Post(path.Join(teacherRoute, become_teacher.Route), become_teacher.MakeHandler(services.BecomeTeacherSrv, log))
		r.Post(path.Join(teacherRoute, add_skill.Route), add_skill.MakeHandler(services.AddSkillSrv, log))
		r.Post(path.Join(teacherRoute, add_time.Route), add_time.MakeHandler(services.AddScheduleTimeSrv, log))
		r.Post(path.Join(lessonRoute, book_lesson.Route), book_lesson.MakeHandler(services.BookLessonSrv, log))
		r.Post(path.Join(reviewRoute, add_review.Route), add_review.MakeHandler(services.AddReviewSrv, log))

		r.Put(path.Join(lessonsRoute, cancel_lesson.Route), cancel_lesson.MakeHandler(services.CancelLessonSrv, log))
		r.Put(path.Join(lessonsRoute, approve_lesson.Route), approve_lesson.MakeHandler(services.ApproveLessonSrv, log))
		r.Put(path.Join(lessonsRoute, start_lesson.Route), start_lesson.MakeHandler(services.StartLessonSrv, log))
		r.Put(path.Join(lessonsRoute, finish_lesson.Route), finish_lesson.MakeHandler(services.FinishLessonSrv, log))

		r.Patch(path.Join(userRoute, edit_user.Route), edit_user.MakeHandler(services.EditProfileSrv, log))
	})

	router.Mount(apiRoute, apiRouter)

	// Добавляем swagger endpoint
	router.Get("/swagger/*", httpSwagger.Handler(
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
	// eg := errgroup.Group{}
	//
	//eg.Go(func() error {
	//	s.logger.Info("starting Rest server", zap.String("address", s.server.Addr))
	//	return s.server.ListenAndServe()
	//})
	s.logger.Info("starting Rest server", zap.String("address", s.server.Addr))

	return s.server.ListenAndServe()
}

// GracefulStop корректная остановка сервера.
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
