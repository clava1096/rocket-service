package config

import (
	"os"

	"github.com/clava1096/rocket-service/assembly/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Logger                 LoggerConfig
	Kafka                  KafkaConfig
	OrderAssembledProducer OrderAssembledProducer
	OrderPaidConsumer      OrderPaidConsumer
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	OrderAssembledProducerCfg, err := env.NewOrderAssembledProducerConfig()
	if err != nil {
		return err
	}

	OrderPaidConsumerCfg, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerCfg,
		Kafka:                  kafkaCfg,
		OrderAssembledProducer: OrderAssembledProducerCfg,
		OrderPaidConsumer:      OrderPaidConsumerCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
