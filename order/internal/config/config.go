package config

import (
	"github.com/clava1096/rocket-service/order/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Server      OrderConfig
	Postgres    Postgres
	GrpcClients GrpcClients
	Logger      LoggerConfig
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

	appConfig = &config{
		Server:      orderCfg,
		Postgres:    postgresCfg,
		GrpcClients: clientsGrpcCfg,
		Logger:      loggerCfg,
	}
	return nil
}

func AppConfig() *config {
	return appConfig
}
