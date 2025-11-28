package model

import (
	"context"
	"github.com/clava1096/rocket-service/order/internal/model"
	converter "github.com/clava1096/rocket-service/order/internal/repository/conveter"
)

func (r *repository) Create(ctx context.Context, order model.Order) (model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.UUID]; exists {
		return model.Order{}, model.ErrThisOrderExists
	}

	repoOrder := converter.OrderToRepoModel(order)
	r.orders[repoOrder.UUID] = &repoOrder

	return converter.OrderFromRepoModel(repoOrder), nil
}
