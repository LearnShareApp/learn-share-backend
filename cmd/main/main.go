package main

import (
	"context"
	"github.com/LearnShareApp/learn-share-backend/internal/application"
	"github.com/LearnShareApp/learn-share-backend/internal/config"
	"github.com/LearnShareApp/learn-share-backend/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	ctx, ctxClose := context.WithCancel(context.Background())
	defer ctxClose()

	// development logger (красивый вывод в консоль)
	log := logger.NewDevelopment()
	defer log.Sync()

	// load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("failed to load config", zap.Error(err))
		panic("failed to load config")
	}

	log.Info(cfg.LogConfig())

	app, err := application.New(ctx, *cfg, log)

	if err != nil {
		log.Error("failed to create application", zap.Error(err))
	}

	app.Run()
}
