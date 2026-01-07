package model

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"

	"github.com/clava1096/rocket-service/order/internal/config"
	"github.com/clava1096/rocket-service/order/internal/model"
)

type RepositorySuite struct {
	suite.Suite
	repo *repository
	ctx  context.Context
}

func (s *RepositorySuite) SetupTest() {
	s.ctx = context.Background()

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal("Error loading config...")
	}
	pool, err := pgxpool.New(s.ctx, cfg.Postgres.GetPostgresUri())

	s.repo = NewRepository(pool)
}

func TestRepository(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (s *RepositorySuite) newOrder() model.Order {
	return model.Order{
		UUID:       gofakeit.UUID(),
		UserUUID:   gofakeit.UUID(),
		PartUUIDs:  []string{gofakeit.UUID(), gofakeit.UUID()},
		TotalPrice: 100.5,
		Status:     model.OrderStatusPendingPayment,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
