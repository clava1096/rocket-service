package service

import (
	"context"
	"github.com/clava1096/rocket-service/order/internal/model"
)

type OrderService interface {
	Cancel(ctx context.Context, uuid string) error

	Create(ctx context.Context, order model.Order) (model.Order, error)

	Get(ctx context.Context, uuid string) (model.Order, error)

	Pay(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (model.Order, error)
}
