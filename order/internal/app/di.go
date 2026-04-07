package app

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/clava1096/rocket-service/order/internal/config"
	"github.com/clava1096/rocket-service/order/internal/converter/kafka/decoder"
	"github.com/clava1096/rocket-service/order/internal/service/consumer/order_consumer"
	"github.com/clava1096/rocket-service/order/internal/service/producer/order_producer"
	"github.com/clava1096/rocket-service/platform/pkg/kafka"
	"github.com/clava1096/rocket-service/platform/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderAPI "github.com/clava1096/rocket-service/order/internal/api/order/v1"
	def "github.com/clava1096/rocket-service/order/internal/client/grpc"
	inventoryv1Client "github.com/clava1096/rocket-service/order/internal/client/grpc/inventory/v1"
	paymentv1Client "github.com/clava1096/rocket-service/order/internal/client/grpc/payment/v1"
	kafkaConverter "github.com/clava1096/rocket-service/order/internal/converter/kafka"
	"github.com/clava1096/rocket-service/order/internal/repository"
	orderRepository "github.com/clava1096/rocket-service/order/internal/repository/order"
	orderService "github.com/clava1096/rocket-service/order/internal/service"
	orderServiceImpl "github.com/clava1096/rocket-service/order/internal/service/order"
	wrappedKafkaConsumer "github.com/clava1096/rocket-service/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/clava1096/rocket-service/platform/pkg/kafka/producer"
	kafkaMiddleware "github.com/clava1096/rocket-service/platform/pkg/middleware/kafka"
	orderV1 "github.com/clava1096/rocket-service/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/clava1096/rocket-service/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/clava1096/rocket-service/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1      *orderV1.Server
	orderService orderService.OrderService

	orderConsumerService orderService.OrderConsumerService
	consumerGroup        sarama.ConsumerGroup

	orderRepository repository.OrderRepository

	postgresPool    *pgxpool.Pool
	inventoryConn   *grpc.ClientConn
	inventoryClient def.InventoryClient

	paymentConn   *grpc.ClientConn
	paymentClient def.PaymentClient

	orderShipAssembledConsumer kafka.Consumer
	orderProducerService       orderService.OrderProducerService

	orderAssembledDecoder kafkaConverter.ShipAssembledDecoder
	orderPaidDecoder      kafkaConverter.OrderPaidDecoder
	orderPaidProducer     kafka.Producer
	syncProducer          sarama.SyncProducer
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
			d.OrderProducerService(),
		)
	}
	return d.orderService
}

func (d *diContainer) OrderProducerService() orderService.OrderProducerService {
	if d.orderPaidProducer == nil {
		d.orderProducerService = order_producer.NewService(d.OrderRecorderProducer())
	}

	return d.orderProducerService
}

func (d *diContainer) OrderRecorderProducer() kafka.Producer {

	if d.orderPaidProducer == nil {
		d.orderPaidProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderPaidProducer.Topic(),
			logger.Logger())
	}

	return d.orderPaidProducer
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {

	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().KafkaConfig.Brokers(),
			config.AppConfig().OrderPaidProducer.Config())
		if err != nil {
			panic(fmt.Sprintf("sarama.NewSyncProducer failed: %v", zap.Error(err)))
		}

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) OrderConsumerService(ctx context.Context) orderService.OrderConsumerService {
	if d.orderConsumerService == nil {
		d.orderConsumerService = order_consumer.NewService(d.OrderShipAssembledConsumer(), d.OrderAssembledDecoder(), d.OrderRepository(ctx))
	}

	return d.orderConsumerService
}

func (d *diContainer) OrderShipAssembledConsumer() kafka.Consumer {
	if d.orderShipAssembledConsumer == nil {
		d.orderShipAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().ShipAssembledConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()))
	}

	return d.orderShipAssembledConsumer
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().KafkaConfig.Brokers(),
			config.AppConfig().ShipAssembledConsumer.GroupID(),
			config.AppConfig().ShipAssembledConsumer.Config())

		if err != nil {
			panic(fmt.Sprintf("sarama.NewConsumer failed: %v", zap.Error(err)))
		}
		d.consumerGroup = consumerGroup
	}
	return d.consumerGroup
}

func (d *diContainer) OrderAssembledDecoder() kafkaConverter.ShipAssembledDecoder {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = decoder.NewShipAssembledDecoder()
	}
	return d.orderAssembledDecoder
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
