package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/clava1096/rocket-service/iam/internal/config"
	grpcMiddleware "github.com/clava1096/rocket-service/iam/internal/middleware"
	"github.com/clava1096/rocket-service/platform/pkg/closer"
	platformLogger "github.com/clava1096/rocket-service/platform/pkg/logger"
	iamv1 "github.com/clava1096/rocket-service/shared/pkg/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runGrpcServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initGrpcServer,
	}
	for _, init := range inits {
		if err := init(ctx); err != nil {
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
	return platformLogger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson())
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(platformLogger.Logger())
	return nil
}

func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().Server.Address())
	if err != nil {
		return err
	}
	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		lerr := listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}
		return nil
	})
	a.listener = listener
	return nil
}

func (a *App) initGrpcServer(ctx context.Context) error {
	IamV1API := a.diContainer.IamV1API(ctx)
	a.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpcMiddleware.LoggerInterceptor()))

	iamv1.RegisterUserServiceServer(a.grpcServer, IamV1API)
	reflection.Register(a.grpcServer)

	closer.AddNamed("GRPC listener", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	return nil
}

func (a *App) runGrpcServer(ctx context.Context) error {
	platformLogger.Info(ctx, fmt.Sprintf("gRPC InventoryService server listening on %s", config.AppConfig().Server.Address()))
	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}
