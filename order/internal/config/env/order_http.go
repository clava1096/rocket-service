package env

import (
	"log"
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type serverEnvConfig struct {
	Host    string        `env:"HTTP_HOST"`
	Port    string        `env:"HTTP_PORT"`
	Timeout time.Duration `env:"HTTP_TIMEOUT"`
}

type serverConfig struct {
	raw serverEnvConfig
}

func NewServerEnvConfig() (*serverConfig, error) {
	var conf serverEnvConfig

	err := env.Parse(&conf)

	log.Print(conf)

	if err != nil {
		return nil, err
	}

	return &serverConfig{conf}, nil
}

func (s *serverConfig) Address() string {
	return net.JoinHostPort(s.raw.Host, s.raw.Port)
}

func (s *serverConfig) HeaderTimeout() time.Duration {
	return s.raw.Timeout
}
