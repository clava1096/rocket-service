package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/clava1096/rocket-service/order/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/clava1096/rocket-service/order/internal/converter"
	orderV1 "github.com/clava1096/rocket-service/shared/pkg/openapi/order/v1"
)

func (a *api) OrderPayment(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.OrderPaymentParams) (orderV1.OrderPaymentRes, error) {

	paymentMethod := converter.PaymentMethodFromOpenAPI(req.PaymentMethod)

	// 2. Вызываем сервис
	updatedOrder, err := a.orderService.Pay(ctx, params.OrderUUID, paymentMethod)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order '%s' not found", params.OrderUUID),
			}, nil
		}

		return nil, status.Errorf(codes.Internal, "payment failed: %v", err)
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: converter.MustParseUUID(*updatedOrder.TransactionUUID),
	}, nil

}
