package order

import (
	"context"
	"github.com/clava1096/rocket-service/order/internal/model"
	"time"
)

func (s *service) Cancel(ctx context.Context, uuid string) error {
	order, err := s.orderRepository.Get(ctx, uuid)
	if err != nil {
		return err
	}

	if order.Status == model.OrderStatusPaid {
		return model.ErrOrderNotPending
	}

	order.Status = model.OrderStatusCancelled
	order.UpdatedAt = time.Now()

	_, err = s.orderRepository.Update(ctx, order)
	return err
}
