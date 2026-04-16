package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderShipConsumerEnvConfig struct {
	Topic   string `env:"ORDER_SHIP_ASSEMBLED_TOPIC_NAME,required"`
	GroupID string `env:"ORDER_SHIP_ASSEMBLED_GROUPID,required"`
}

type orderShipConsumerConfig struct {
	raw orderShipConsumerEnvConfig
}

func NewOrderShipAssembledConsumerConfig() (*orderShipConsumerConfig, error) {
	var raw orderShipConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderShipConsumerConfig{raw: raw}, nil
}

func (conf *orderShipConsumerConfig) Topic() string {
	return conf.raw.Topic
}

func (conf *orderShipConsumerConfig) GroupID() string {
	return conf.raw.GroupID
}

func (conf *orderShipConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
