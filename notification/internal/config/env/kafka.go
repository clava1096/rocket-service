package env

import "github.com/caarlos0/env/v11"

type kafkaEnvConfig struct {
	Brokers []string `env:"KAFKA_BROKERS,required"`
}

type kafkaConfig struct {
	raw kafkaEnvConfig
}

func NewKafkaEnvConfig() (*kafkaConfig, error) {
	var raw kafkaEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &kafkaConfig{raw: raw}, nil
}

func (conf *kafkaConfig) Brokers() []string {
	return conf.raw.Brokers
}
