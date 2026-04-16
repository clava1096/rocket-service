package consumer

import (
	"github.com/IBM/sarama"
	"github.com/clava1096/rocket-service/platform/pkg/kafka"
	"go.uber.org/zap"
)

type Middleware func(next kafka.MessageHandler) kafka.MessageHandler

type GroupHandler struct {
	handler kafka.MessageHandler
	logger  Logger
}

func NewGroupHandler(handler kafka.MessageHandler, logger Logger, middlewares ...Middleware) *GroupHandler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return &GroupHandler{
		handler: handler,
		logger:  logger,
	}
}

func (g *GroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (g *GroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (g *GroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				g.logger.Info(session.Context(), "Kafka message channel close")
				return nil
			}

			msg := kafka.Message{
				Key:            message.Key,
				Value:          message.Value,
				Topic:          message.Topic,
				Partition:      message.Partition,
				Offset:         message.Offset,
				Timestamp:      message.Timestamp,
				BlockTimestamp: message.BlockTimestamp,
				Headers:        extractHeaders(message.Headers),
			}

			if err := g.handler(session.Context(), msg); err != nil {
				g.logger.Error(session.Context(), "Kafka handler error", zap.Error(err))
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			g.logger.Info(session.Context(), "Kafka context done")
			return nil
		}
	}

}

func extractHeaders(headers []*sarama.RecordHeader) map[string][]byte {
	res := make(map[string][]byte)
	for _, header := range headers {
		if header != nil && header.Key != nil {
			res[string(header.Key)] = header.Value
		}
	}
	return res
}
