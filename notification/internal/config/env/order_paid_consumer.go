package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderPaidConsumerEnvConfig struct {
	Topic   string `env:"ORDER_PAID_TOPIC,required"`
	GroupID string `env:"ORDER_PAID_GROUP_ID,required"`
}

type orderPaidConsumerConfig struct {
	raw orderPaidConsumerEnvConfig
}

func NewOrderPaidConsumerConfig() (*orderPaidConsumerConfig, error) {
	var raw orderPaidConsumerEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderPaidConsumerConfig{raw: raw}, nil
}

func (conf *orderPaidConsumerConfig) Topic() string {
	return conf.raw.Topic
}

func (conf *orderPaidConsumerConfig) GroupID() string {
	return conf.raw.GroupID
}

func (conf *orderPaidConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
