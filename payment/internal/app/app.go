package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentAPI "github.com/clava1096/rocket-service/payment/internal/api/payment/v1"
	"github.com/clava1096/rocket-service/payment/internal/config"
	logger "github.com/clava1096/rocket-service/payment/internal/middleware"
	"github.com/clava1096/rocket-service/platform/pkg/closer"
	platformLogger "github.com/clava1096/rocket-service/platform/pkg/logger"
	paymentv1 "github.com/clava1096/rocket-service/shared/pkg/proto/payment/v1"
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
	a.diContainer = NewDIContainer()
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
	listener, err := net.Listen("tcp", config.AppConfig().Payment.Address())
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
	service := a.diContainer.PaymentService(ctx)
	api := paymentAPI.NewAPI(service)

	a.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(logger.LoggerInterceptor()),
	)

	paymentv1.RegisterPaymentServiceServer(a.grpcServer, api)
	reflection.Register(a.grpcServer)

	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	return nil
}

func (a *App) runGrpcServer(ctx context.Context) error {
	platformLogger.Info(ctx, fmt.Sprintf("gRPC InventoryService server listening on %s", config.AppConfig().Payment.Address()))
	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}
