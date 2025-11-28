package model

import (
	"context"
	"github.com/clava1096/rocket-service/order/internal/model"
	converter "github.com/clava1096/rocket-service/order/internal/repository/conveter"
)

func (r *repository) Update(ctx context.Context, order model.Order) (model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[order.UUID]; !ok {
		return model.Order{}, model.ErrOrderNotFound
	}

	repoOrder := converter.OrderToRepoModel(order)

	r.orders[order.UUID] = &repoOrder

	return converter.OrderFromRepoModel(repoOrder), nil
}
