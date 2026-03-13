// internal/service/order/create_test.go
package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	"github.com/clava1096/rocket-service/order/internal/model"
)

func (s *ServiceSuite) TestCreate_Always() {
	inputOrder := model.Order{
		UUID:       gofakeit.UUID(),
		UserUUID:   gofakeit.UUID(),
		PartUUIDs:  []string{gofakeit.UUID()},
		TotalPrice: 100.0,
		Status:     model.OrderStatusPendingPayment,
	}

	s.inventoryClient.On("ListParts", s.ctx, mock.Anything).
		Return([]model.Part{{Uuid: inputOrder.PartUUIDs[0]}}, nil)

	s.orderRepository.On("Get", s.ctx, inputOrder.UUID).
		Return(model.Order{}, model.ErrOrderNotFound)

	s.orderRepository.On("Create", s.ctx, mock.AnythingOfType("model.Order")).
		Return(inputOrder, nil)

	result, err := s.service.Create(s.ctx, inputOrder)

	s.NoError(err)
	s.Equal(inputOrder.UUID, result.UUID)
}
