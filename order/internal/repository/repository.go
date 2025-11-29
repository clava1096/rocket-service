package repository

import (
	"context"

	"github.com/clava1096/rocket-service/order/internal/model"
)

type OrderRepository interface {
	Get(ctx context.Context, uuid string) (model.Order, error)
	Delete(ctx context.Context, uuid string) error
	Create(ctx context.Context, order model.Order) (model.Order, error)
	Update(ctx context.Context, order model.Order) (model.Order, error)
}
