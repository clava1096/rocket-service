package config

import (
	"github.com/clava1096/rocket-service/inventory/internal/config/env"
	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Server  InventoryConfig
	MongoDb Mongo
	Logger  LoggerConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)

	if err != nil {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	mongoCfg, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	inventoryCfg, err := env.NewInventoryEnvGrpc()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:  loggerCfg,
		MongoDb: mongoCfg,
		Server:  inventoryCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
