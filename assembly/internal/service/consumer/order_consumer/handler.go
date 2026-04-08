package order_consumer

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/clava1096/rocket-service/assembly/internal/model"
	"github.com/clava1096/rocket-service/platform/pkg/kafka"
	"github.com/clava1096/rocket-service/platform/pkg/logger"
	eventsv1 "github.com/clava1096/rocket-service/shared/pkg/proto/events/v1"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *service) HandleOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	go s.processAssembly(ctx, event)
	return nil
}

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)

	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid")
		return err
	}

	return s.HandleOrderPaid(ctx, event)

}

func (s *service) processAssembly(ctx context.Context, event model.OrderPaidEvent) {

	delaySec := rand.IntN(9) + 1
	delay := time.Duration(delaySec) * time.Second

	logger.Info(ctx, "Assembly started, waiting...",
		zap.String("order_uuid", event.OrderUUID),
		zap.Int("delay_sec", delaySec))

	select {
	case <-time.After(delay):

	case <-ctx.Done():
		logger.Warn(ctx, "Assembly canceled")
		return
	}

	assemblyEvent := &eventsv1.ShipAssembled{
		EventUuid:    uuid.New().String(),
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: int64(delaySec),
	}

	payload, err := proto.Marshal(assemblyEvent)
	if err != nil {
		logger.Error(ctx, "Failed to marshal assembly event", zap.Error(err))
		return
	}

	err = s.producer.Send(ctx, []byte(event.OrderUUID), payload)
	if err != nil {
		logger.Error(ctx, "Failed to send to Kafka",
			zap.Error(err))
		return
	}

	logger.Info(ctx, "ShipAssembled sent",
		zap.String("order_uuid", event.OrderUUID))
}
