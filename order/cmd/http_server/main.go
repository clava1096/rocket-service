package main

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"github.com/clava1096/rocket-service/order/internal/app"
	"github.com/clava1096/rocket-service/order/internal/config"
	"github.com/clava1096/rocket-service/platform/pkg/closer"
	"github.com/clava1096/rocket-service/platform/pkg/logger"
	"go.uber.org/zap"
)

const shutdownTimeout = 10 * time.Second

const configPath = "./deploy/compose/order/.env"

func main() {
	err := config.Load(configPath)

	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	runApp()
}

func runApp() {
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.NewApp(appCtx)

	if err != nil {
		logger.Error(appCtx, "cannot init app", zap.Error(err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, "cannot run app", zap.Error(err))
		return
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(context.Background(), "failed to close all graceful shutdown", zap.Error(err))
	}
}
