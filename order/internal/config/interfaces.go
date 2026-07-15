package config

import (
	"time"

	"github.com/IBM/sarama"
)

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type Postgres interface {
	URI() string
	MigrationsDir() string
}

type OrderConfig interface {
	Address() string
	HeaderTimeout() time.Duration
}

type GrpcClients interface {
	InventoryURI() string
	PaymentURI() string
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidProducer interface {
	Topic() string
	Config() *sarama.Config
}

type ShipAssembledConsumer interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}
