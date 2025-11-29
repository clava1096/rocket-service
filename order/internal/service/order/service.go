package order

import (
	"github.com/clava1096/rocket-service/order/internal/client/grpc"
	"github.com/clava1096/rocket-service/order/internal/repository"
	def "github.com/clava1096/rocket-service/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepository repository.OrderRepository

	inventory grpc.InventoryClient
	payment   grpc.PaymentClient
}

func NewService(orderRepository repository.OrderRepository, inventory grpc.InventoryClient, payment grpc.PaymentClient) *service {
	return &service{
		orderRepository,
		inventory,
		payment,
	}
}
