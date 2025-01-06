package rest

import (
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"

	"github.com/LearnShareApp/learn-share-backend/internal/use_cases/registration"
	"github.com/gorilla/mux"
)

const (
	defaultHTTPServerWriteTimeout = time.Second * 15
	defaultHTTPServerReadTimeout  = time.Second * 15
)

type ServerConfig struct {
	Port int
}

type Services struct {
	RegSrv *registration.Service
}

type Server struct {
	server *http.Server
	logger *zap.Logger
}

func NewServer(services *Services, config ServerConfig, log *zap.Logger) *Server {
	router := mux.NewRouter()

	router.Use()

	regHandler := registration.MakeHandler(services.RegSrv, log.Named("registration_service"))

	apiRouter := mux.NewRouter().PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/signup", regHandler).Methods("POST")

	router.PathPrefix("/api").Handler(apiRouter)

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
		s.logger.Info("starting Rest server", zap.String("port", s.server.Addr))
		return s.server.ListenAndServe()
	})

	return eg.Wait()
}
