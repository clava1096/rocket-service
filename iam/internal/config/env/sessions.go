package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type sessionEnvConfig struct {
	TTL                time.Duration `env:"SESSION_TTL,required" envDefault:"24h"`
	KeyPrefix          string        `env:"SESSION_KEY_PREFIX" envDefault:"session:"`
	UserSessionsPrefix string        `env:"USER_SESSIONS_PREFIX" envDefault:"user_sessions:"`
}

type sessionConfig struct {
	raw sessionEnvConfig
}

func NewSessionConfig() (*sessionConfig, error) {
	var raw sessionEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &sessionConfig{raw: raw}, nil
}

func (cfg *sessionConfig) TTL() time.Duration {
	return cfg.raw.TTL
}

func (cfg *sessionConfig) KeyPrefix() string {
	return cfg.raw.KeyPrefix
}

func (cfg *sessionConfig) MakeKey(sessionID string) string {
	return cfg.raw.KeyPrefix + sessionID
}

func (cfg *sessionConfig) UserSessionsKeyPrefix() string {
	return cfg.raw.UserSessionsPrefix
}

func (cfg *sessionConfig) MakeUserSessionsKey(userID string) string {
	return cfg.raw.UserSessionsPrefix + userID
}
