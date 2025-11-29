package order

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	"github.com/clava1096/rocket-service/order/internal/model"
)

func (s *ServiceSuite) TestCancel_Success() {
	uuid := gofakeit.UUID()
	order := model.Order{
		UUID:      uuid,
		UserUUID:  gofakeit.UUID(),
		Status:    model.OrderStatusPendingPayment,
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	s.orderRepository.On("Get", s.ctx, uuid).Return(order, nil)

	updatedOrder := order
	updatedOrder.Status = model.OrderStatusCancelled
	s.orderRepository.On("Update", s.ctx, mock.MatchedBy(func(o model.Order) bool {
		return o.UUID == uuid &&
			o.Status == model.OrderStatusCancelled &&
			!o.UpdatedAt.IsZero()
	})).Return(updatedOrder, nil)

	err := s.service.Cancel(s.ctx, uuid)

	s.NoError(err)

	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCancel_OrderNotFound() {
	uuid := gofakeit.UUID()

	s.orderRepository.On("Get", s.ctx, uuid).Return(model.Order{}, model.ErrOrderNotFound)

	err := s.service.Cancel(s.ctx, uuid)

	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)

	s.orderRepository.AssertNotCalled(s.T(), "Update")
}

func (s *ServiceSuite) TestCancel_OrderAlreadyPaid() {
	uuid := gofakeit.UUID()
	order := model.Order{
		UUID:   uuid,
		Status: model.OrderStatusPaid,
	}

	s.orderRepository.On("Get", s.ctx, uuid).Return(order, nil)

	err := s.service.Cancel(s.ctx, uuid)

	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotPending)

	s.orderRepository.AssertNotCalled(s.T(), "Update")
}

func (s *ServiceSuite) TestCancel_RepositoryUpdateError() {
	uuid := gofakeit.UUID()
	order := model.Order{
		UUID:   uuid,
		Status: model.OrderStatusPendingPayment,
	}

	s.orderRepository.On("Get", s.ctx, uuid).Return(order, nil)

	updateError := gofakeit.Error()
	s.orderRepository.On("Update", s.ctx, mock.AnythingOfType("model.Order")).Return(model.Order{}, updateError)

	err := s.service.Cancel(s.ctx, uuid)

	s.Error(err)
	s.ErrorIs(err, updateError)

	s.orderRepository.AssertExpectations(s.T())
}
