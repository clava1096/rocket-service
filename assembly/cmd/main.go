package main

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"github.com/clava1096/rocket-service/assembly/internal/app"
	"github.com/clava1096/rocket-service/assembly/internal/config"
	"github.com/clava1096/rocket-service/platform/pkg/closer"
	"github.com/clava1096/rocket-service/platform/pkg/logger"
	"go.uber.org/zap"
)

const configPath = "./deploy/compose/assembly/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "failed to initialize application", zap.Error(err))
		return
	}

	err = a.Run(appCtx)

	if err != nil {
		logger.Error(appCtx, "failed to run application", zap.Error(err))
		return
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "graceful shutdown failed", zap.Error(err))
	}
}
