package rest

import (
	"context"
	"fmt"
	"github.com/LearnShareApp/learn-share-backend/internal/service/jwt"
	"net/http"
	"time"

	"github.com/LearnShareApp/learn-share-backend/internal/transport/rest/middlewares"
	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/login"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/registration"
	"github.com/go-chi/chi/v5"
)

const (
	defaultHTTPServerWriteTimeout = time.Second * 15
	defaultHTTPServerReadTimeout  = time.Second * 15
)

type ServerConfig struct {
	Port int
}

type Services struct {
	RegSrv     *registration.Service
	LoginSrv   *login.Service
	JwtService *jwt.Service
}

type Server struct {
	server *http.Server
	logger *zap.Logger
}

func NewServer(services *Services, config ServerConfig, log *zap.Logger) *Server {
	router := chi.NewRouter()

	router.Use(middlewares.LoggerMiddleware(log.Named("log_middleware")))

	regHandler := registration.MakeHandler(services.RegSrv, log)
	loginHandler := login.MakeHandler(services.LoginSrv, log)

	apiRouter := chi.NewRouter()

	apiRouter.Post("/signup", regHandler)
	apiRouter.Post("/login", loginHandler)

	apiRouter.Group(func(apiRouter chi.Router) {
		apiRouter.Use(middlewares.JWTMiddleware(services.JwtService, log.Named("jwt_middleware")))
		//apiRouter.Post("/manage", CreateAsset)
	})

	router.Mount("/api", apiRouter)

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
