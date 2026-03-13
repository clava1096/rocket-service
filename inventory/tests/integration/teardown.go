package integration

import (
	"context"

	"github.com/clava1096/rocket-service/platform/pkg/logger"
	"go.uber.org/zap"
)

func teardownTestEnvironment(ctx context.Context, env *TestEnvironment) {
	log := logger.Logger()
	log.Info(ctx, "Cleanup test environment...")

	cleanupTestEnvironment(ctx, env)

	log.Info(ctx, "Test environment successful clean")
}

func cleanupTestEnvironment(ctx context.Context, env *TestEnvironment) {
	if env.App != nil {
		if err := env.App.Terminate(ctx); err != nil {
			logger.Error(ctx, "Failed to terminate app container", zap.Error(err))
		} else {
			logger.Info(ctx, "Container application terminated")
		}
	}

	if env.Mongo != nil {
		if err := env.Mongo.Terminate(ctx); err != nil {
			logger.Error(ctx, "failed to terminate mongo container", zap.Error(err))
		} else {
			logger.Info(ctx, "Mongo container terminated")
		}
	}

	if env.Network != nil {
		if err := env.Network.Remove(ctx); err != nil {
			logger.Error(ctx, "Failed to remove network container", zap.Error(err))
		} else {
			logger.Info(ctx, "Network removed")
		}
	}
}
