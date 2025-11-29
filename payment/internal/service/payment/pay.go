package payment

import (
	"context"
	"github.com/google/uuid"
)

func (s *service) Pay(ctx context.Context) string {
	transactionUuid := uuid.New().String()
	return transactionUuid
}
