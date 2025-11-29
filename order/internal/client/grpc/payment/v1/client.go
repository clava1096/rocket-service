package v1

import (
	def "github.com/clava1096/rocket-service/order/internal/client/grpc"
	paymentv1 "github.com/clava1096/rocket-service/shared/pkg/proto/payment/v1"
)

var _ def.PaymentClient = (*client)(nil)

type client struct {
	paymentClient paymentv1.PaymentServiceClient
}

func NewClient(paymentClient paymentv1.PaymentServiceClient) *client {
	return &client{
		paymentClient: paymentClient,
	}
}
