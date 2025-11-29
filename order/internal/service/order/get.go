package order

import (
	"context"
	"errors"

	"github.com/clava1096/rocket-service/order/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (model.Order, error) {
	order, err := s.orderRepository.Get(ctx, uuid)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return model.Order{}, model.ErrOrderNotFound
		}
		return model.Order{}, err
	}
	return order, nil
}
