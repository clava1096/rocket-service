package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig `yaml:"server"`
	MongoDb MongoDb      `yaml:"mongodb"`
}

type ServerConfig struct {
	GrpcPort string `yaml:"grpc_port"`
}

type MongoDb struct {
	Username   string `yaml:"init_root_username"`
	Password   string `yaml:"init_root_password"`
	Database   string `yaml:"init_database"`
	Port       string `yaml:"port"`
	Host       string `yaml:"host"`
	AuthDb     string `yaml:"auth_db"`
	Collection string `yaml:"collection"`
}

func GetConfig() (*Config, error) {
	var config Config

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	configPath := filepath.Join(dir, "config.yaml")

	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (m *MongoDb) GetMongodbUri() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s",
		m.Username,
		m.Password,
		m.Host,
		m.Port,
		m.Database,
		m.AuthDb)
}
