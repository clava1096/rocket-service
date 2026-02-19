package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type Postgres interface {
	URI() string
	MigrationsDir() string
}

type OrderConfig interface {
	Address() string
	HeaderTimeout() time.Duration
}

type GrpcClients interface {
	InventoryURI() string
	PaymentURI() string
}
