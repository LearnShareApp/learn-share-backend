package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LearnShareApp/learn-share-backend/internal/application"
	"github.com/LearnShareApp/learn-share-backend/internal/config"
	"github.com/LearnShareApp/learn-share-backend/pkg/logger"

	"go.uber.org/zap"
)

// @title        Learn-Share API
// @version      1.0
// @description  back-end part for mobile application.

// @contact.name   Ruslan's Support
// @contact.url    https://t.me/Ruslan20007
// @contact.email  ruslanrbb8@gmail.com

// @host         adoe.ru:81
// @BasePath     /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	ctx, ctxClose := context.WithCancel(context.Background())
	defer ctxClose()

	log := logger.NewDevelopment()
	defer func() {
		if err := log.Sync(); err != nil {
			log.Error("failed to sync logger", zap.Error(err))
		}
	}()

	envPaths := []string{
		".env",
		"./config/.env",
		"./internal/config/.env",
	}

	cfg, err := config.LoadConfig(envPaths)
	if err != nil {
		log.Error("failed to load config", zap.Error(err))

		return
	}

	marshaledCfg, err := cfg.LogConfig()
	if err != nil {
		log.Error("failed to marshal config to log", zap.Error(err))
	} else {
		log.Info(marshaledCfg)
	}

	app, err := application.New(ctx, cfg, log)
	if err != nil {
		log.Error("failed to create application", zap.Error(err))

		return
	}

	graceCh := make(chan os.Signal, 1)
	errCh := make(chan error, 1)

	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(graceCh) // clear subscribe to signal
	defer close(graceCh)       // close channel
	defer close(errCh)         // close error channel

	go func() {
		if err := app.Run(); err != nil {
			select {
			case errCh <- err: // Non blocked sending error to channel
			default:
				log.Error("failed to send error to channel", zap.Error(err))
			}
		}
	}()

	// waiting or signal to close or error
	select {
	case <-graceCh:
		log.Info("Received shutdown signal")

	case err := <-errCh:
		log.Error("Application error", zap.Error(err))
	}

	// create context with timeout to close
	const timeout = 10 * time.Second

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, timeout)
	defer shutdownCancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		log.Error("failed to shutdown", zap.Error(err))
	}

	log.Info("App stopped")
}
