package model

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/clava1096/rocket-service/order/internal/model"
	"github.com/clava1096/rocket-service/platform/pkg/logger"
	"github.com/clava1096/rocket-service/platform/pkg/testcontainers/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
)

type RepositorySuite struct {
	suite.Suite

	ctx    context.Context
	cancel context.CancelFunc

	pool        *pgxpool.Pool
	repo        *repository
	pgContainer *postgres.Container
}

func (s *RepositorySuite) SetupSuite() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 5*time.Minute)

	pgContainer, err := postgres.NewContainer(s.ctx,
		postgres.WithImageName("postgres:18.2"),
		postgres.WithDatabase("test_orders"),
		postgres.WithAuth("test_user", "test_password"),
		postgres.WithLogger(logger.Logger()))

	s.NoError(err)
	s.pgContainer = pgContainer

	uri := pgContainer.URI()
	s.pool, err = pgxpool.New(s.ctx, uri)
	s.NoError(err)

	sqlDB := stdlib.OpenDB(*s.pool.Config().ConnConfig)
	defer sqlDB.Close()

	err = goose.Up(sqlDB, "../../../migrations")
	s.NoError(err, "failed to apply migrations")

	s.repo = NewRepository(s.pool)
}

func (s *RepositorySuite) SetupTest() {
	_, err := s.pool.Exec(s.ctx, "TRUNCATE TABLE orders RESTART IDENTITY CASCADE")
	s.NoError(err)
}

func (s *RepositorySuite) TearDownSuite() {
	if s.pool != nil {
		s.pool.Close()
	}
	if s.pgContainer != nil {
		_ = s.pgContainer.Terminate(s.ctx)
	}
	if s.cancel != nil {
		s.cancel()
	}
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
