// internal/service/order/create_test.go
package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/clava1096/rocket-service/order/internal/model"
)

func (s *ServiceSuite) TestCreate_AlwaysFailsWithOrderNotFound() {
	inputOrder := model.Order{
		UUID:      gofakeit.UUID(),
		UserUUID:  gofakeit.UUID(),
		PartUUIDs: []string{gofakeit.UUID()},
	}

	s.orderRepository.On("Get", s.ctx, inputOrder.UUID).
		Return(model.Order{}, model.ErrOrderNotFound)

	_, err := s.service.Create(s.ctx, inputOrder)

	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)

	s.inventoryClient.AssertNotCalled(s.T(), "ListParts")
	s.orderRepository.AssertNotCalled(s.T(), "Create")
}
