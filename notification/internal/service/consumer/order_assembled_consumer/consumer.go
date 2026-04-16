package order_assembled_consumer

import (
	"context"

	kafkaConverter "github.com/clava1096/rocket-service/notification/internal/converter/kafka"
	service2 "github.com/clava1096/rocket-service/notification/internal/service"
	"github.com/clava1096/rocket-service/platform/pkg/kafka"
	"github.com/clava1096/rocket-service/platform/pkg/logger"
	"go.uber.org/zap"
)

type service struct {
	orderAssembledConsumer kafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssemblyDecoder
	telegramService        service2.TelegramService
}

func NewService(orderAssembledConsumer kafka.Consumer, orderAssembledDecoder kafkaConverter.OrderAssemblyDecoder, telegramService service2.TelegramService) *service {
	return &service{
		orderAssembledConsumer: orderAssembledConsumer,
		orderAssembledDecoder:  orderAssembledDecoder,
		telegramService:        telegramService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order orderPaidConsumer service")

	err := s.orderAssembledConsumer.Consume(ctx, s.OrderAssemblyHandler)

	if err != nil {
		logger.Info(ctx, "Consume from order.paid topic error", zap.Error(err))
		return err
	}

	return nil
}
