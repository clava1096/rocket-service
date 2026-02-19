package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type mongoEnvConfig struct {
	Username     string `env:"MONGO_INITDB_ROOT_USERNAME,required"`
	Password     string `env:"MONGO_INITDB_ROOT_PASSWORD,required"`
	Database     string `env:"MONGO_DATABASE"`
	Port         string `env:"MONGO_PORT"`
	ExternalPort string `env:"MONGO_EXTERNAL_PORT"`
	Hostname     string `env:"MONGO_HOSTNAME"`
	Host         string `env:"MONGO_HOST"`
	AuthDB       string `env:"MONGO_AUTH_DB"`
	Collection   string `env:"MONGO_COLLECTION"`
}

type mongoConfig struct {
	raw mongoEnvConfig
}

func NewMongoConfig() (*mongoConfig, error) {
	var raw mongoEnvConfig

	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &mongoConfig{raw: raw}, nil
}

func (mongo *mongoConfig) URI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s",
		mongo.raw.Username,
		mongo.raw.Password,
		mongo.raw.Host,
		mongo.raw.Port,
		mongo.raw.Database,
		mongo.raw.AuthDB)
}

func (mongo *mongoConfig) DatabaseName() string {
	return mongo.raw.Database
}

func (mongo *mongoConfig) CollectionName() string {
	return mongo.raw.Collection
}
