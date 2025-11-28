package model

import (
	"context"
	"github.com/clava1096/rocket-service/order/internal/model"
	converter "github.com/clava1096/rocket-service/order/internal/repository/conveter"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	order, ok := r.orders[uuid]
	if !ok {
		return model.Order{}, model.ErrOrderNotFound
	}

	return converter.OrderFromRepoModel(*order), nil
}
