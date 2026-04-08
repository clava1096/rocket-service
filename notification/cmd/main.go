package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/clava1096/rocket-service/notification/internal/app"
	"github.com/clava1096/rocket-service/notification/internal/config"
	"github.com/clava1096/rocket-service/platform/pkg/closer"
	"github.com/clava1096/rocket-service/platform/pkg/logger"
	"go.uber.org/zap"
)

const configPath = "./deploy/compose/notifiction/.env" //todo добавить dockercompose и .env

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Sprintf("Error loading config: %v", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "failed to run app", zap.Error(err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, "failed to run app", zap.Error(err))
		return
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "failed to gracefully shutdown", zap.Error(err))
	}
}
