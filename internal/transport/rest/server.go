package rest

import (
	"fmt"
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
	RegSrv   *registration.Service
	LoginSrv *login.Service
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
