package v1

import (
	"context"

	"github.com/clava1096/rocket-service/order/internal/model"
	paymentv1 "github.com/clava1096/rocket-service/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod model.PaymentMethod) (string, error) {
	pbMethod := paymentMethodToProto(paymentMethod)

	req := &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: pbMethod,
	}

	resp, err := c.paymentClient.PayOrder(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.TransactionUuid, nil
}

// paymentMethodToProto конвертирует доменный PaymentMethod в Protobuf.
func paymentMethodToProto(method model.PaymentMethod) paymentv1.PaymentMethod {
	switch method {
	case model.PaymentMethodCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCreditCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
