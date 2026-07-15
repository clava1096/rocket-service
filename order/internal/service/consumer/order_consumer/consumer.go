package order_consumer

import (
	"context"

	kafkaConverter "github.com/clava1096/rocket-service/order/internal/converter/kafka"
	"github.com/clava1096/rocket-service/order/internal/repository"
	"github.com/clava1096/rocket-service/platform/pkg/kafka"
	"github.com/clava1096/rocket-service/platform/pkg/logger"
	"go.uber.org/zap"
)

type service struct {
	orderPaidConsumer    kafka.Consumer
	producer             kafka.Producer
	shipAssembledDecoder kafkaConverter.ShipAssembledDecoder
	orderRepo            repository.OrderRepository
}

func NewService(orderPaidConsumer kafka.Consumer, shipAssembledDecoder kafkaConverter.ShipAssembledDecoder, orderRepo repository.OrderRepository) *service {
	return &service{
		orderPaidConsumer:    orderPaidConsumer,
		shipAssembledDecoder: shipAssembledDecoder,
		orderRepo:            orderRepo,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order orderPaidConsumer service")

	err := s.orderPaidConsumer.Consume(ctx, s.OrderHandler)

	if err != nil {
		logger.Error(ctx, "Consume from order.paid topic error", zap.Error(err))
		return err
	}

	return nil
}
