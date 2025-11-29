package v1

import (
	def "github.com/clava1096/rocket-service/order/internal/client/grpc"
	generatedInventoryV1 "github.com/clava1096/rocket-service/shared/pkg/proto/inventory/v1"
)

var _ def.InventoryClient = (*client)(nil)

type client struct {
	generatedClient generatedInventoryV1.InventoryServiceClient
}

func NewClient(generated generatedInventoryV1.InventoryServiceClient) *client {
	return &client{generatedClient: generated}
}
