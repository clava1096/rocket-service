package service

import (
	"context"
)

type PaymentService interface {
	Pay(ctx context.Context) string
}
