package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/clava1096/rocket-service/order/internal/model"
)

func (s *service) Create(ctx context.Context, order model.Order) (model.Order, error) {
	_, err := s.orderRepository.Get(ctx, order.UUID)
	if errors.Is(err, model.ErrOrderNotFound) {
		return model.Order{}, model.ErrOrderNotFound
	}

	parts, err := s.inventory.ListParts(ctx, model.PartsFilter{
		Uuids: order.PartUUIDs,
	})
	if err != nil {
		return model.Order{}, model.ErrPartNotFound
	}

	found := make(map[string]bool, len(parts))
	var totalPrice float64

	for _, part := range parts {
		found[part.Uuid] = true
		totalPrice += part.Price
	}

	for _, uuid := range order.PartUUIDs {
		if !found[uuid] {
			return model.Order{}, fmt.Errorf("%w: part %s not found", model.ErrPartNotFound, uuid)
		}
	}

	orderToCreate := model.Order{
		UUID:       order.UUID,
		UserUUID:   order.UserUUID,
		PartUUIDs:  order.PartUUIDs,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPendingPayment,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	createdOrder, err := s.orderRepository.Create(ctx, orderToCreate)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to create order: %w", err)
	}

	return createdOrder, nil
}
