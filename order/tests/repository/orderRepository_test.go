//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/clava1096/rocket-service/order/internal/migrator"
	"github.com/clava1096/rocket-service/platform/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/clava1096/rocket-service/order/internal/model"
	"github.com/clava1096/rocket-service/order/internal/repository"
	orderRepository "github.com/clava1096/rocket-service/order/internal/repository/order"
	"github.com/clava1096/rocket-service/platform/pkg/testcontainers/postgres"
)

func TestRepositoryIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Order Repository Integration Suite")
}

var _ = Describe("OrderRepository", Ordered, func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		repo   repository.OrderRepository
		//pool        *pgxpool.Pool
		pgContainer *postgres.Container
	)

	BeforeAll(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Minute)

		var err error
		pgContainer, err = postgres.NewContainer(ctx,
			postgres.WithDatabase("order_test"),
			postgres.WithAuth("test_user", "test_pass"),
			postgres.WithLogger(logger.Logger()),
		)
		Expect(err).NotTo(HaveOccurred(), "Failed to start PostgreSQL container")

		testDBURI := pgContainer.URI()

		pool, err := pgxpool.New(ctx, testDBURI)
		Expect(err).NotTo(HaveOccurred())

		pgxConfig, err := pgx.ParseConfig(testDBURI)
		Expect(err).NotTo(HaveOccurred(), "Failed to parse test DB URI")

		sqlDB := stdlib.OpenDB(*pgxConfig)
		defer sqlDB.Close()
		migrationsDir := "../../migrations"

		migratorRunner := migrator.NewMigrator(
			sqlDB,
			migrationsDir,
		)
		err = migratorRunner.Up()

		Expect(err).NotTo(HaveOccurred(), "Failed to run migrations")

		repo = orderRepository.NewRepository(pool)
	})

	AfterEach(func() {
		_, err := pgContainer.Client().Exec(ctx, "TRUNCATE TABLE orders RESTART IDENTITY CASCADE;")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterAll(func() {
		cancel()
		if pgContainer != nil {
			_ = pgContainer.Terminate(context.Background())
		}
	})

	Describe("Create", func() {
		It("should create a new order and persist it to PostgreSQL", func() {
			newOrder := model.Order{
				UserUUID:   gofakeit.UUID(),
				Status:     model.OrderStatusPendingPayment,
				TotalPrice: 1234.56,
			}

			created, err := repo.Create(ctx, newOrder)

			Expect(err).NotTo(HaveOccurred())
			Expect(created.UUID).NotTo(BeEmpty())
			Expect(created.UserUUID).To(Equal(newOrder.UserUUID))
			Expect(created.Status).To(Equal(newOrder.Status))

			var exists bool
			err = pgContainer.Client().QueryRow(ctx,
				"SELECT EXISTS(SELECT 1 FROM orders WHERE uuid = $1)",
				created.UUID,
			).Scan(&exists)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
		})
	})
})
