package kafka

import (
	"context"

	"github.com/clava1096/rocket-service/platform/pkg/kafka"
	"github.com/clava1096/rocket-service/platform/pkg/kafka/consumer"
	"go.uber.org/zap"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
}

func Logging(logger Logger) consumer.Middleware {
	return func(next kafka.MessageHandler) kafka.MessageHandler {
		return func(ctx context.Context, msg kafka.Message) error {
			logger.Info(ctx, "Message received", zap.Any("topic", msg.Topic))
			return next(ctx, msg)
		}
	}
}
