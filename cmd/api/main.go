package main

import (
	"application-service/internal/app"
	"application-service/internal/config"
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("failed to init logger:", err)
	}
	defer func() { _ = logger.Sync() }()

	application, err := app.New(context.Background(), cfg, logger)
	if err != nil {
		logger.Fatal("app init failed", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := application.Run(ctx); err != nil {
		logger.Fatal("app run failed", zap.Error(err))
	}

	// extra wait for logs flush in some environments
	time.Sleep(100 * time.Millisecond)
}
