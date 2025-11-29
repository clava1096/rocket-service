package model

import "github.com/clava1096/rocket-service/order/internal/model"

func (s *RepositorySuite) TestCreate_Success() {
	order := s.newOrder()

	createdOrder, err := s.repo.Create(s.ctx, order)

	s.NoError(err)
	s.Equal(order.UUID, createdOrder.UUID)
	s.Equal(order.UserUUID, createdOrder.UserUUID)
	s.Equal(order.TotalPrice, createdOrder.TotalPrice)
	s.Equal(order.Status, createdOrder.Status)

	retrievedOrder, err := s.repo.Get(s.ctx, order.UUID)
	s.NoError(err)
	s.Equal(createdOrder, retrievedOrder)
}

func (s *RepositorySuite) TestCreate_DuplicateUUID() {
	order := s.newOrder()

	_, err := s.repo.Create(s.ctx, order)
	s.NoError(err)

	_, err = s.repo.Create(s.ctx, order)
	s.Error(err)
	s.ErrorIs(err, model.ErrThisOrderExists)
}
