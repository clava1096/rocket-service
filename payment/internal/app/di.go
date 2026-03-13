package app

import (
	"context"

	paymentAPI "github.com/clava1096/rocket-service/payment/internal/api/payment/v1"
	"github.com/clava1096/rocket-service/payment/internal/service"
	paymentService "github.com/clava1096/rocket-service/payment/internal/service/payment"
	paymentv1 "github.com/clava1096/rocket-service/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1      paymentv1.PaymentServiceServer
	paymentService service.PaymentService
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1(ctx context.Context) paymentv1.PaymentServiceServer {
	if d.paymentV1 == nil {
		d.paymentV1 = paymentAPI.NewAPI(d.PaymentService(ctx))
	}
	return d.paymentV1
}

func (d *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewService()
	}
	return d.paymentService
}
