package order_consumer

import (
	"context"

	"github.com/clava1096/rocket-service/order/internal/model"
	"github.com/clava1096/rocket-service/platform/pkg/kafka"
	"github.com/clava1096/rocket-service/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.shipAssembledDecoder.Decode(msg.Value)

	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid")
		return err
	}

	return s.HandleShipAssembled(ctx, event)
}

func (s *service) HandleShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error {
	err := s.orderRepo.UpdateStatus(ctx, event.OrderUUID, model.OrderStatusAssembled)

	if err != nil {
		logger.Error(ctx, "Failed to update order status", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Order marked as ASSEMBLED", zap.String("order_uuid", event.OrderUUID))
	return nil
}
