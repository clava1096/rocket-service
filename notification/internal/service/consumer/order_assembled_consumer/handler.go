package order_assembled_consumer

import (
	"context"

	"github.com/clava1096/rocket-service/notification/internal/model"
	"github.com/clava1096/rocket-service/platform/pkg/kafka"
	"github.com/clava1096/rocket-service/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) OrderAssemblyHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderAssembledDecoder.Decode(msg.Value)

	if err != nil {
		logger.Info(ctx, "Error decoding event", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Received event", zap.Any("event", event))

	return s.ProcessMessage(ctx, event)
}

func (s *service) ProcessMessage(ctx context.Context, event model.ShipAssembledEvent) error {
	logger.Info(ctx, "Processing order assembly event",
		zap.String("order_uuid", event.OrderUUID))

	if err := s.telegramService.SendAssembledNotification(ctx, event); err != nil {
		logger.Error(ctx, "Failed to send order paid notification", zap.Error(err))
	}

	logger.Info(ctx, "Telegram Notification sent",
		zap.String("order_uuid", event.OrderUUID))

	return nil
}
