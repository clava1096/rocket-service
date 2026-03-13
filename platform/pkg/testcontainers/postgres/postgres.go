package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
)

const (
	postgresPort           = "5432"
	postgresStartupTimeout = 5 * time.Minute

	postgresEnvUsernameKey = "POSTGRES_USER"
	postgresEnvPasswordKey = "POSTGRES_PASSWORD"
)

type Container struct {
	container testcontainers.Container
	client    *pgx.Conn
	cfg       *Config
}

func NewContainer(ctx context.Context, opts ...Option) (*Container, error) {
	cfg := buildConfig(opts...)

	container, err := startPostgresContainer(ctx, cfg)
	if err != nil {
		return nil, err
	}

	success := false

	defer func() {
		if !success {
			if err = container.Terminate(ctx); err != nil {
				fmt.Errorf("failed to terminate postgres container: %w", err)
			}
		}
	}()

	cfg.Host, cfg.Port, err = getContainerHostPort(ctx, container)

	if err != nil {
		return nil, err
	}

	uri := buildPostgresURI(cfg)

	client, err := connectPostgresClient(ctx, uri)
	if err != nil {
		return nil, err
	}

	//cfg.Logger.Info(ctx, "postgres container started", zap.String("uri", uri))
	fmt.Println("postgres container started")
	fmt.Print("uri:")
	fmt.Println(uri)

	success = true

	return &Container{
		container: container,
		client:    client,
		cfg:       cfg,
	}, nil
}

func (c *Container) Client() *pgx.Conn {
	return c.client
}

func (c *Container) Config() *Config {
	return c.cfg
}

func (c *Container) Terminate(ctx context.Context) error {
	if c.client != nil {
		if err := c.client.Close(ctx); err != nil {
			fmt.Errorf("failed to close postgres client: %w", err)
		}
	}

	if c.container != nil {
		if err := c.container.Terminate(ctx); err != nil {
			fmt.Errorf("failed to terminate postgres container: %w", err)
		}
	}

	fmt.Println("Postgres container terminated")

	return nil
}

func (c *Container) URI() string {
	return buildPostgresURI(c.cfg)
}
