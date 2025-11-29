package order

import (
	"context"
	"errors"
	"github.com/clava1096/rocket-service/order/internal/model"
	"time"
)

func (s *service) Pay(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (model.Order, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return model.Order{}, model.ErrOrderNotFound
		}
		return model.Order{}, err
	}
	if order.Status != model.OrderStatusPendingPayment {
		return model.Order{}, model.ErrOrderNotPending
	}
	transactionUUID, err := s.payment.PayOrder(ctx, orderUUID, order.UserUUID, paymentMethod)
	if err != nil {
		return model.Order{}, err
	}

	order.Status = model.OrderStatusPaid
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &paymentMethod
	order.UpdatedAt = time.Now()

	return order, nil
}
