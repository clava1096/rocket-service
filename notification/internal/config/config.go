package config

import (
	"github.com/clava1096/rocket-service/notification/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	KafkaConfig                 KafkaConfig
	LoggerConfig                LoggerConfig
	OrderAssemblyConsumerConfig OrderAssemblyConsumerConfig
	OrderPaidConsumerConfig     OrderPaidConsumerConfig
	TelegramConfig              TelegramConfig
	ProxyConfig                 Proxy
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	println()

	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaEnvConfig()
	if err != nil {
		return err
	}

	loggerCfg, err := env.NewLoggerEnvConfig()
	if err != nil {
		return err
	}

	orderAssemblyConsumerCfg, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	orderPaidConsumerCfg, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	TelegramCfg, err := env.NewTelegramConfig()
	if err != nil {
		return err
	}

	ProxyCfg, err := env.NewProxyEnvConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		KafkaConfig:                 kafkaCfg,
		LoggerConfig:                loggerCfg,
		TelegramConfig:              TelegramCfg,
		OrderPaidConsumerConfig:     orderPaidConsumerCfg,
		OrderAssemblyConsumerConfig: orderAssemblyConsumerCfg,
		ProxyConfig:                 ProxyCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
