package config

import (
	"github.com/clava1096/rocket-service/payment/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Logger  LoggerConfig
	Payment PaymentConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)

	logger, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	payment, err := env.NewPaymentEnvConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:  logger,
		Payment: payment,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
