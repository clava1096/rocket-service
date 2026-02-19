package app

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	orderAPI "github.com/clava1096/rocket-service/order/internal/api/order/v1"
	"github.com/clava1096/rocket-service/order/internal/config"
	m "github.com/clava1096/rocket-service/order/internal/middleware"
	"github.com/clava1096/rocket-service/order/internal/migrator"
	"github.com/clava1096/rocket-service/platform/pkg/closer"
	"github.com/clava1096/rocket-service/platform/pkg/logger"
	orderV1 "github.com/clava1096/rocket-service/shared/pkg/openapi/order/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
	listener    net.Listener
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	err := a.initDeps(ctx)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runHttpServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDi,
		a.initLogger,
		a.initCloser,
		a.initMigration,
		a.initListener,
		a.runHttpServer,
	}
	for _, init := range inits {
		if err := init(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDi(_ context.Context) error {
	a.diContainer = NewDIContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson())
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initMigration(_ context.Context) error {
	pgxConfig, err := pgx.ParseConfig(config.AppConfig().Postgres.URI())
	if err != nil {
		return err
	}

	sqlDB := stdlib.OpenDB(*pgxConfig)

	closer.AddNamed("SQL DB migration", func(ctx context.Context) error {
		return sqlDB.Close()
	})

	migratorRunner := migrator.NewMigrator(
		sqlDB,
		config.AppConfig().Postgres.MigrationsDir(),
	)

	if err = migratorRunner.Up(); err != nil {
		return err
	}

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

func (a *App) runHttpServer(ctx context.Context) error {
	service := a.diContainer.OrderService(ctx)
	api := orderAPI.NewAPI(service)

	storageServer, err := orderV1.NewServer(api)
	if err != nil {
		return err
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(m.RequestLogger)

	r.Mount("/", storageServer)

	a.httpServer = &http.Server{
		Handler:           r,
		ReadHeaderTimeout: config.AppConfig().Server.HeaderTimeout(),
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})

	return nil
}
