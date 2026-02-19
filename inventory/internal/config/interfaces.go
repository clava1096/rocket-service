package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type Mongo interface {
	URI() string
	DatabaseName() string
	CollectionName() string
}

type InventoryConfig interface {
	Address() string
}
