package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type inventoryEnvGrpc struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type inventoryGrpc struct {
	raw inventoryEnvGrpc
}

func NewInventoryEnvGrpc() (*inventoryGrpc, error) {
	var raw inventoryEnvGrpc

	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &inventoryGrpc{raw: raw}, nil
}

func (inventory *inventoryGrpc) Address() string {
	return net.JoinHostPort(inventory.raw.Host, inventory.raw.Port)
}
