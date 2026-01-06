package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

type PaymentMethod string

const (
	PaymentMethodUnknown       PaymentMethod = "UNKNOWN"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

type Order struct {
	UUID            uuid.UUID      `db:"uuid"`
	UserUUID        uuid.UUID      `db:"user_uuid"`
	PartUUIDs       []uuid.UUID    `db:"part_uuids"`
	TotalPrice      float64        `db:"total_price"`
	Status          OrderStatus    `db:"status"`
	TransactionUUID *string        `db:"transaction_uuid"`
	PaymentMethod   *PaymentMethod `db:"payment_method"`
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at"`
}
