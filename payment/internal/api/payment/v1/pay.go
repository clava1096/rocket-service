package v1

import (
	"context"
	"log"

	paymentv1 "github.com/clava1096/rocket-service/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, payOrder *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	transactionUuid := a.paymentService.Pay(ctx)
	log.Printf("Payment processed: txn=%s, order=%s, user=%s, method=%s",
		transactionUuid,
		payOrder.OrderUuid,
		payOrder.UserUuid,
		payOrder.PaymentMethod.String())
	return &paymentv1.PayOrderResponse{TransactionUuid: transactionUuid}, nil
}
