package config

import (
	"github.com/clava1096/rocket-service/iam/internal/config/env"

	"github.com/joho/godotenv"
)

var appConfig *config

type config struct {
	Server   IamConfig
	Postgres Postgres
	Logger   LoggerConfig
	Redis    RedisConfig
	Session  SessionConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)

	if err != nil {
		return err
	}

	serverCfg, err := env.NewIamEnvGrpc()
	if err != nil {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	redisCfg, err := env.NewRedisConfig()
	if err != nil {
		return err
	}

	sessionCfg, err := env.NewSessionConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Server:     serverCfg,
		PostgresDb: postgresCfg,
		Logger:     loggerCfg,
		Redis:      redisCfg,
		Session:    sessionCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
