package order

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	"github.com/clava1096/rocket-service/order/internal/model"
)

func (s *ServiceSuite) TestPay_Success() {
	uuid := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	paymentMethod := model.PaymentMethodCard
	transactionUUID := gofakeit.UUID()

	// Исходный заказ
	order := model.Order{
		UUID:      uuid,
		UserUUID:  userUUID,
		Status:    model.OrderStatusPendingPayment,
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	s.orderRepository.On("Get", s.ctx, uuid).Return(order, nil)

	s.payClient.On("PayOrder", s.ctx, uuid, userUUID, paymentMethod).
		Return(transactionUUID, nil)

	now := time.Now()
	updatedOrder := model.Order{
		UUID:            uuid,
		UserUUID:        userUUID,
		Status:          model.OrderStatusPaid,
		PaymentMethod:   &paymentMethod,
		TransactionUUID: &transactionUUID,
		UpdatedAt:       now,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		CreatedAt:       order.CreatedAt,
	}

	s.orderRepository.On("Update", s.ctx, mock.MatchedBy(func(o model.Order) bool {
		return o.UUID == updatedOrder.UUID &&
			o.Status == updatedOrder.Status &&
			o.TransactionUUID != nil && *o.TransactionUUID == transactionUUID
	})).Return(updatedOrder, nil)

	result, err := s.service.Pay(s.ctx, uuid, paymentMethod)

	s.NoError(err)
	s.Equal(uuid, result.UUID)
	s.Equal(model.OrderStatusPaid, result.Status)
	s.Equal(&transactionUUID, result.TransactionUUID)
	s.Equal(&paymentMethod, result.PaymentMethod)
	s.True(result.UpdatedAt.After(order.UpdatedAt))

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
