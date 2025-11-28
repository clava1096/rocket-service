package converter

import (
	"github.com/clava1096/rocket-service/order/internal/model"
	orderV1 "github.com/clava1096/rocket-service/shared/pkg/openapi/order/v1"
	"github.com/google/uuid"
)

// OrderToOpenAPI конвертирует доменную модель заказа в OpenAPI-DTO.
func OrderToOpenAPI(order model.Order) *orderV1.OrderDto {
	dto := &orderV1.OrderDto{
		OrderUUID:  MustParseUUID(order.UUID),
		UserUUID:   MustParseUUID(order.UserUUID),
		PartUuids:  stringSliceToUUIDSlice(order.PartUUIDs),
		TotalPrice: order.TotalPrice,
		Status:     orderStatusToOpenAPI(order.Status),
	}

	// Опциональные поля
	if order.TransactionUUID != nil {
		txnUUID := MustParseUUID(*order.TransactionUUID)
		dto.TransactionUUID = orderV1.NewOptNilUUID(txnUUID)
	}
	if order.PaymentMethod != nil {
		dto.PaymentMethod = orderV1.NewOptPaymentMethod(
			paymentMethodToOpenAPI(*order.PaymentMethod),
		)
	}

	return dto
}

// CreateOrderRequestToModel конвертирует OpenAPI-запрос в доменную модель.
func CreateOrderRequestToModel(req *orderV1.CreateOrderRequest) model.Order {
	return model.Order{
		UUID:      generateUUID(), // UUID заказа генерируется на сервере
		UserUUID:  req.UUID,
		PartUUIDs: req.PartsUUID,
	}
}

// PayOrderRequestToModel конвертирует OpenAPI-запрос оплаты в доменный способ оплаты.
func PayOrderRequestToModel(req *orderV1.PayOrderRequest) model.PaymentMethod {
	return PaymentMethodFromOpenAPI(req.PaymentMethod)
}

// --- Вспомогательные функции ---

func orderStatusToOpenAPI(status model.OrderStatus) orderV1.OrderStatus {
	switch status {
	case model.OrderStatusPendingPayment:
		return orderV1.OrderStatusPENDINGPAYMENT
	case model.OrderStatusPaid:
		return orderV1.OrderStatusPAID
	case model.OrderStatusCancelled:
		return orderV1.OrderStatusCANCELLED
	default:
		return orderV1.OrderStatusPENDINGPAYMENT
	}
}

func paymentMethodToOpenAPI(method model.PaymentMethod) orderV1.PaymentMethod {
	switch method {
	case model.PaymentMethodCard:
		return orderV1.PaymentMethodCARD
	case model.PaymentMethodSBP:
		return orderV1.PaymentMethodSBP
	case model.PaymentMethodCreditCard:
		return orderV1.PaymentMethodCREDITCARD
	case model.PaymentMethodInvestorMoney:
		return orderV1.PaymentMethodINVESTORMONEY
	default:
		return orderV1.PaymentMethodUNKNOWN
	}
}

func PaymentMethodFromOpenAPI(method orderV1.PaymentMethod) model.PaymentMethod {
	switch method {
	case orderV1.PaymentMethodCARD:
		return model.PaymentMethodCard
	case orderV1.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case orderV1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCreditCard
	case orderV1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnknown
	}
}

func stringSliceToUUIDSlice(ss []string) []uuid.UUID {
	uuids := make([]uuid.UUID, len(ss))
	for i, s := range ss {
		uuids[i] = MustParseUUID(s)
	}
	return uuids
}

func MustParseUUID(s string) uuid.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		// В production лучше возвращать ошибку, но для простоты — паника
		panic("invalid UUID: " + s)
	}
	return u
}

func generateUUID() string {
	return uuid.New().String()
}
