package model

import (
	"context"
	"github.com/clava1096/rocket-service/order/internal/model"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type RepositorySuite struct {
	suite.Suite
	repo *repository
	ctx  context.Context
}

func (s *RepositorySuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = NewRepository()
}

func TestRepository(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (s *RepositorySuite) newOrder() model.Order {
	return model.Order{
		UUID:       "test-uuid",
		UserUUID:   "user-123",
		PartUUIDs:  []string{"part-1", "part-2"},
		TotalPrice: 100.5,
		Status:     model.OrderStatusPendingPayment,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
