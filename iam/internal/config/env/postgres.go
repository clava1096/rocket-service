package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type postgresEnv struct {
	Host          string `env:"POSTGRES_HOST"`
	Port          string `env:"EXTERNAL_POSTGRES_PORT"`
	User          string `env:"POSTGRES_USER"`
	Password      string `env:"POSTGRES_PASSWORD"`
	Database      string `env:"POSTGRES_DB"`
	MigrationsDir string `env:"MIGRATION_DIRECTORY"`
}

type postgres struct {
	raw postgresEnv
}

func NewPostgresConfig() (*postgres, error) {
	var raw postgresEnv

	err := env.Parse(&raw)

	if err != nil {
		return nil, err
	}
	return &postgres{raw: raw}, nil
}

func (p *postgres) URI() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		p.raw.User,
		p.raw.Password,
		p.raw.Host,
		p.raw.Port,
		p.raw.Database,
	)
}

func (p *postgres) MigrationsDir() string {
	return p.raw.MigrationsDir
}
