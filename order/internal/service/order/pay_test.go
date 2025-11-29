package order

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/clava1096/rocket-service/order/internal/model"
)

func (s *ServiceSuite) TestPay_Success() {
	uuid := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	paymentMethod := model.PaymentMethodCard
	transactionUUID := gofakeit.UUID()

	order := model.Order{
		UUID:      uuid,
		UserUUID:  userUUID,
		Status:    model.OrderStatusPendingPayment,
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	s.orderRepository.On("Get", s.ctx, uuid).Return(order, nil)

	s.payClient.On("PayOrder", s.ctx, uuid, userUUID, paymentMethod).
		Return(transactionUUID, nil)

	updatedOrder, err := s.service.Pay(s.ctx, uuid, paymentMethod)

	s.NoError(err)
	s.Equal(uuid, updatedOrder.UUID)
	s.Equal(model.OrderStatusPaid, updatedOrder.Status)
	s.Equal(&transactionUUID, updatedOrder.TransactionUUID)
	s.Equal(&paymentMethod, updatedOrder.PaymentMethod)
	s.True(updatedOrder.UpdatedAt.After(order.UpdatedAt))

	s.orderRepository.AssertExpectations(s.T())
	s.payClient.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestPay_OrderNotFound() {
	uuid := gofakeit.UUID()
	paymentMethod := model.PaymentMethodCard

	s.orderRepository.On("Get", s.ctx, uuid).
		Return(model.Order{}, model.ErrOrderNotFound)

	_, err := s.service.Pay(s.ctx, uuid, paymentMethod)

	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)

	s.payClient.AssertNotCalled(s.T(), "PayOrder")
}

func (s *ServiceSuite) TestPay_OrderNotPending() {
	uuid := gofakeit.UUID()
	paymentMethod := model.PaymentMethodCard

	// Заказ в статусе PAID
	order := model.Order{
		UUID:   uuid,
		Status: model.OrderStatusPaid,
	}

	s.orderRepository.On("Get", s.ctx, uuid).Return(order, nil)

	_, err := s.service.Pay(s.ctx, uuid, paymentMethod)

	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotPending)

	s.payClient.AssertNotCalled(s.T(), "PayOrder")
}

func (s *ServiceSuite) TestPay_PaymentServiceError() {
	uuid := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	paymentMethod := model.PaymentMethodCard

	order := model.Order{
		UUID:     uuid,
		UserUUID: userUUID,
		Status:   model.OrderStatusPendingPayment,
	}

	s.orderRepository.On("Get", s.ctx, uuid).Return(order, nil)

	paymentError := gofakeit.Error()
	s.payClient.On("PayOrder", s.ctx, uuid, userUUID, paymentMethod).
		Return("", paymentError)

	_, err := s.service.Pay(s.ctx, uuid, paymentMethod)

	s.Error(err)
	s.ErrorIs(err, paymentError)

	s.orderRepository.AssertNotCalled(s.T(), "Update")
}
