package v1

import (
	"context"
	"errors"

	"github.com/clava1096/rocket-service/order/internal/converter"
	"github.com/clava1096/rocket-service/order/internal/model"
	orderV1 "github.com/clava1096/rocket-service/shared/pkg/openapi/order/v1"
)

func (a *api) GetInfoOrderByUUID(ctx context.Context, params orderV1.GetInfoOrderByUUIDParams) (orderV1.GetInfoOrderByUUIDRes, error) {
	order, err := a.orderService.Get(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "Order not found",
			}, nil
		}
	}

	return converter.OrderToOpenAPI(order), nil
}
