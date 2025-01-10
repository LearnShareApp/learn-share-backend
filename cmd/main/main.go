package main

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/application"
	"github.com/LearnShareApp/learn-share-backend/internal/config"
	"github.com/LearnShareApp/learn-share-backend/pkg/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, ctxClose := context.WithCancel(context.Background())
	defer ctxClose()

	log := logger.NewDevelopment()
	defer log.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("failed to load config", zap.Error(err))
		panic("failed to load config")
	}

	log.Info(cfg.LogConfig())

	app, err := application.New(ctx, *cfg, log)
	if err != nil {
		log.Error("failed to create application", zap.Error(err))
		return
	}

	graceCh := make(chan os.Signal, 1)
	errCh := make(chan error, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(graceCh) // Очищаем подписку на сигналы
	defer close(graceCh)       // Закрываем канал
	defer close(errCh)         // Закрываем канал с ошибками

	go func() {
		if err := app.Run(); err != nil {
			select {
			case errCh <- err: // Неблокирующая отправка ошибки
			default:
				log.Error("failed to send error to channel", zap.Error(err))
			}
		}
	}()

	// Ожидаем либо сигнал завершения, либо ошибку
	select {
	case <-graceCh:
		log.Info("Received shutdown signal")
	case err := <-errCh:
		log.Error("Application error", zap.Error(err))
	}

	// Создаем контекст с таймаутом для graceful shutdown
	const timeout = 10 * time.Second
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, timeout)
	defer shutdownCancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		log.Error("failed to shutdown", zap.Error(err))
	}

	log.Info("App stopped")
}
