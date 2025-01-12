package rest

import (
	"context"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/service/jwt"
	get_categories2 "github.com/LearnShareApp/learn-share-backend/internal/use_cases/get_categories"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/user/get_profile"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"time"

	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/middlewares"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/auth/login"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/auth/registration"
	"github.com/go-chi/chi/v5"

	_ "github.com/LearnShareApp/learn-share-backend/docs"
)

const (
	defaultHTTPServerWriteTimeout = time.Second * 15
	defaultHTTPServerReadTimeout  = time.Second * 15

	authRoute = "/auth"
	userRoute = "/user"
	apiRoute  = "/api"
)

type ServerConfig struct {
	Port int
}

type Services struct {
	JwtSrv           *jwt.Service
	RegSrv           *registration.Service
	LoginSrv         *login.Service
	GetCategoriesSrv *get_categories2.Service
	GetProfileSrv    *get_profile.Service
}

type Server struct {
	server *http.Server
	logger *zap.Logger
}

func NewServices(jwtSrv *jwt.Service,
	reg *registration.Service,
	login *login.Service,
	getCategories *get_categories2.Service,
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

	regHandler := registration.MakeHandler(services.RegSrv, log)
	loginHandler := login.MakeHandler(services.LoginSrv, log)
	getCategoriesHandler := get_categories2.MakeHandler(services.GetCategoriesSrv, log)
	getProfileHandler := get_profile.MakeHandler(services.GetProfileSrv, log)

	apiRouter := chi.NewRouter()

	// auth routes
	authRouter := chi.NewRouter()
	authRouter.Post(registration.Route, regHandler)
	authRouter.Post(login.Route, loginHandler)
	apiRouter.Mount(authRoute, authRouter)

	// Categories
	apiRouter.Get(get_categories2.Route, getCategoriesHandler)

	apiRouter.Group(func(apiRouter chi.Router) {
		apiRouter.Use(middlewares.JWTMiddleware(services.JwtSrv, log.Named("jwt_middleware")))

		// user routes
		userRouter := chi.NewRouter()
		userRouter.Get(get_profile.Route, getProfileHandler)
		apiRouter.Mount(userRoute, userRouter)

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
	eg := errgroup.Group{}

	eg.Go(func() error {
		s.logger.Info("starting Rest server", zap.String("address", s.server.Addr))
		return s.server.ListenAndServe()
	})

	return eg.Wait()
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
