package app

import (
	"context"

	"github.com/clava1096/rocket-service/assembly/internal/config"
	"github.com/clava1096/rocket-service/platform/pkg/closer"
	"github.com/clava1096/rocket-service/platform/pkg/logger"
	"github.com/go-faster/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
}

func New(ctx context.Context) (*App, error) {
	a := &App{}
	err := a.initDeps(ctx)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errCh <- errors.Errorf("consumer crashed: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "shutting down gracefully")
	case err := <-errCh:
		logger.Error(ctx, "Component crashed", zap.Error(err))
		cancel()
		<-ctx.Done()
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
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

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "starting order paid consumer...")

	err := a.diContainer.OrderConsumerService().RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}
