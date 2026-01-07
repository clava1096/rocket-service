package model

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/clava1096/rocket-service/order/internal/model"
)

func (s *RepositorySuite) TestUpdate_Success() {
	order := s.newOrder()
	createdOrder, err := s.repo.Create(s.ctx, order)
	s.NoError(err)

	updatedOrderInput := createdOrder
	updatedOrderInput.Status = model.OrderStatusPaid
	updatedOrderInput.TotalPrice = 999.99

	updatedOrder, err := s.repo.Update(s.ctx, updatedOrderInput)
	s.NoError(err)

	s.Equal(model.OrderStatusPaid, updatedOrder.Status)
	s.Equal(999.99, updatedOrder.TotalPrice)

	retrievedOrder, err := s.repo.Get(s.ctx, updatedOrder.UUID)
	s.NoError(err)
	s.Equal(model.OrderStatusPaid, retrievedOrder.Status)
	s.Equal(999.99, retrievedOrder.TotalPrice)
}

func (s *RepositorySuite) TestUpdate_OrderNotFound() {
	order := s.newOrder()
	order.UUID = gofakeit.UUID()

	_, err := s.repo.Update(s.ctx, order)

	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)
}
