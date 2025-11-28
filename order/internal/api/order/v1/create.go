package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/clava1096/rocket-service/order/internal/converter"
	"github.com/clava1096/rocket-service/order/internal/model"
	orderV1 "github.com/clava1096/rocket-service/shared/pkg/openapi/order/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) CreateOrderRequest(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRequestRes, error) {
	if req.UUID == "" {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "order UUID cannot be empty",
		}, nil
	}

	if len(req.PartsUUID) == 0 {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "parts UUID list cannot be empty",
		}, nil
	}

	domainOrder := converter.CreateOrderRequestToModel(req)
	order, err := a.orderService.Create(ctx, domainOrder)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return &orderV1.BadRequestError{
				Code:    400,
				Message: fmt.Sprintf("one or more parts not found: %v", err),
			}, nil
		}
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}
	return &orderV1.CreateOrderResponse{
		OrderUUID:  converter.MustParseUUID(order.UUID),
		TotalPrice: order.TotalPrice,
	}, nil
}
