package model

import (
	"context"
	"github.com/clava1096/rocket-service/order/internal/model"
	repo "github.com/clava1096/rocket-service/order/internal/repository/model"
)

func (r *repository) Delete(ctx context.Context, uuid string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	repoOrder, ok := r.orders[uuid]
	if !ok {
		return model.ErrOrderNotFound
	}
	repoOrder.Status = repo.OrderStatus(model.OrderStatusCancelled)

	return nil
}
