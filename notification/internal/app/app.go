package app

import (
	"context"

	"github.com/clava1096/rocket-service/notification/internal/config"
	"github.com/clava1096/rocket-service/platform/pkg/closer"
	"github.com/clava1096/rocket-service/platform/pkg/logger"
	"github.com/go-faster/errors"
	"go.uber.org/zap"
)

type App struct {
	diContainer *diContainer
}

func New(ctx context.Context) (*App, error) {
	a := &App{}
	err := a.InitDeps(ctx)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) InitDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 2)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := a.runPaidConsumer(ctx); err != nil {
			errCh <- errors.Errorf("paid consumer crashed: %v", err)
		}
	}()

	go func() {
		if err := a.runAssemblyConsumer(ctx); err != nil {
			errCh <- errors.Errorf("assembly consumer crashed: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "shutting down signal received")
	case err := <-errCh:
		logger.Info(ctx, "component crashed shutting down", zap.Error(err))

		cancel()

		<-ctx.Done()
		return err
	}

	return nil
}

func (a *App) runPaidConsumer(ctx context.Context) error {
	logger.Info(ctx, "paid consumer running...")
	err := a.diContainer.OrderPaidConsumer().RunConsumer(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (a *App) runAssemblyConsumer(ctx context.Context) error {
	logger.Info(ctx, "assembly consumer running...")
	err := a.diContainer.OrderAssembledConsumer().RunConsumer(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDIContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().LoggerConfig.Level(),
		config.AppConfig().LoggerConfig.AsJson())
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}
