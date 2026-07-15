package order_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/clava1096/rocket-service/assembly/internal/model"
	"github.com/clava1096/rocket-service/platform/pkg/kafka"
	"github.com/clava1096/rocket-service/platform/pkg/logger"
	eventsv1 "github.com/clava1096/rocket-service/shared/pkg/proto/events/v1"
)

type service struct {
	orderAssemblyProducer kafka.Producer
}

func NewService(orderAssemblyProducer kafka.Producer) *service {
	return &service{
		orderAssemblyProducer: orderAssemblyProducer,
	}
}

func (p *service) ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error {
	msg := &eventsv1.ShipAssembled{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	payload, err := proto.Marshal(msg)

	if err != nil {
		logger.Error(ctx, "Failed to marshal order payload", zap.Error(err))
	}

	err = p.orderAssemblyProducer.Send(ctx, []byte(event.OrderUUID), payload)

	if err != nil {
		logger.Error(ctx, "Failed to publish order payload", zap.Error(err))
		return err
	}

	logger.Info(ctx, "ShipAssembled event sent successfully",
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.Int64("build_time_sec", event.BuildTimeSec))

	return nil
}
