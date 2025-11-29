package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	grpc "github.com/clava1096/rocket-service/order/internal/client/grpc/mocks"
	rep "github.com/clava1096/rocket-service/order/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	orderRepository *rep.OrderRepository

	inventoryClient *grpc.InventoryClient

	payClient *grpc.PaymentClient

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.orderRepository = rep.NewOrderRepository(s.T())

	s.inventoryClient = grpc.NewInventoryClient(s.T())

	s.payClient = grpc.NewPaymentClient(s.T())

	s.service = NewService(s.orderRepository, s.inventoryClient, s.payClient)
}
func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
