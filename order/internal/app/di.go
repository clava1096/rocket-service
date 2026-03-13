package app

import (
	"context"
	"log"

	"github.com/clava1096/rocket-service/order/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderAPI "github.com/clava1096/rocket-service/order/internal/api/order/v1"
	def "github.com/clava1096/rocket-service/order/internal/client/grpc"
	inventoryv1Client "github.com/clava1096/rocket-service/order/internal/client/grpc/inventory/v1"
	paymentv1Client "github.com/clava1096/rocket-service/order/internal/client/grpc/payment/v1"
	"github.com/clava1096/rocket-service/order/internal/repository"
	orderRepository "github.com/clava1096/rocket-service/order/internal/repository/order"
	orderService "github.com/clava1096/rocket-service/order/internal/service"
	orderServiceImpl "github.com/clava1096/rocket-service/order/internal/service/order"
	orderV1 "github.com/clava1096/rocket-service/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/clava1096/rocket-service/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/clava1096/rocket-service/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1         *orderV1.Server
	orderService    orderService.OrderService
	orderRepository repository.OrderRepository

	postgresPool    *pgxpool.Pool
	inventoryConn   *grpc.ClientConn
	inventoryClient def.InventoryClient

	paymentConn   *grpc.ClientConn
	paymentClient def.PaymentClient
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderServer(ctx context.Context) orderService.OrderService {
	var err error
	if d.orderV1 == nil {
		api := orderAPI.NewAPI(d.OrderService(ctx))
		d.orderV1, err = orderV1.NewServer(api)
		if err != nil {
			panic(err)
		}
	}
	return d.orderService
}

func (d *diContainer) OrderService(ctx context.Context) orderService.OrderService {
	if d.orderService == nil {
		d.orderService = orderServiceImpl.NewService(
			d.OrderRepository(ctx),
			d.InventoryClient(ctx),
			d.PaymentClient(ctx),
		)
	}
	return d.orderService
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewRepository(d.PostgresPool(ctx))
	}
	return d.orderRepository
}

func (d *diContainer) PostgresPool(ctx context.Context) *pgxpool.Pool {
	if d.postgresPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}

		err = pool.Ping(ctx)
		if err != nil {
			log.Fatalf("Error pinging database: %v", err)
		}
		d.postgresPool = pool
	}
	return d.postgresPool
}

func (d *diContainer) PaymentClient(_ context.Context) def.PaymentClient {
	if d.paymentClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().GrpcClients.PaymentURI(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatalf("Failed to connect to payment service: %v", err)
		}

		rawClient := paymentv1.NewPaymentServiceClient(conn)

		d.paymentConn = conn
		d.paymentClient = paymentv1Client.NewClient(rawClient)
	}

	return d.paymentClient
}

func (d *diContainer) InventoryClient(_ context.Context) def.InventoryClient {
	if d.inventoryClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().GrpcClients.InventoryURI(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatalf("failed to connect to inventory service: %v", err)
		}

		rawClient := inventoryv1.NewInventoryServiceClient(conn)

		d.inventoryConn = conn
		d.inventoryClient = inventoryv1Client.NewClient(rawClient)
	}

	return d.inventoryClient
}
