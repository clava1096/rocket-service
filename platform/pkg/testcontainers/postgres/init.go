package postgres

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-faster/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func startPostgresContainer(ctx context.Context, cfg *Config) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Name:     cfg.ContainerName,
		Image:    cfg.ImageName,
		Networks: []string{cfg.NetworkName},
		Env: map[string]string{
			postgresEnvUsernameKey: cfg.Username,
			postgresEnvPasswordKey: cfg.Password,
			"POSTGRES_DB":          cfg.Database,
		},
		WaitingFor:         wait.ForListeningPort(postgresPort + "/tcp").WithStartupTimeout(postgresStartupTimeout),
		HostConfigModifier: defaultHostConfig(),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, errors.Errorf("failed to start postgres container: %v", err)
	}
	return container, nil
}

func getContainerHostPort(ctx context.Context, container testcontainers.Container) (string, string, error) {
	host, err := container.Host(ctx)

	if err != nil {
		return "", "", errors.Errorf("failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, postgresPort+"/tcp") // todo может не работать из-за /tcp т.к по идеи нету функции

	return host, port.Port(), err
}

func buildPostgresURI(cfg *Config) string {
	username := url.QueryEscape(cfg.Username)
	password := url.QueryEscape(cfg.Password)

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username, password, cfg.Host, cfg.Port, cfg.Database)
}
