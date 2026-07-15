package order_paid_consumer

import (
	"context"

	kafkaConverter "github.com/clava1096/rocket-service/notification/internal/converter/kafka"
	service2 "github.com/clava1096/rocket-service/notification/internal/service"
	"github.com/clava1096/rocket-service/platform/pkg/kafka"
	"github.com/clava1096/rocket-service/platform/pkg/logger"
	"go.uber.org/zap"
)

type service struct {
	orderPaidConsumer kafka.Consumer
	orderPaidDecoder  kafkaConverter.OrderPaidDecoder
	telegramService   service2.TelegramService
}

func NewService(orderPaidConsumer kafka.Consumer, orderPaidDecoder kafkaConverter.OrderPaidDecoder, telegramService service2.TelegramService) *service {
	return &service{
		orderPaidConsumer: orderPaidConsumer,
		orderPaidDecoder:  orderPaidDecoder,
		telegramService:   telegramService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order orderPaidConsumer service")

	err := s.orderPaidConsumer.Consume(ctx, s.OrderPaidHandler)

	if err != nil {
		logger.Info(ctx, "Consume from order.paid topic error", zap.Error(err))
		return err
	}

	return nil
}
