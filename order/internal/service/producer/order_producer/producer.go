package order_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/clava1096/rocket-service/order/internal/model"
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

func (p *service) ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	msg := &eventsv1.OrderPaid{
		EventUuid:       event.EventUUID,
		OrderUuid:       event.OrderUUID,
		UserUuid:        event.UserUUID,
		PaymentMethod:   event.PaymentMethod,
		TransactionUuid: event.TransactionUUID,
	}

	payload, err := proto.Marshal(msg)

	if err != nil {
		logger.Error(ctx, "Failed to marshal order payload", zap.Error(err))
		return err
	}

	err = p.orderAssemblyProducer.Send(ctx, []byte(event.OrderUUID), payload)

	if err != nil {
		logger.Error(ctx, "Failed to publish order payload", zap.Error(err))
		return err
	}

	logger.Info(ctx, "OrderPaid event sent successfully",
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID))

	return nil
}
