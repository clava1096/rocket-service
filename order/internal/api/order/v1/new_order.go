package v1

import (
	"context"
	"net/http"

	orderV1 "github.com/clava1096/rocket-service/shared/pkg/openapi/order/v1"
)

func (a *api) NewError(ctx context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: 500,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}
