package model

import "github.com/clava1096/rocket-service/order/internal/model"

func (s *RepositorySuite) TestDelete_Success() {
	order := s.newOrder()
	createdOrder, err := s.repo.Create(s.ctx, order)
	s.NoError(err)

	err = s.repo.Delete(s.ctx, createdOrder.UUID)
	s.NoError(err)

	updatedOrder, err := s.repo.Get(s.ctx, createdOrder.UUID)
	s.NoError(err)
	s.Equal(model.OrderStatusCancelled, updatedOrder.Status)
}

func (s *RepositorySuite) TestDelete_OrderNotFound() {
	err := s.repo.Delete(s.ctx, "non-existent-uuid")

	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)
}
