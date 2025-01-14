package rest

import (
	"context"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/service/jwt"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/categories/get_categories"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/users/get_profile"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/sync/errgroup"
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

	authRoute  = "/auth"
	userRoute  = "/user"
	usersRoute = "/users"
	apiRoute   = "/api"
)

type ServerConfig struct {
	Port int
}

type Services struct {
	JwtSrv           *jwt.Service
	RegSrv           *registration.Service
	LoginSrv         *login.Service
	GetCategoriesSrv *get_categories.Service
	GetProfileSrv    *get_profile.Service
}

type Server struct {
	server *http.Server
	logger *zap.Logger
}

func NewServices(jwtSrv *jwt.Service,
	reg *registration.Service,
	login *login.Service,
	getCategories *get_categories.Service,
	getProfile *get_profile.Service) *Services {
	return &Services{
		JwtSrv:           jwtSrv,
		RegSrv:           reg,
		LoginSrv:         login,
		GetCategoriesSrv: getCategories,
		GetProfileSrv:    getProfile,
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
	usersRouter.Get(get_profile.PublicRoute, get_profile.MakePublicHandler(services.GetProfileSrv, log))

	// protected routes
	apiRouter.Group(func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware(services.JwtSrv, log.Named("jwt_middleware")))

		// protected routes
		r.Get(path.Join(userRoute, get_profile.ProtectedRoute), get_profile.MakeProtectedHandler(services.GetProfileSrv, log))

	})

	apiRouter.Mount(usersRoute, usersRouter)

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
	eg := errgroup.Group{}

	eg.Go(func() error {
		s.logger.Info("starting Rest server", zap.String("address", s.server.Addr))
		return s.server.ListenAndServe()
	})
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
