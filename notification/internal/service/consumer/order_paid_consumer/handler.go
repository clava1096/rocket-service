package order_paid_consumer

import (
	"context"

	"github.com/clava1096/rocket-service/notification/internal/model"
	"github.com/clava1096/rocket-service/platform/pkg/kafka"
	"github.com/clava1096/rocket-service/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) OrderPaidHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)

	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid")
		return err
	}

	return s.ProcessMessage(ctx, event)
}

func (s *service) ProcessMessage(ctx context.Context, event model.OrderPaidEvent) error {
	logger.Info(ctx, "Processing order paid event",
		zap.String("order_uuid", event.OrderUUID))

	if err := s.telegramService.SendPaidNotification(ctx, event); err != nil {
		logger.Error(ctx, "Failed to send order paid notification", zap.Error(err))
	}

	logger.Info(ctx, "Telegram Notification sent",
		zap.String("order_uuid", event.OrderUUID))

	return nil
}
