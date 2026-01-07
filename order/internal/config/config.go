package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   Server   `yaml:"server"`
	Postgres Postgres `yaml:"postgres"`
}

type Server struct {
	HttpPort string `yaml:"http_port"`

	InventoryGrpcPort string `yaml:"inventory_grpc_port"`
	PaymentGrpcPort   string `yaml:"payment_grpc_port"`
}

type Postgres struct {
	Host          string `yaml:"host"`
	Port          string `yaml:"port"`
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	Database      string `yaml:"database"`
	MigrationsDir string `yaml:"migrations_dir"`
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

func (p *Postgres) GetPostgresUri() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		p.User,
		p.Password,
		p.Host,
		p.Port,
		p.Database,
	)
}
