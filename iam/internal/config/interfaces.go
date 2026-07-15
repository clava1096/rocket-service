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

type IamConfig interface {
	Address() string
}

type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
	CacheTTL() time.Duration
}

type SessionConfig interface {
	TTL() time.Duration
	KeyPrefix() string
	MakeKey(sessionID string) string
	UserSessionsKeyPrefix() string
	MakeUserSessionsKey(userID string) string
}
