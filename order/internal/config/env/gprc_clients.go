package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type grpcClientEnvConfig struct {
	InventoryGrpcHost string `env:"INVENTORY_GRPC_HOST"`
	InventoryGrpcPort string `env:"INVENTORY_GRPC_PORT"`
	PaymentGrpcHost   string `env:"PAYMENT_GRPC_HOST"`
	PaymentGrpcPort   string `env:"PAYMENT_GRPC_PORT"`
}

type grpcConfig struct {
	client grpcClientEnvConfig
}

func NewGrpcConfig() (*grpcConfig, error) {
	var client grpcClientEnvConfig

	err := env.Parse(&client)
	if err != nil {
		return nil, err
	}
	return &grpcConfig{client: client}, nil
}

func (p *grpcConfig) InventoryURI() string {
	return net.JoinHostPort(p.client.InventoryGrpcHost, p.client.InventoryGrpcPort)
}

func (p *grpcConfig) PaymentURI() string {
	return net.JoinHostPort(p.client.PaymentGrpcHost, p.client.PaymentGrpcPort)
}
