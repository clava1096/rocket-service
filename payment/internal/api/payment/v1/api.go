package v1

import (
	"github.com/clava1096/rocket-service/payment/internal/service"
	paymentv1 "github.com/clava1096/rocket-service/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentv1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func NewAPI(paymentService service.PaymentService) *api {
	return &api{paymentService: paymentService}
}
