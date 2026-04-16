package consumer

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/clava1096/rocket-service/platform/pkg/kafka"
	"go.uber.org/zap"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type consumer struct {
	ConsumerGroup sarama.ConsumerGroup
	topics        []string
	Logger        Logger
	middlewares   []Middleware
}

func NewConsumer(group sarama.ConsumerGroup, topics []string, logger Logger, middlewares ...Middleware) *consumer {
	return &consumer{
		ConsumerGroup: group,
		topics:        topics,
		Logger:        logger,
		middlewares:   middlewares,
	}
}

func (c *consumer) Consume(ctx context.Context, handler kafka.MessageHandler) error {
	newGroupHandler := NewGroupHandler(handler, c.Logger, c.middlewares...)

	for {
		if err := c.ConsumerGroup.Consume(ctx, c.topics, newGroupHandler); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}
			c.Logger.Error(ctx, "Kafka consume error", zap.Error(err))

			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		c.Logger.Info(ctx, "Kafka consumer group rebalancing...")
	}
}
