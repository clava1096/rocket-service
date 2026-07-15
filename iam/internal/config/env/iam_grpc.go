package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type iamEnvGrpc struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type iamGrpc struct {
	raw iamEnvGrpc
}

func NewIamEnvGrpc() (*iamGrpc, error) {
	var raw iamEnvGrpc

	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}
	return &iamGrpc{raw: raw}, nil
}

func (iag *iamGrpc) Address() string {
	return net.JoinHostPort(iag.raw.Host, iag.raw.Port)
}
