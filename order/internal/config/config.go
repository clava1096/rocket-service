package config

import (
	"github.com/clava1096/rocket-service/order/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Server                OrderConfig
	Postgres              Postgres
	GrpcClients           GrpcClients
	Logger                LoggerConfig
	OrderPaidProducer     OrderPaidProducer
	ShipAssembledConsumer ShipAssembledConsumer
	KafkaConfig           KafkaConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)

	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()

	if err != nil {
		return err
	}

	clientsGrpcCfg, err := env.NewGrpcConfig()

	if err != nil {
		return err
	}

	orderCfg, err := env.NewServerEnvConfig()
	if err != nil {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	OrderPaidProducerCfg, err := env.NewOrderAssembledProducerConfig()
	if err != nil {
		return err
	}

	ShipAssembledConsumerCfg, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaEnvConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Server:                orderCfg,
		Postgres:              postgresCfg,
		GrpcClients:           clientsGrpcCfg,
		Logger:                loggerCfg,
		OrderPaidProducer:     OrderPaidProducerCfg,
		ShipAssembledConsumer: ShipAssembledConsumerCfg,
		KafkaConfig:           kafkaCfg,
	}
	return nil
}

func AppConfig() *config {
	return appConfig
}
