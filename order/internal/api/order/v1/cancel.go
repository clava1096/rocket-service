package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/clava1096/rocket-service/order/internal/model"
	orderV1 "github.com/clava1096/rocket-service/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	err := a.orderService.Cancel(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{Message: "Order '" + params.OrderUUID + "' not found"}, nil
		}
		if errors.Is(err, model.ErrOrderNotPending) {
			return &orderV1.ConflictError{Message: "Order '" + params.OrderUUID + "' was paid"}, nil
		}
		return nil, status.Errorf(codes.Internal, "failed to cancel order: %v", err)
	}

	return &orderV1.CancelOrderNoContent{}, nil
}
