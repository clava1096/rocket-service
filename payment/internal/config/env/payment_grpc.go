package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type paymentEnvConfig struct {
	Port string `env:"GRPC_PORT"`
	Host string `env:"GRPC_HOST"`
}

type paymentConfig struct {
	raw paymentEnvConfig
}

func NewPaymentEnvConfig() (*paymentConfig, error) {
	var config paymentEnvConfig
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &paymentConfig{raw: config}, nil
}

func (config *paymentConfig) Address() string {
	return net.JoinHostPort(config.raw.Host, config.raw.Port)
}
