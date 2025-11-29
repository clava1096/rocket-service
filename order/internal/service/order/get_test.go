package order

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/clava1096/rocket-service/order/internal/model"
)

func (s *ServiceSuite) TestGet_Success() {
	expectedOrder := model.Order{
		UUID:       gofakeit.UUID(),
		UserUUID:   gofakeit.UUID(),
		PartUUIDs:  []string{gofakeit.UUID(), gofakeit.UUID()},
		TotalPrice: gofakeit.Float64Range(10, 10000),
		Status:     model.OrderStatusPendingPayment,
		CreatedAt:  time.Now().Truncate(time.Second),
		UpdatedAt:  time.Now().Truncate(time.Second),
	}

	s.orderRepository.On("Get", s.ctx, expectedOrder.UUID).Return(expectedOrder, nil)

	actualOrder, err := s.service.Get(s.ctx, expectedOrder.UUID)

	s.NoError(err)
	s.Equal(expectedOrder.UUID, actualOrder.UUID)
	s.Equal(expectedOrder.UserUUID, actualOrder.UserUUID)
	s.Equal(expectedOrder.PartUUIDs, actualOrder.PartUUIDs)
	s.Equal(expectedOrder.TotalPrice, actualOrder.TotalPrice)
	s.Equal(expectedOrder.Status, actualOrder.Status)
	s.Equal(expectedOrder.CreatedAt, actualOrder.CreatedAt)
	s.Equal(expectedOrder.UpdatedAt, actualOrder.UpdatedAt)

	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGet_OrderNotFound() {
	uuid := gofakeit.UUID()

	s.orderRepository.On("Get", s.ctx, uuid).Return(model.Order{}, model.ErrOrderNotFound)

	order, err := s.service.Get(s.ctx, uuid)

	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)
	s.Empty(order)

	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGet_RepositoryInternalError() {
	uuid := gofakeit.UUID()

	internalError := gofakeit.Error()
	s.orderRepository.On("Get", s.ctx, uuid).Return(model.Order{}, internalError)

	order, err := s.service.Get(s.ctx, uuid)

	s.Error(err)
	s.ErrorIs(err, internalError)
	s.Empty(order)

	s.orderRepository.AssertExpectations(s.T())
}
