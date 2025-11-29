package model

import "github.com/clava1096/rocket-service/order/internal/model"

func (s *RepositorySuite) TestGet_Success() {
	order := s.newOrder()
	createdOrder, err := s.repo.Create(s.ctx, order)
	s.NoError(err)

	retrievedOrder, err := s.repo.Get(s.ctx, createdOrder.UUID)

	s.NoError(err)
	s.Equal(createdOrder, retrievedOrder)
}

func (s *RepositorySuite) TestGet_OrderNotFound() {
	_, err := s.repo.Get(s.ctx, "non-existent-uuid")

	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)
}
